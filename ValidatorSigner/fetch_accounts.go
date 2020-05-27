package ValidatorSigner

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func (signer *SimpleSigner) ListAccounts(req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	var ret []*pb.Account
	for account := range signer.wallet.Accounts() {
		ret = append(ret, &pb.Account{
			Name:                 account.Name(),
			PublicKey:            account.PublicKey().Marshal(),
		})
	}

	return &pb.ListAccountsResponse{
		State:                pb.ResponseState_SUCCEEDED,
		Accounts:             ret,
	},nil
}