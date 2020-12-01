package validator_signer

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/herumi/bls-eth-go-binary/bls"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
	prot "github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/bloxapp/eth2-key-manager/stores/in_memory"
)

func inmemStorage() *in_memory.InMemStore {
	return in_memory.NewInMemStore(core.MainNetwork)
}

func setupNoSlashingProtection(seed []byte) (ValidatorSigner, error) {
	noProtection := &prot.NoProtection{}
	store := inmemStorage()
	wallet, err := walletWithSeed(seed, store)
	if err != nil {
		return nil, err
	}
	return NewSimpleSigner(wallet, noProtection, core.PyrmontNetwork), nil
}

func setupWithSlashingProtection(seed []byte, setLatestAttestation bool) (ValidatorSigner, error) {
	store := inmemStorage()
	protector := prot.NewNormalProtection(store)
	wallet, err := walletWithSeed(seed, store)
	if err != nil {
		return nil, err
	}

	// update highest attestation
	acc, err := wallet.AccountByPublicKey("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")
	if err != nil {
		log.Fatal(err)
	}
	if setLatestAttestation {
		protector.UpdateLatestAttestation(acc.ValidatorPublicKey(), &eth.AttestationData{
			Slot:            0,
			CommitteeIndex:  0,
			BeaconBlockRoot: ignoreError(hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000")).([]byte),
			},
		})
	}
	return NewSimpleSigner(wallet, protector, core.PyrmontNetwork), nil
}

func walletWithSeed(seed []byte, store core.Storage) (core.Wallet, error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil, err
	}

	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	vault, err := eth2keymanager.NewKeyVault(options)
	if err != nil {
		return nil, err
	}

	wallet, err := vault.Wallet()
	if err != nil {
		return nil, err
	}

	_, err = wallet.CreateValidatorAccount(seed, nil)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func TestSlotSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupNoSlashingProtection(seed)
	require.NoError(t, err)

	derivedPriv, err := util.PrivateKeyFromSeedAndPath(seed, "m/12381/3600/0/0/0") // TODO - refactor to remte wealdetch dependency
	require.NoError(t, err)

	accountPriv := &bls.SecretKey{}
	require.NoError(t, accountPriv.SetHexString(hex.EncodeToString(derivedPriv.Marshal())))

	tests := []struct {
		name          string
		slot          uint64
		pubKey        []byte
		domain        []byte
		expectedError error
		accountPriv   *bls.SecretKey
		msg           string
	}{
		{
			name:          "simple sign",
			slot:          1,
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			domain:        []byte("domain"),
			expectedError: nil,
			accountPriv:   accountPriv,
			msg:           "c47e6c550b583a4bce0f2504d81045042d7c4bf439f769e8838f8686a93993f7",
		},
		{
			name:          "unknown account, should error",
			slot:          1,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        []byte("domain"),
			expectedError: errors.New("account not found"),
			accountPriv:   nil,
			msg:           "",
		},
		{
			name:          "nil account, should error",
			slot:          1,
			pubKey:        nil,
			domain:        []byte("domain"),
			expectedError: errors.New("account was not supplied"),
			accountPriv:   nil,
			msg:           "",
		},
		{
			name:          "empty account, should error",
			slot:          1,
			pubKey:        _byteArray(""),
			domain:        []byte("domain"),
			expectedError: errors.New("account was not supplied"),
			accountPriv:   nil,
			msg:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := signer.SignSlot(test.slot, test.domain, test.pubKey)
			if test.expectedError != nil {
				if err != nil {
					require.Equal(t, test.expectedError.Error(), err.Error())
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				require.NoError(t, err)

				sig := &bls.Sign{}
				err := sig.SetHexString(hex.EncodeToString(res))
				require.NoError(t, err)
				require.True(t, sig.Verify(test.accountPriv.GetPublicKey(), test.msg))
			}
		})
	}
}
