package validator_signer

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-ssz"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func (signer *SimpleSigner) Sign(req *pb.SignRequest) (*pb.SignResponse, error) {
	// 1. check we can even sign this
	// TODO - should we?

	// 2. get the account
	if req.GetPublicKey() == nil {
		return nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(req.GetPublicKey()))
	if err != nil {
		return nil, err
	}

	forSig, err := PrepareReqForSigning(req)
	if err != nil {
		return nil, err
	}

	sig, err := account.ValidationKeySign(forSig[:])
	if err != nil {
		return nil, err
	}

	return &pb.SignResponse{
		State:     pb.ResponseState_SUCCEEDED,
		Signature: sig.Marshal(),
	}, nil
}

// PrepareReqForSigning prepares sign request for signing.
func PrepareReqForSigning(req *pb.SignRequest) ([32]byte, error) {
	signingData := struct {
		Hash   []byte `ssz-size:"32"`
		Domain []byte `ssz-size:"32"`
	}{
		Hash:   req.Data,
		Domain: req.Domain,
	}
	return ssz.HashTreeRoot(signingData)
}