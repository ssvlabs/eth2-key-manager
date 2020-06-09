package validator_signer

import (
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func (signer *SimpleSigner) SignBeaconAttestation(req *pb.SignBeaconAttestationRequest) (*pb.SignResponse, error) {
	// 1. get the account
	if req.GetAccount() == "" { // TODO by public key
		return nil, fmt.Errorf("account was not supplied")
	}
	account,error := signer.wallet.AccountByName(req.GetAccount())
	if error != nil {
		return nil,error
	}
	if account == nil {
		return nil,fmt.Errorf("account not found")
	}

	// 2. lock for current account
	signer.lock(account.ID(), "attestation")
	defer func () {
		signer.unlockAndDelete(account.ID(), "attestation")
	}()

	// 3. check we can even sign this
	if val,err := signer.slashingProtector.IsSlashableAttestation(account,req); err != nil || len(val) != 0 {
		if err != nil {
			return nil,err
		}
		return nil, fmt.Errorf("slashable attestation, not signing")
	}

	// 4. add to protection storage
	err := signer.slashingProtector.SaveAttestation(account,req)
	if err != nil {
		return nil, err
	}

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
	data := core.ToCoreAttestationData(req)
	forSig,err := prepareForSig(data, req.Domain)
	if err != nil {
		return nil, err
	}
	return forSig[:],nil
}