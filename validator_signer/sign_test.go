package validator_signer

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
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
	return NewSimpleSigner(wallet, noProtection), nil
}

func setupWithSlashingProtection(seed []byte) (ValidatorSigner, error) {
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
	protector.UpdateLatestAttestation(acc.ValidatorPublicKey(), &pb.SignBeaconAttestationRequest{
		Id:     &pb.SignBeaconAttestationRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
		Domain: ignoreError(hex.DecodeString("01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac")).([]byte),
		Data: &pb.AttestationData{
			Slot:            0,
			CommitteeIndex:  0,
			BeaconBlockRoot: ignoreError(hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000")).([]byte),
			Source: &pb.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000")).([]byte),
			},
			Target: &pb.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000")).([]byte),
			},
		},
	})

	return NewSimpleSigner(wallet, protector), nil
}

func walletWithSeed(seed []byte, store core.Storage) (core.Wallet, error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil, err
	}

	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	options.SetSeed(seed)
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

func TestSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupNoSlashingProtection(seed)
	require.NoError(t, err)

	accountPriv, err := util.PrivateKeyFromSeedAndPath(seed, "m/12381/3600/0/0/0")
	require.NoError(t, err)

	tests := []struct {
		name          string
		req           *pb.SignRequest
		expectedError error
		accountPriv   *e2types.BLSPrivateKey
		msg           string
	}{
		{
			name: "simple sign",
			req: &pb.SignRequest{
				Id:     &pb.SignRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
				Data:   []byte("data"),
				Domain: []byte("domain"),
			},
			expectedError: nil,
			accountPriv:   accountPriv,
			msg:           "c47e6c550b583a4bce0f2504d81045042d7c4bf439f769e8838f8686a93993f7",
		},
		{
			name: "unknown account, should error",
			req: &pb.SignRequest{
				Id:     &pb.SignRequest_PublicKey{PublicKey: _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c")},
				Data:   []byte("data"),
				Domain: []byte("domain"),
			},
			expectedError: errors.New("account not found"),
			accountPriv:   nil,
			msg:           "",
		},
		{
			name: "empty account, should error",
			req: &pb.SignRequest{
				Id:     &pb.SignRequest_Account{Account: ""},
				Data:   []byte("data"),
				Domain: []byte("domain"),
			},
			expectedError: errors.New("account was not supplied"),
			accountPriv:   nil,
			msg:           "",
		},
		{
			name: "nil account, should error",
			req: &pb.SignRequest{
				Id:     nil,
				Data:   []byte("data"),
				Domain: []byte("domain"),
			},
			expectedError: errors.New("account was not supplied"),
			accountPriv:   nil,
			msg:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := signer.Sign(test.req)
			if test.expectedError != nil {
				if err != nil {
					require.Equal(t, test.expectedError.Error(), err.Error())
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				require.NoError(t, err)

				sig, err := e2types.BLSSignatureFromBytes(res.Signature)
				require.NoError(t, err)

				msgBytes, err := hex.DecodeString(test.msg)
				require.NoError(t, err)
				require.True(t, sig.Verify(msgBytes, test.accountPriv.PublicKey()))
			}
		})
	}
}
