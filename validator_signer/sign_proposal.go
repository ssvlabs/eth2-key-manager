package validator_signer

import (
	"fmt"

	"github.com/bloxapp/KeyVault/core"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func (signer *SimpleSigner) SignBeaconProposal(req *pb.SignBeaconProposalRequest) (*pb.SignResponse, error) {
	// 1. get the account
	if req.GetAccount() == "" { // TODO by public key
		return nil, fmt.Errorf("account was not supplied")
	}
	account, err := signer.wallet.AccountByName(req.GetAccount())
	if err != nil {
		return nil, err
	}

	// 2. lock for current account
	signer.lock(account.ID(), "proposal")
	defer func() {
		signer.unlockAndDelete(account.ID(), "proposal")
	}()

	// 3. check we can even sign this
	if val, err := signer.slashingProtector.IsSlashableProposal(account, req); err != nil || val != nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("slashable proposal, not signing")
	}

	// 4. add to protection storage
	if err := signer.slashingProtector.SaveProposal(account, req); err != nil {
		return nil, err
	}

	// 5. generate ssz root hash and sign
	forSig, err := PrepareProposalReqForSigning(req)
	if err != nil {
		return nil, err
	}
	sig, err := account.Sign(forSig)
	if err != nil {
		return nil, err
	}
	res := &pb.SignResponse{
		State:     pb.ResponseState_SUCCEEDED,
		Signature: sig.Marshal(),
	}

	return res, nil
}

// PrepareProposalReqForSigning prepares the given proposal request for signing.
// This is exported to allow use it by custom signing mechanism.
func PrepareProposalReqForSigning(req *pb.SignBeaconProposalRequest) ([]byte, error) {
	data := core.ToCoreBlockData(req)
	forSig, err := prepareForSig(data, req.Domain)
	if err != nil {
		return nil, err
	}
	return forSig[:], nil
}
