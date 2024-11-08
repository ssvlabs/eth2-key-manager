package handler

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	types "github.com/wealdtech/go-eth2-types/v2"

	eth2keymanager "github.com/ssvlabs/eth2-key-manager"
	rootcmd "github.com/ssvlabs/eth2-key-manager/cli/cmd"
	"github.com/ssvlabs/eth2-key-manager/cli/cmd/wallet/cmd/account/flag"
	"github.com/ssvlabs/eth2-key-manager/core"
	"github.com/ssvlabs/eth2-key-manager/signer"
	"github.com/ssvlabs/eth2-key-manager/stores/inmemory"
)

// VoluntaryExitFlagValues keeps all collected values for seed
type VoluntaryExitFlagValues struct {
	index              int
	seedBytes          []byte
	currentForkVersion phase0.Version
	epoch              int
	validator          *core.ValidatorInfo
	network            core.Network
	responseType       rootcmd.ResponseType
}

// SignRequestEncoded is the encoded sign request
type SignRequestEncoded struct {
	PublicKey       []byte   `json:"public_key,omitempty"`
	SignatureDomain [32]byte `json:"signature_domain,omitempty"`
	Data            []byte
	ObjectType      string
}

// VoluntaryExit creates a new wallet account(s) and prints the storage.
func (h *Account) VoluntaryExit(cmd *cobra.Command, args []string) error {
	err := core.InitBLS()
	if err != nil {
		return errors.Wrap(err, "failed to init BLS")
	}

	voluntaryExitFlags, err := CollectVoluntaryExitFlags(cmd)
	if err != nil {
		return errors.Wrap(err, "failed to collect voluntary exit flags")
	}

	// Initialize store
	store := inmemory.NewInMemStore(voluntaryExitFlags.network)
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)

	// Create new key vault
	_, err = eth2keymanager.NewKeyVault(options)
	if err != nil {
		return errors.Wrap(err, "failed to create key vault")
	}

	wallet, err := store.OpenWallet()
	if err != nil {
		return errors.Wrap(err, "failed to open wallet")
	}

	// Compute domain
	genesisValidatorsRoot := store.Network().GenesisValidatorsRoot()
	domainBytes, err := types.ComputeDomain(types.DomainVoluntaryExit, voluntaryExitFlags.currentForkVersion[:], genesisValidatorsRoot[:])
	if err != nil {
		return errors.Wrap(err, "failed to calculate domain")
	}
	var domain phase0.Domain
	copy(domain[:], domainBytes)

	voluntaryExit := &phase0.VoluntaryExit{
		Epoch:          phase0.Epoch(voluntaryExitFlags.epoch),
		ValidatorIndex: voluntaryExitFlags.validator.Index,
	}

	if voluntaryExitFlags.responseType == rootcmd.ObjectResponseType {
		simpleSigner := signer.NewSimpleSigner(wallet, nil, store.Network())
		acc, err := wallet.CreateValidatorAccount(voluntaryExitFlags.seedBytes, &voluntaryExitFlags.index)
		if err != nil {
			return errors.Wrap(err, "failed to create validator account")
		}

		// validation that the derived validation public key is the same as the one in the validator info
		if !bytes.Equal(acc.ValidatorPublicKey(), voluntaryExitFlags.validator.Pubkey[:]) {
			derivedPubKey := "0x" + hex.EncodeToString(acc.ValidatorPublicKey())
			providedPubKey := voluntaryExitFlags.validator.Pubkey.String()
			return errors.Errorf("derived validator public key: %s, does not match with the provided one: %s", derivedPubKey, providedPubKey)
		}

		signature, _, err := simpleSigner.SignVoluntaryExit(voluntaryExit, domain, acc.ValidatorPublicKey())
		if err != nil {
			return errors.Wrap(err, "failed to sign voluntary exit")
		}

		signedVoluntaryExit := &phase0.SignedVoluntaryExit{
			Message: voluntaryExit,
		}
		copy(signedVoluntaryExit.Signature[:], signature)

		err = h.printer.JSON(signedVoluntaryExit)
		if err != nil {
			return errors.Wrap(err, "failed to print signed voluntary exit JSON")
		}
		return nil
	}

	// Sign request
	marshalSSZ, err := voluntaryExit.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "failed to marshal voluntary exit")
	}

	signRequest := &SignRequestEncoded{
		PublicKey:       voluntaryExitFlags.validator.Pubkey[:],
		SignatureDomain: domain,
		Data:            marshalSSZ,
		ObjectType:      "*models.SignRequestVoluntaryExit",
	}

	byts, err := json.Marshal(signRequest)
	if err != nil {
		return errors.Wrap(err, "failed to marshal sign request")
	}

	h.printer.Text(hex.EncodeToString(byts))

	return nil
}

// CollectVoluntaryExitFlags returns collected flags for seed
func CollectVoluntaryExitFlags(cmd *cobra.Command) (*VoluntaryExitFlagValues, error) {
	voluntaryExitFlagValues := VoluntaryExitFlagValues{}

	// Get network flag value.
	network, err := rootcmd.GetNetworkFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the network flag value")
	}
	voluntaryExitFlagValues.network = network

	// Get response-type flag value.
	responseType, err := rootcmd.GetResponseTypeFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the response type value")
	}
	voluntaryExitFlagValues.responseType = responseType

	if responseType == rootcmd.ObjectResponseType {
		// Get seed flag value.
		seedFlagValue, err := rootcmd.GetSeedFlagValue(cmd)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve the seed flag value")
		}

		if seedFlagValue == "" {
			return nil, errors.New("seed flag is required for object response type")
		}

		// Get seed bytes
		seedBytes, err := hex.DecodeString(seedFlagValue)
		if err != nil {
			return nil, errors.Wrap(err, "failed to HEX decode seed")
		}
		voluntaryExitFlagValues.seedBytes = seedBytes
	}

	// Get index flag value.
	indexFlagValue, err := rootcmd.GetIndexFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the index flag value")
	}
	voluntaryExitFlagValues.index = indexFlagValue

	// Get current fork version flag value.
	currentForkVersionFlagValue, err := flag.GetCurrentForkVersionFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the current fork version flag value")
	}
	voluntaryExitFlagValues.currentForkVersion = currentForkVersionFlagValue

	// Get epoch flag value.
	epochFlagValue, err := flag.GetEpochFlagValue(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the index flag value")
	}
	voluntaryExitFlagValues.epoch = epochFlagValue

	// Get validator info flag value.
	validator, err := flag.GetVoluntaryExitInfoFlagValue(cmd)
	if err != nil {
		return nil, err
	}
	voluntaryExitFlagValues.validator = validator

	return &voluntaryExitFlagValues, nil
}
