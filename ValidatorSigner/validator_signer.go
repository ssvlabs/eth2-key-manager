package ValidatorSigner

import (
	"github.com/google/uuid"
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

type signingRoot struct {
	Hash   [32]byte `ssz-size:"32"`
	Domain []byte `ssz-size:"32"`
}

type SimpleSigner struct {
	wallet types.Wallet
	slashingProtector VaultSlashingProtector
	signLocks map[string]*sync.RWMutex
}

// if already locked, will lock until released
func (signer *SimpleSigner) lock (accountId uuid.UUID, operation string) {
	k := accountId.String() + "_" + operation
	signer.signLocks[k] = &sync.RWMutex{}
	signer.signLocks[k].Lock()
}

func (signer *SimpleSigner) unlockAndDelete (accountId uuid.UUID, operation string) {
	k := accountId.String() + "_" + operation
	signer.signLocks[k].Unlock()
	delete(signer.signLocks,k)
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



