package eth1_deposit

import (
	"github.com/bloxapp/KeyVault/core"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-ssz"
	types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
)

const (
	MaxEffectiveBalanceInGwei uint64 = 32000000000
	BLSWithdrawalPrefixByte byte = byte(0)
)

// this is basically copied from https://github.com/prysmaticlabs/prysm/blob/master/shared/keystore/deposit_input.go
func DepositData(validationAccount core.Account, withdrawalAccount core.Account, amount uint64) (*ethpb.Deposit_Data, [32]byte, error) {
	di := &ethpb.Deposit_Data{
		PublicKey:             validationAccount.PublicKey().Marshal(),
		WithdrawalCredentials: withdrawalCredentialsHash(withdrawalAccount),
		Amount:                amount,
	}

	sr, err := ssz.SigningRoot(di)
	if err != nil {
		return nil, [32]byte{}, err
	}

	domain := types.Domain(types.DomainDeposit,nil /*forkVersion*/, nil /*genesisValidatorsRoot*/)

	// prepare for sig
	signingContainer := struct {
		Root   []byte `json:"object_root,omitempty" ssz-size:"32"`
		Domain []byte `json:"domain,omitempty" ssz-size:"32"`
	}{
		Root:sr[:],
		Domain:domain,
	}
	root,err := ssz.HashTreeRoot(signingContainer)
	if err != nil {
		return nil, [32]byte{}, err
	}

	// sign
	sig,err := validationAccount.Sign(root[:])
	if err != nil {
		return nil, [32]byte{}, err
	}
	di.Signature = sig.Marshal()

	// root with sig
	dr,err := ssz.HashTreeRoot(di)
	if err != nil {
		return nil, [32]byte{}, err
	}

	return di,dr,nil
}

// withdrawalCredentialsHash forms a 32 byte hash of the withdrawal public
// address.
//
// The specification is as follows:
//   withdrawal_credentials[:1] == BLS_WITHDRAWAL_PREFIX_BYTE
//   withdrawal_credentials[1:] == hash(withdrawal_pubkey)[1:]
// where withdrawal_credentials is of type bytes32.
func withdrawalCredentialsHash(withdrawalAccount core.Account) []byte {
	h := util.SHA256(withdrawalAccount.PublicKey().Marshal())
	return append([]byte{BLSWithdrawalPrefixByte}, h[1:]...)[:32]
}