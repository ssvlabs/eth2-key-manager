package eth1_deposit

import (
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	ssz "github.com/prysmaticlabs/go-ssz"
	types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"

	"github.com/bloxapp/eth2-key-manager/core"
)

const (
	MaxEffectiveBalanceInGwei uint64 = 32000000000
	BLSWithdrawalPrefixByte   byte   = byte(0)
)

// DepositData is basically copied from https://github.com/prysmaticlabs/prysm/blob/master/shared/keystore/deposit_input.go
func DepositData(validationKey *core.HDKey, withdrawalPubKey []byte, network core.Network, amountInGwei uint64) (*ethpb.Deposit_Data, [32]byte, error) {
	depositData := struct {
		PublicKey             []byte `ssz-size:"48"`
		WithdrawalCredentials []byte `ssz-size:"32"`
		Amount                uint64
	}{
		PublicKey:             validationKey.PublicKey().Serialize(),
		WithdrawalCredentials: withdrawalCredentialsHash(withdrawalPubKey),
		Amount:                amountInGwei,
	}
	objRoot, err := ssz.HashTreeRoot(depositData)
	if err != nil {
		return nil, [32]byte{}, errors.Wrap(err, "failed to determine the root hash of deposit data")
	}

	// Create domain
	domain := types.Domain(types.DomainDeposit, network.ForkVersion(), types.ZeroGenesisValidatorsRoot)

	// Prepare for sig
	signingContainer := struct {
		Root   []byte `json:"object_root,omitempty" ssz-size:"32"`
		Domain []byte `json:"domain,omitempty" ssz-size:"32"`
	}{
		Root:   objRoot[:],
		Domain: domain,
	}
	root, err := ssz.HashTreeRoot(signingContainer)
	if err != nil {
		return nil, [32]byte{}, errors.Wrap(err, "failed to determine the root hash of signing container")
	}

	// Sign
	sig, err := validationKey.Sign(root[:])
	if err != nil {
		return nil, [32]byte{}, errors.Wrap(err, "failed to sign the root")
	}

	signedDepositData := &ethpb.Deposit_Data{
		PublicKey:             validationKey.PublicKey().Serialize(),
		WithdrawalCredentials: withdrawalCredentialsHash(withdrawalPubKey),
		Amount:                amountInGwei,
		Signature:             sig,
	}

	depositDataRoot, err := ssz.HashTreeRoot(signedDepositData)
	if err != nil {
		return nil, [32]byte{}, errors.Wrap(err, "failed to determine the root hash of deposit data")
	}

	return signedDepositData, depositDataRoot, nil
}

// withdrawalCredentialsHash forms a 32 byte hash of the withdrawal public
// address.
//
// The specification is as follows:
//   withdrawal_credentials[:1] == BLS_WITHDRAWAL_PREFIX_BYTE
//   withdrawal_credentials[1:] == hash(withdrawal_pubkey)[1:]
// where withdrawal_credentials is of type bytes32.
func withdrawalCredentialsHash(withdrawalPubKey []byte) []byte {
	h := util.SHA256(withdrawalPubKey)
	return append([]byte{BLSWithdrawalPrefixByte}, h[1:]...)[:32]
}
