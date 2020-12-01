package signer

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func (signer *SimpleSigner) ListAccounts() (*pb.ListAccountsResponse, error) {
	var ret []*pb.Account
	for _, account := range signer.wallet.Accounts() {
		ret = append(ret, &pb.Account{
			Name:      account.Name(),
			PublicKey: account.ValidatorPublicKey(),
		})
	}

	return &pb.ListAccountsResponse{
		State:    pb.ResponseState_SUCCEEDED,
		Accounts: ret,
	}, nil
}
