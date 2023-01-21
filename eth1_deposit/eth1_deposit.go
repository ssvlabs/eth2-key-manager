package eth1deposit

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssvtypes "github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
	util "github.com/wealdtech/go-eth2-util"

	"github.com/bloxapp/eth2-key-manager/core"
)

const (
	// MaxEffectiveBalanceInGwei is the max effective balance
	MaxEffectiveBalanceInGwei phase0.Gwei = 32000000000

	// BLSWithdrawalPrefixByte is the BLS withdrawal prefix
	BLSWithdrawalPrefixByte = byte(0)
)

// IsSupportedDepositNetwork returns true if the given network is supported
var IsSupportedDepositNetwork = func(network core.Network) bool {
	return network == core.PyrmontNetwork || network == core.PraterNetwork || network == core.MainNetwork
}

// DepositData is basically copied from https://github.com/prysmaticlabs/prysm/blob/master/shared/keystore/deposit_input.go
func DepositData(validationKey *core.HDKey, withdrawalPubKey []byte, network core.Network, amount phase0.Gwei) (*phase0.DepositData, [32]byte, error) {
	if !IsSupportedDepositNetwork(network) {
		return nil, [32]byte{}, errors.Errorf("Network %s is not supported", network)
	}

	depositData := &phase0.DepositData{
		WithdrawalCredentials: withdrawalCredentialsHash(withdrawalPubKey),
		Amount:                amount,
	}
	copy(depositData.PublicKey[:], validationKey.PublicKey().Serialize())

	objRoot, err := depositData.HashTreeRoot()
	if err != nil {
		return nil, [32]byte{}, errors.Wrap(err, "failed to determine the root hash of deposit data")
	}

	// Create domain
	domain, err := ssvtypes.ComputeETHDomain(ssvtypes.DomainDeposit, network.ForkVersion(), ssvtypes.GenesisValidatorsRoot)
	if err != nil {
		return nil, [32]byte{}, errors.Wrap(err, "failed to calculate domain")
	}

	signingData := phase0.SigningData{
		ObjectRoot: objRoot,
		Domain:     domain,
	}

	root, err := signingData.HashTreeRoot()
	if err != nil {
		return nil, [32]byte{}, errors.Wrap(err, "failed to determine the root hash of signing container")
	}

	// Sign
	sig, err := validationKey.Sign(root[:])
	if err != nil {
		return nil, [32]byte{}, errors.Wrap(err, "failed to sign the root")
	}

	copy(depositData.Signature[:], sig)
	depositDataRoot, err := depositData.HashTreeRoot()
	if err != nil {
		return nil, [32]byte{}, errors.Wrap(err, "failed to determine the root hash of deposit data")
	}

	return depositData, depositDataRoot, nil
}

// withdrawalCredentialsHash forms a 32 byte hash of the withdrawal public
// address.
//
// The specification is as follows:
//
//	withdrawal_credentials[:1] == BLS_WITHDRAWAL_PREFIX_BYTE
//	withdrawal_credentials[1:] == hash(withdrawal_pubkey)[1:]
//
// where withdrawal_credentials is of type bytes32.
func withdrawalCredentialsHash(withdrawalPubKey []byte) []byte {
	h := util.SHA256(withdrawalPubKey)
	return append([]byte{BLSWithdrawalPrefixByte}, h[1:]...)[:32]
}
