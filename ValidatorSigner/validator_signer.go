package ValidatorSigner

import (
	"github.com/prysmaticlabs/go-ssz"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"sync"
)

type ValidatorSigner interface {
	ListAccounts(req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error)
	SignBeaconProposal(req *pb.SignBeaconProposalRequest) (*pb.SignResponse, error)
	SignBeaconAttestation(req *pb.SignBeaconAttestationRequest) (*pb.SignResponse, error)
	Sign(req *pb.SignRequest) (*pb.SignResponse, error)
}

type SimpleSigner struct {
	wallet types.Wallet
	slashingProtector VaultSlashingProtector

	proposalLocks map[string]*sync.RWMutex
	attestationLocks map[string]*sync.RWMutex
}

func prepareForSig(data interface{}, domain []byte) ([32]byte,error) {
	root, err := ssz.HashTreeRoot(data)
	if err != nil {
		return [32]byte{}, err
	}
	signingRoot := &signingRoot{
		Hash:   root,
		Domain: domain,
	}
	forsig,err := ssz.HashTreeRoot(signingRoot)
	if err != nil {
		return [32]byte{}, err
	}

	return forsig,nil
}



