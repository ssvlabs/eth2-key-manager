package signer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

// SignBeaconBlock signs the given beacon block
func (signer *SimpleSigner) SignBeaconBlock(b *spec.VersionedBeaconBlock, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	slot, err := b.Slot()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get block slot")
	}

	var block ssz.HashRoot
	switch b.Version {
	case spec.DataVersionPhase0:
		block = b.Phase0
	case spec.DataVersionAltair:
		block = b.Altair
	case spec.DataVersionBellatrix:
		block = b.Bellatrix
	case spec.DataVersionCapella:
		block = b.Capella
	case spec.DataVersionDeneb:
		block = b.Deneb
	case spec.DataVersionElectra:
		block = b.Electra

	default:
		return nil, nil, errors.Errorf("unsupported block version %d", b.Version)
	}

	return signer.SignBlock(block, slot, domain, pubKey)
}
