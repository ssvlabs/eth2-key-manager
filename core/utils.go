package core

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

// ComputeETHDomain returns computed domain
func ComputeETHDomain(domainType phase0.DomainType, fork phase0.Version, genesisValidatorRoot phase0.Root) (phase0.Domain, error) {
	ret := phase0.Domain{}

	forkData := phase0.ForkData{
		CurrentVersion:        fork,
		GenesisValidatorsRoot: genesisValidatorRoot,
	}
	forkDataRoot, err := forkData.HashTreeRoot()
	if err != nil {
		return ret, errors.Wrap(err, "failed to calculate signature domain")
	}
	copy(ret[:], domainType[:])
	copy(ret[4:], forkDataRoot[:])
	return ret, nil
}

// ComputeETHSigningRoot returns computed root for eth signing
func ComputeETHSigningRoot(obj ssz.HashRoot, domain phase0.Domain) (phase0.Root, error) {
	root, err := obj.HashTreeRoot()
	if err != nil {
		return phase0.Root{}, err
	}
	signingContainer := phase0.SigningData{
		ObjectRoot: root,
		Domain:     domain,
	}
	return signingContainer.HashTreeRoot()
}
