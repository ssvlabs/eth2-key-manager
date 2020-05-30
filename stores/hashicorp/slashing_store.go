package hashicorp

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

const (
	WalletSlashingProtectionPathStr = WalletIdsBaseePath + "%s/protection"
)

func (store *HashicorpVaultStore) SaveAttestation(req *pb.SignBeaconAttestationRequest) {

}

func (store *HashicorpVaultStore) RetrieveAttestation(epoch uint64, slot uint64) (*pb.SignBeaconAttestationRequest, error) {

}

func (store *HashicorpVaultStore) SaveProposal(req *pb.SignBeaconProposalRequest) {

}

func (store *HashicorpVaultStore) RetrieveProposal(epoch uint64, slot uint64) (*pb.SignBeaconProposalRequest, error) {

}

