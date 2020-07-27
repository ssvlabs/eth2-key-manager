package validator_signer

import (
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	prot "github.com/bloxapp/KeyVault/slashing_protection"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
	"testing"
)

func inmemStorage() *in_memory.InMemStore {
	return in_memory.NewInMemStore()
}

func setupNoSlashingProtection(seed []byte) (ValidatorSigner,error) {
	noProtection := &prot.NoProtection{}
	store := inmemStorage()
	wallet,err := walletWithSeed(seed,store)
	if err != nil {
		return nil,err
	}
	return NewSimpleSigner(wallet,noProtection),nil
}

func setupWithSlashingProtection(seed []byte) (ValidatorSigner,error) {
	store := inmemStorage()
	protector := prot.NewNormalProtection(store)
	wallet,err := walletWithSeed(seed,store)
	if err != nil {
		return nil,err
	}
	return NewSimpleSigner(wallet,protector),nil
}

func walletWithSeed(seed []byte, store core.Storage) (core.Wallet,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,err
	}


	options := &KeyVault.WalletOptions{}
	options.SetStorage(store)
	options.SetSeed(seed)
	vault,err := KeyVault.ImportKeyVault(options)
	if err != nil {
		return nil,err
	}

	wallet,err := vault.Wallet()
	if err != nil {
		return nil,err
	}

	_,err = wallet.CreateValidatorAccount("1")
	if err != nil {
		return nil,err
	}

	return wallet,nil
}

func TestSignatures(t *testing.T) {
	seed,_ := hex.DecodeString("f51883a4c56467458c3b47d06cd135f862a6266fabdfb9e9e4702ea5511375d7")
	signer,err := setupNoSlashingProtection(seed)
	if err != nil {
		t.Error(err)
		return
	}
	accountPriv,err := util.PrivateKeyFromSeedAndPath(seed,"m/12381/3600/0/0/0")
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name string
		req *pb.SignRequest
		expectedError error
		accountPriv *e2types.BLSPrivateKey
		msg string
	}{
		{
			name:"simple sign",
			req: &pb.SignRequest{
				Id:                   &pb.SignRequest_Account{Account:"1"},
				Data:                 []byte("data"),
				Domain:               []byte("domain"),
			},
			expectedError: nil,
			accountPriv: accountPriv,
			msg: "c47e6c550b583a4bce0f2504d81045042d7c4bf439f769e8838f8686a93993f7",
		},
		{
			name:"unknown account, should error",
			req: &pb.SignRequest{
				Id:                   &pb.SignRequest_Account{Account:"10"},
				Data:                 []byte("data"),
				Domain:               []byte("domain"),
			},
			expectedError: fmt.Errorf("account not found"),
			accountPriv: nil,
			msg: "",
		},
		{
			name:"empty account, should error",
			req: &pb.SignRequest{
				Id:                   &pb.SignRequest_Account{Account:""},
				Data:                 []byte("data"),
				Domain:               []byte("domain"),
			},
			expectedError: fmt.Errorf("account was not supplied"),
			accountPriv: nil,
			msg: "",
		},
		{
			name:"nil account, should error",
			req: &pb.SignRequest{
				Id:                  nil,
				Data:                 []byte("data"),
				Domain:               []byte("domain"),
			},
			expectedError: fmt.Errorf("account was not supplied"),
			accountPriv: nil,
			msg: "",
		},
	}

	for _,test := range tests {
		t.Run(test.name,func(t *testing.T) {
			res,err := signer.Sign(test.req)
			if test.expectedError != nil {
				if err != nil {
					if err.Error() != test.expectedError.Error() {
						t.Errorf("wrong error returned: %s, expected: %s", err.Error(),test.expectedError.Error())
					}
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				if err != nil {
					t.Error(err)
					return
				}

				sig,err := e2types.BLSSignatureFromBytes(res.Signature)
				if err != nil {
					t.Error(err)
					return
				}
				msgBytes,err := hex.DecodeString(test.msg)
				if err != nil {
					t.Error(err)
					return
				}
				if !sig.Verify(msgBytes,test.accountPriv.PublicKey()) {
					t.Errorf("signature does not verify against pubkey",)
				}
			}
		})
	}
}


func TestWalletLocked(t *testing.T) {

}

func TestAccountLocked(t *testing.T) {

}
