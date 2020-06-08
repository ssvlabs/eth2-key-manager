package slashing_protection

import (
	"encoding/hex"
	"github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/bloxapp/KeyVault/wallet_hd"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"reflect"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func store () *in_memory.InMemStore {
	return in_memory.NewInMemStore(
		reflect.TypeOf(&KeyVault.KeyVault{}),
		reflect.TypeOf(&wallet_hd.HDWallet{}),
		reflect.TypeOf(&wallet_hd.HDAccount{}),
	)
}

func vault() (*KeyVault.KeyVault,error) {
	options := &KeyVault.PortfolioOptions{}
	options.SetStorage(store())
	options.SetSeed(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"))
	return KeyVault.NewKeyVault(options)
}

func setupProposal() (core.VaultSlashingProtector, []core.Account,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,nil,err
	}

	// create an account to use
	vault,err := vault()
	if err != nil {
		return nil,nil,err
	}
	w,err := vault.CreateWallet("test")
	if err != nil {
		return nil,nil,err
	}
	account1,err := w.CreateValidatorAccount("1")
	if err != nil {
		return nil,nil,err
	}
	account2,err := w.CreateValidatorAccount("2")
	if err != nil {
		return nil,nil,err
	}

	protector := core.NewNormalProtection(vault.Context.Storage.(core.SlashingStore))
	protector.SaveProposal(account1, &pb.SignBeaconProposalRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.BeaconBlockHeader{
			Slot:                 100,
			ProposerIndex:        2,
			ParentRoot:           []byte("A"),
			StateRoot:            []byte("A"),
			BodyRoot:             []byte("A"),
		},
	})
	protector.SaveProposal(account1, &pb.SignBeaconProposalRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.BeaconBlockHeader{
			Slot:                 101,
			ProposerIndex:        2,
			ParentRoot:           []byte("B"),
			StateRoot:            []byte("B"),
			BodyRoot:             []byte("B"),
		},
	})
	protector.SaveProposal(account1, &pb.SignBeaconProposalRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.BeaconBlockHeader{
			Slot:                 102,
			ProposerIndex:        2,
			ParentRoot:           []byte("C"),
			StateRoot:            []byte("C"),
			BodyRoot:             []byte("C"),
		},
	})

	return protector,[]core.Account{account1,account2},nil
}


func TestDoubleProposal(t *testing.T) {
	protector,accounts,err := setupProposal()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("New proposal, should not slash",func(t *testing.T) {
		res, err := protector.IsSlashableProposal(accounts[0], &pb.SignBeaconProposalRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.BeaconBlockHeader{
				Slot:                 99,
				ProposerIndex:        2,
				ParentRoot:           []byte("Z"),
				StateRoot:            []byte("Z"),
				BodyRoot:             []byte("Z"),
			},
		})

		if err != nil {
			t.Error(err)
			return
		}
		if res != nil {
			t.Errorf("found too many/few slashed proposals: %d, expected: %d", 1,0)
			return
		}
	})

	t.Run("different proposer index, should not slash",func(t *testing.T) {
		res, err := protector.IsSlashableProposal(accounts[1], &pb.SignBeaconProposalRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.BeaconBlockHeader{
				Slot:                 100,
				ProposerIndex:        3,
				ParentRoot:           []byte("A"),
				StateRoot:            []byte("A"),
				BodyRoot:             []byte("A"),
			},
		})

		if err != nil {
			t.Error(err)
			return
		}
		if res != nil {
			t.Errorf("found too many/few slashed proposals: %d, expected: %d", 1,0)
			return
		}
	})

	t.Run("double proposal (different body root), should slash",func(t *testing.T) {
		res, err := protector.IsSlashableProposal(accounts[0], &pb.SignBeaconProposalRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.BeaconBlockHeader{
				Slot:                 100,
				ProposerIndex:        2,
				ParentRoot:           []byte("A"),
				StateRoot:            []byte("A"),
				BodyRoot:             []byte("B"),
			},
		})

		if err != nil {
			t.Error(err)
			return
		}
		if res == nil {
			t.Errorf("found too many/few slashed proposals: %d, expected: %d", 0,1)
			return
		}
		if res.Status != core.DoubleProposal {
			t.Errorf("wrong proposal status returned, expected DoubleProposal")
			return
		}
	})
}