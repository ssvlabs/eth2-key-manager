package ValidatorSigner

import (
	"fmt"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

// copy from prysm
type Checkpoint struct {
	Epoch uint64
	Root  []byte `ssz-size:"32"`
}

// copy from prysm
type BeaconAttestation struct {
	Slot            uint64
	CommitteeIndex  uint64
	BeaconBlockRoot []byte `ssz-size:"32"`
	Source          *Checkpoint
	Target          *Checkpoint
}

func (signer *SimpleSigner) SignBeaconAttestation(req *pb.SignBeaconAttestationRequest) (*pb.SignResponse, error) {
	// 1. check we can even sign this
	if val,err := signer.slashingProtector.IsSlashableAttestation(req); err != nil || !val {
		if err != nil {
			return nil,err
		}
		return nil, fmt.Errorf("slashable attestation, not signing")
	}

	// 2. add to protection storage
	err := signer.slashingProtector.SaveAttestation(req)
	if err != nil {
		return nil, err
	}

	// 3. get the account
	if req.GetAccount() == "" { // TODO by public key
		return nil, fmt.Errorf("account was not supplied")
	}
	account,error := signer.wallet.AccountByName(req.GetAccount())
	if error != nil {
		return nil,error
	}

	// 4. lock for current account
	signer.lock(account.ID(), "attestation")
	defer func () {
		signer.unlockAndDelete(account.ID(), "attestation")
	}()

	// 5.
	forSig,err := prepareAttestationReqForSigning(req)
	if err != nil {
		return nil, err
	}
	sig,err := account.Sign(forSig)
	if err != nil {
		return nil, err
	}
	res := &pb.SignResponse{
		State:                pb.ResponseState_SUCCEEDED,
		Signature:            sig.Marshal(),
	}

	return res,nil
}

func prepareAttestationReqForSigning(req *pb.SignBeaconAttestationRequest) ([]byte,error) {
	data := &BeaconAttestation{ // Create a local copy of the data; we need ssz size information to calculate the correct root.
		Slot:            req.Data.Slot,
		CommitteeIndex:  req.Data.CommitteeIndex,
		BeaconBlockRoot: req.Data.BeaconBlockRoot,
		Source: &Checkpoint{
			Epoch: req.Data.Source.Epoch,
			Root:  req.Data.Source.Root,
		},
		Target: &Checkpoint{
			Epoch: req.Data.Target.Epoch,
			Root:  req.Data.Target.Root,
		},
	}
	forSig,err := prepareForSig(data, req.Domain)
	if err != nil {
		return nil, err
	}
	return forSig[:],nil
}