package ValidatorSigner

import (
	"fmt"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	"sync"
)

// copy from prysm https://github.com/prysmaticlabs/prysm/blob/master/validator/client/validator_propose.go#L220-L226
type beaconBlockHeader struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    []byte `ssz-size:"32"`
	StateRoot     []byte `ssz-size:"32"`
	BodyRoot      []byte `ssz-size:"32"`
}

type signingRoot struct {
	Hash   [32]byte `ssz-size:"32"`
	Domain []byte `ssz-size:"32"`
}

func (signer *SimpleSigner) SignBeaconProposal(req *pb.SignBeaconProposalRequest) (*pb.SignResponse, error) {
	// 1. check we can even sign this
	if val,err := signer.slashingProtector.IsSlashablePropose(req); err != nil || !val {
		if err != nil {
			return nil,err
		}
		return nil, fmt.Errorf("slashable proposal, not signing")
	}

	// 2. get the account
	if req.GetAccount() == "" { // TODO by public key
		return nil, fmt.Errorf("account was not supplied")
	}
	account,error := signer.wallet.AccountByName(req.GetAccount())
	if error != nil {
		return nil,error
	}

	// 3. lock for current account
	signer.proposalLocks[account.ID().String()] = &sync.RWMutex{}
	signer.proposalLocks[account.ID().String()].Lock()
	defer func () {
		signer.proposalLocks[account.ID().String()].Unlock()
		delete(signer.proposalLocks,account.ID().String())
	}()

	// 4. generate ssz root hash and sign
	forSig,err := prepareProposalReqForSigning(req)
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
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	// 5. add to protection storage
	err = signer.slashingProtector.SaveProposal(req)
	if err != nil {
		return nil, err
	}

	return res,nil
}

func prepareProposalReqForSigning(req *pb.SignBeaconProposalRequest) ([]byte,error) {
	data := &beaconBlockHeader{ // Create a local copy of the data; we need ssz size information to calculate the correct root.
		Slot:          req.Data.Slot,
		ProposerIndex: req.Data.ProposerIndex,
		ParentRoot:    req.Data.ParentRoot,
		StateRoot:     req.Data.StateRoot,
		BodyRoot:      req.Data.BodyRoot,
	}
	forSig,err := prepareForSig(data, req.Domain)
	if err != nil {
		return nil, err
	}
	return forSig[:],nil
}