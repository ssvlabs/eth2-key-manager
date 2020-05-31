package slashing_protectors

import (
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/encryptors"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	hd "github.com/wealdtech/go-eth2-wallet-hd/v2"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"testing"
)

func setupProposal() (VaultSlashingProtector, types.Account,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,nil,err
	}

	store := in_memory.NewInMemStore()
	wallet,err := hd.CreateWallet("test",[]byte(""), store, encryptors.NewPlainTextEncryptor())
	if err != nil {
		return nil,nil,err
	}
	err = wallet.Unlock([]byte(""))
	if err != nil {
		return nil,nil,err
	}

	account1,err := wallet.CreateAccount("1",[]byte(""))
	if err != nil {
		return nil,nil,err
	}

	protector := NewNormalProtection(store)
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

	return protector,account1,nil
}


func TestDoubleProposal(t *testing.T) {
	protector,account,err := setupProposal()
	account2 := core.NewSimpleAccount()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("New proposal, should not slash",func(t *testing.T) {
		res, err := protector.IsSlashableProposal(account, &pb.SignBeaconProposalRequest{
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
		res, err := protector.IsSlashableProposal(account2, &pb.SignBeaconProposalRequest{
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
		res, err := protector.IsSlashableProposal(account, &pb.SignBeaconProposalRequest{
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