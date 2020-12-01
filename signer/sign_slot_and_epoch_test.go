package signer

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/herumi/bls-eth-go-binary/bls"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

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
	if err := core.InitBLS(); err != nil { // very important!
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

	pk := &bls.PublicKey{}
	require.NoError(t, pk.Deserialize(_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")))

	tests := []struct {
		name          string
		slot          uint64
		pubKey        []byte
		domain        []byte
		expectedError error
		msg           string
	}{
		{
			name:          "simple sign",
			slot:          1,
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: nil,
			msg:           "7920f65abe2efb506d0ec763e227ab58978b6e2dda41d4bc2ceb785b4084b0fa",
		},
		{
			name:          "unknown account, should error",
			slot:          1,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			msg:           "",
		},
		{
			name:          "nil account, should error",
			slot:          1,
			pubKey:        nil,
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account was not supplied"),
			msg:           "",
		},
		{
			name:          "empty account, should error",
			slot:          1,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
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
				err := sig.Deserialize(res)
				require.NoError(t, err)
				require.True(t, sig.VerifyByte(pk, _byteArray(test.msg)))
			}
		})
	}
}

func TestEpochSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupNoSlashingProtection(seed)
	require.NoError(t, err)

	pk := &bls.PublicKey{}
	require.NoError(t, pk.Deserialize(_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")))

	tests := []struct {
		name          string
		epoch         uint64
		pubKey        []byte
		domain        []byte
		expectedError error
		msg           string
	}{
		{
			name:          "simple sign",
			epoch:         1,
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: nil,
			msg:           "7920f65abe2efb506d0ec763e227ab58978b6e2dda41d4bc2ceb785b4084b0fa",
		},
		{
			name:          "unknown account, should error",
			epoch:         1,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			msg:           "",
		},
		{
			name:          "nil account, should error",
			epoch:         1,
			pubKey:        nil,
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account was not supplied"),
			msg:           "",
		},
		{
			name:          "empty account, should error",
			epoch:         1,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			msg:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := signer.SignEpoch(test.epoch, test.domain, test.pubKey)
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
				err := sig.Deserialize(res)
				require.NoError(t, err)
				require.True(t, sig.VerifyByte(pk, _byteArray(test.msg)))
			}
		})
	}
}
