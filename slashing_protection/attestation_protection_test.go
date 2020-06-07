package slashing_protection

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	"testing"
)

func setupAttestation() (core.VaultSlashingProtector, []core.Account,error) {
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
	protector.SaveAttestation(account1, &pb.SignBeaconAttestationRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.AttestationData{
			Slot:                 30,
			CommitteeIndex:       5,
			BeaconBlockRoot:      []byte("A"),
			Source:               &pb.Checkpoint{
				Epoch:                1,
				Root:                 []byte("B"),
			},
			Target:               &pb.Checkpoint{
				Epoch:                2,
				Root:                 []byte("C"),
			},
		},
	})
	protector.SaveAttestation(account1, &pb.SignBeaconAttestationRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.AttestationData{
			Slot:                 30,
			CommitteeIndex:       5,
			BeaconBlockRoot:      []byte("A"),
			Source:               &pb.Checkpoint{
				Epoch:                2,
				Root:                 []byte("B"),
			},
			Target:               &pb.Checkpoint{
				Epoch:                3,
				Root:                 []byte("C"),
			},
		},
	})
	protector.SaveAttestation(account1, &pb.SignBeaconAttestationRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.AttestationData{
			Slot:                 30,
			CommitteeIndex:       5,
			BeaconBlockRoot:      []byte("B"),
			Source:               &pb.Checkpoint{
				Epoch:                3,
				Root:                 []byte("C"),
			},
			Target:               &pb.Checkpoint{
				Epoch:                4,
				Root:                 []byte("D"),
			},
		},
	})
	protector.SaveAttestation(account1, &pb.SignBeaconAttestationRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.AttestationData{
			Slot:                 30,
			CommitteeIndex:       5,
			BeaconBlockRoot:      []byte("B"),
			Source:               &pb.Checkpoint{
				Epoch:                4,
				Root:                 []byte("C"),
			},
			Target:               &pb.Checkpoint{
				Epoch:                10,
				Root:                 []byte("D"),
			},
		},
	})
	protector.SaveAttestation(account1, &pb.SignBeaconAttestationRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.AttestationData{
			Slot:                 30,
			CommitteeIndex:       5,
			BeaconBlockRoot:      []byte("B"),
			Source:               &pb.Checkpoint{
				Epoch:                5,
				Root:                 []byte("C"),
			},
			Target:               &pb.Checkpoint{
				Epoch:                9,
				Root:                 []byte("D"),
			},
		},
	})
	return protector, []core.Account{account1,account2},nil
}

func TestSurroundingVote(t *testing.T) {
	protector,accounts,err := setupAttestation()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("1 Surrounded vote",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0], &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       4,
				BeaconBlockRoot:      []byte("A"),
				Source:               &pb.Checkpoint{
					Epoch:                2,
					Root:                 []byte("B"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                5,
					Root:                 []byte("C"),
				},
			},
		})

		if err != nil {
			t.Error(err)
			return
		}
		if len(res) != 1 {
			t.Errorf("found too many/few slashed attestations: %d, expected: %d", len(res),1)
			return
		}
		if res[0].Status != core.SurroundedVote {
			t.Errorf("wrong attestation status returned, expected SurroundingVote")
			return
		}
	})

	t.Run("2 Surrounded votes",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0], &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       4,
				BeaconBlockRoot:      []byte("A"),
				Source:               &pb.Checkpoint{
					Epoch:                1,
					Root:                 []byte("B"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                7,
					Root:                 []byte("C"),
				},
			},
		})

		if err != nil {
			t.Error(err)
			return
		}
		if len(res) != 2 {
			t.Errorf("found too many/few slashed attestations: %d, expected: %d", len(res),2)
			return
		}
		if res[0].Status != core.SurroundedVote || res[1].Status != core.SurroundedVote {
			t.Errorf("wrong attestation status returned, expected SurroundingVote")
			return
		}
	})

	t.Run("1 Surrounding vote",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0], &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       4,
				BeaconBlockRoot:      []byte("A"),
				Source:               &pb.Checkpoint{
					Epoch:                5,
					Root:                 []byte("B"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                7,
					Root:                 []byte("C"),
				},
			},
		})
		if err != nil {
			t.Error(err)
			return
		}
		if len(res) != 1 {
			t.Errorf("found too many/few slashed attestations: %d, expected: %d", len(res),1)
			return
		}
		if res[0].Status != core.SurroundingVote {
			t.Errorf("wrong attestation status returned, expected SurroundedVote")
			return
		}
	})

	t.Run("2 Surrounding vote",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0], &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       4,
				BeaconBlockRoot:      []byte("A"),
				Source:               &pb.Checkpoint{
					Epoch:                6,
					Root:                 []byte("B"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                7,
					Root:                 []byte("C"),
				},
			},
		})
		if err != nil {
			t.Error(err)
			return
		}
		if len(res) != 2 {
			t.Errorf("found too many/few slashed attestations: %d, expected: %d", len(res),2)
			return
		}
		if res[0].Status != core.SurroundingVote || res[1].Status != core.SurroundingVote {
			t.Errorf("wrong attestation status returned, expected SurroundedVote")
			return
		}
	})
}

func TestDoubleAttestationVote(t *testing.T) {
	protector,accounts,err := setupAttestation()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Different committee index, should slash",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0], &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       4,
				BeaconBlockRoot:      []byte("A"),
				Source:               &pb.Checkpoint{
					Epoch:                2,
					Root:                 []byte("B"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                3,
					Root:                 []byte("C"),
				},
			},
		})

		if err != nil {
			t.Error(err)
		}
		if len(res) != 1 {
			t.Errorf("found too many/few slashed attestations: %d, expected: %d", len(res),1)
		}
		if res[0].Status != core.DoubleVote {
			t.Errorf("wrong attestation status returned, expected DoubleVote")
		}
	})

	t.Run("Different block root, should slash",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0], &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       5,
				BeaconBlockRoot:      []byte("AA"),
				Source:               &pb.Checkpoint{
					Epoch:                2,
					Root:                 []byte("B"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                3,
					Root:                 []byte("C"),
				},
			},
		})

		if err != nil {
			t.Error(err)
		}
		if len(res) != 1 {
			t.Errorf("found too many/few slashed attestations: %d, expected: %d", len(res),1)
		}
		if res[0].Status != core.DoubleVote {
			t.Errorf("wrong attestation status returned, expected DoubleVote")
		}
	})

	t.Run("Same attestation, should not error",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0], &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       5,
				BeaconBlockRoot:      []byte("B"),
				Source:               &pb.Checkpoint{
					Epoch:                3,
					Root:                 []byte("C"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                4,
					Root:                 []byte("D"),
				},
			},
		})

		if err != nil || len(res) != 0 {
			t.Error(fmt.Errorf("correct attestation found slashable"))
		}
	})

	t.Run("new attestation, should not error",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0], &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       5,
				BeaconBlockRoot:      []byte("E"),
				Source:               &pb.Checkpoint{
					Epoch:                10,
					Root:                 []byte("I"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                11,
					Root:                 []byte("H"),
				},
			},
		})

		if err != nil || len(res) != 0 {
			t.Error(fmt.Errorf("correct attestation found slashable"))
		}
	})
}
