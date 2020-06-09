package validator_signer

import (
	"fmt"
	"github.com/prysmaticlabs/go-ssz"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func (signer *SimpleSigner) Sign(req *pb.SignRequest) (*pb.SignResponse, error) {
	// 1. check we can even sign this
	// TODO - should we?

	// 2. get the account
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

	// 4.
	forSig,err := prepareReqForSigning(req)
	if err != nil {
		return nil, err
	}
	sig,err := account.Sign(forSig[:])
	if err != nil {
		return nil, err
	}
	res := &pb.SignResponse{
		State:                pb.ResponseState_SUCCEEDED,
		Signature:            sig.Marshal(),
	}

	return res,nil
}

func prepareReqForSigning(req *pb.SignRequest) ([32]byte,error) {
	signingData := struct {
		Hash   []byte `ssz-size:"32"`
		Domain []byte `ssz-size:"32"`
	}{
		Hash:   req.Data,
		Domain: req.Domain,
	}
	return ssz.HashTreeRoot(signingData)
}