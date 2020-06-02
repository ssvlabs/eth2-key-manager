package validator_signer

import (
	"bytes"
	"encoding/hex"
	"github.com/bloxapp/KeyVault/encryptors"
	prot "github.com/bloxapp/KeyVault/slashing_protectors"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	hd "github.com/wealdtech/go-eth2-wallet-hd/v2"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"testing"
)

func setupNoSlashingProtection(seed []byte) (ValidatorSigner,error) {
	noProtection := &prot.NoProtection{}
	store := in_memory.NewInMemStore()
	wallet,err := walletWithSeed(seed,store)
	if err != nil {
		return nil,err
	}
	return NewSimpleSigner(wallet,noProtection),nil
}

func walletWithSeed(seed []byte, store types.Store) (types.Wallet,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,err
	}

	wallet, err := hd.CreateWalletFromSeed("test",[]byte(""),store,encryptors.NewPlainTextEncryptor(),seed)
	if err != nil {
		return nil,err
	}
	err = wallet.Unlock([]byte(""))
	if err != nil {
		return nil,err
	}
	account,err := wallet.CreateAccount("1",[]byte(""))
	if err != nil {
		return nil,err
	}
	err = account.Unlock([]byte(""))
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

	tests := []struct {
		name string
		req *pb.SignRequest
		expectedSig []byte
		expectedError error
	}{
		{
			name:"simple sign",
			req: &pb.SignRequest{
				Id:                   &pb.SignRequest_Account{Account:"1"},
				Data:                 []byte("data"),
				Domain:               []byte("domain"),
			},
			expectedSig: []byte(""),
			expectedError: nil,
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
			} else if bytes.Compare(res.Signature,test.expectedSig) != 0 {
				t.Errorf("returned signature different that expectd.")
			}

		})
	}
}


func TestWalletLocked(t *testing.T) {

}

func TestAccountLocked(t *testing.T) {

}