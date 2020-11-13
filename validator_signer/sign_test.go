package validator_signer

import (
	"encoding/hex"
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
	seed, _ := hex.DecodeString("f51883a4c56467458c3b47d06cd135f862a6266fabdfb9e9e4702ea5511375d7")
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
				Id:     &pb.SignRequest_PublicKey{PublicKey: _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4b")},
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
