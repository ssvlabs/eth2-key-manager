package signer

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"testing"

	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/stretchr/testify/require"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
	prot "github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/bloxapp/eth2-key-manager/stores/inmemory"
	"github.com/bloxapp/eth2-key-manager/wallets"
)

func inmemStorage() *inmemory.InMemStore {
	return inmemory.NewInMemStore(core.MainNetwork)
}

func setupNoSlashingProtection(seed []byte) (ValidatorSigner, error) {
	noProtection := &prot.NoProtection{}
	store := inmemStorage()
	wallet, err := walletWithSeed(seed, store)
	if err != nil {
		return nil, err
	}
	return NewSimpleSigner(wallet, noProtection, core.PyrmontNetwork), nil
}

func setupWithSlashingProtection(seed []byte, setLatestAttestation bool, setLatestProposal bool) (ValidatorSigner, error) {
	store := inmemStorage()
	protector := prot.NewNormalProtection(store)
	wallet, err := walletWithSeed(seed, store)
	if err != nil {
		return nil, err
	}

	// update highest attestation
	acc, err := wallet.AccountByPublicKey("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")
	if err != nil {
		log.Fatal(err)
	}
	if setLatestAttestation {
		protector.UpdateHighestAttestation(acc.ValidatorPublicKey(), &eth.AttestationData{
			Slot:            0,
			CommitteeIndex:  0,
			BeaconBlockRoot: ignoreError(hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000")).([]byte),
			Source: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000")).([]byte),
			},
			Target: &eth.Checkpoint{
				Epoch: 0,
				Root:  ignoreError(hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000")).([]byte),
			},
		})
	}

	if setLatestProposal {
		protector.UpdateHighestProposal(acc.ValidatorPublicKey(), &eth.BeaconBlock{
			Slot: 0,
		})
	}

	return NewSimpleSigner(wallet, protector, core.PyrmontNetwork), nil
}

func walletWithSeed(seed []byte, store core.Storage) (core.Wallet, error) {
	if err := core.InitBLS(); err != nil { // very important!
		return nil, err
	}

	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	vault, err := eth2keymanager.NewKeyVault(options)
	if err != nil {
		return nil, err
	}

	wallet, err := vault.Wallet()
	if err != nil {
		return nil, err
	}

	_, err = wallet.CreateValidatorAccount(seed, nil)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func TestSlotSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupNoSlashingProtection(seed)
	require.NoError(t, err)

	pk := &bls.PublicKey{}
	require.NoError(t, pk.Deserialize(_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")))

	tests := []struct {
		name          string
		slot          uint64
		pubKey        []byte
		domain        []byte
		expectedError error
		msg           string
	}{
		{
			name:          "simple sign",
			slot:          1,
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: nil,
			msg:           "7920f65abe2efb506d0ec763e227ab58978b6e2dda41d4bc2ceb785b4084b0fa",
		},
		{
			name:          "unknown account, should error",
			slot:          1,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			msg:           "",
		},
		{
			name:          "nil account, should error",
			slot:          1,
			pubKey:        nil,
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account was not supplied"),
			msg:           "",
		},
		{
			name:          "empty account, should error",
			slot:          1,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			msg:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := signer.SignSlot(test.slot, test.domain, test.pubKey)
			if test.expectedError != nil {
				if err != nil {
					require.Equal(t, test.expectedError.Error(), err.Error())
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				require.NoError(t, err)

				sig := &bls.Sign{}
				err := sig.Deserialize(res)
				require.NoError(t, err)
				require.True(t, sig.VerifyByte(pk, _byteArray(test.msg)))
			}
		})
	}
}

// tested against a block and sig generated from https://github.com/prysmaticlabs/prysm/blob/develop/beacon-chain/sync/validate_aggregate_proof_test.go#L300
func TestAggregateProofReferenceSignatures(t *testing.T) {
	sk := _byteArray("6327b1e58c41d60dd7c3c8b9634204255707c2d12e2513c345001d8926745eea")
	pk := _byteArray("954eb88ed1207f891dc3c28fa6cfdf8f53bf0ed3d838f3476c0900a61314d22d4f0a300da3cd010444dd5183e35a593c")
	domain := _byteArray32("050000008c84cda94176cc2b1268357c57c3160131874a4408e155b0db826d11")
	slot := uint64(0)
	sigByts := _byteArray("a1167cdbebeae876b3fa82d4f4c35fc3dc4706c7ae20cee359919fdbc93a2588c3f7a15c80d12a20c78ac6381a9fe35d06f6b8ae7e95fb87fa2195511bd53ce6f385aa71dda52b38771f954348a57acad9dde225da614c50c02173314417b096")

	// setup KeyVault
	store := inmemStorage()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	options.SetWalletType(core.NDWallet)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)
	wallet, err := vault.Wallet()
	require.NoError(t, err)
	k, err := core.NewHDKeyFromPrivateKey(sk, "")
	require.NoError(t, err)
	acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
	require.NoError(t, wallet.AddValidatorAccount(acc))

	// setup signer
	protector := prot.NewNormalProtection(store)
	signer := NewSimpleSigner(wallet, protector, core.PyrmontNetwork)

	sig, err := signer.SignSlot(slot, domain, pk)
	require.NoError(t, err)
	require.EqualValues(t, sigByts, sig)
}

// tested against a block and sig generated from https://github.com/prysmaticlabs/prysm/blob/develop/beacon-chain/sync/validate_aggregate_proof_test.go#L300
func TestAggregateAndProofReferenceSignatures(t *testing.T) {
	sk := _byteArray("6327b1e58c41d60dd7c3c8b9634204255707c2d12e2513c345001d8926745eea")
	pk := _byteArray("954eb88ed1207f891dc3c28fa6cfdf8f53bf0ed3d838f3476c0900a61314d22d4f0a300da3cd010444dd5183e35a593c")
	domain := _byteArray32("060000008c84cda94176cc2b1268357c57c3160131874a4408e155b0db826d11")
	aggAttByts := _byteArray("08161260a1167cdbebeae876b3fa82d4f4c35fc3dc4706c7ae20cee359919fdbc93a2588c3f7a15c80d12a20c78ac6381a9fe35d06f6b8ae7e95fb87fa2195511bd53ce6f385aa71dda52b38771f954348a57acad9dde225da614c50c02173314417b0961ad1010a0109126a1a20eade62f0457b2fdf48e7d3fc4b60736688286be7c7a3ac4c9a16a5e0600bd9e42222122068656c6c6f2d776f726c640000000000000000000000000000000000000000002a221220eade62f0457b2fdf48e7d3fc4b60736688286be7c7a3ac4c9a16a5e0600bd9e41a60b101ab9cd396472716e5334ecbaf797078452117d73596bc5893480ae48f94eee6d5d7dfd67dad69969771f73b75c10816ce412a385cb85cb556d23649d5587cfc7758d95ee5b0ad33ae1a23ecad7fc08a86eba222497d7ed123a46b893393cd")
	sigByts := _byteArray("8bf29e58a5b594415ce220c3a9f0d64a4cfa44397f92138f8f31849100149e18e0418ed0cb6068f38909b01e9950d7360a8ba1504bd7451c74add42acd82b148ac0b5f3687c429cc571b96307a8902e9976a24747ad68ad21e372302236aab25")

	// setup KeyVault
	store := inmemStorage()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	options.SetWalletType(core.NDWallet)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)
	wallet, err := vault.Wallet()
	require.NoError(t, err)
	k, err := core.NewHDKeyFromPrivateKey(sk, "")
	require.NoError(t, err)
	acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
	require.NoError(t, wallet.AddValidatorAccount(acc))

	// setup signer
	protector := prot.NewNormalProtection(store)
	signer := NewSimpleSigner(wallet, protector, core.PyrmontNetwork)

	// decode aggregated att proof
	aggAndProof := &eth.AggregateAttestationAndProof{}
	require.NoError(t, aggAndProof.Unmarshal(aggAttByts))

	sig, err := signer.SignAggregateAndProof(aggAndProof, domain, pk)
	require.NoError(t, err)
	require.EqualValues(t, sigByts, sig)
}

// tested against a block and sig generated from  https://github.com/prysmaticlabs/prysm/blob/develop/shared/testutil/block.go#L170
func TestRandaoReferenceSignatures(t *testing.T) {
	sk := _byteArray("5470813f7deef638dc531188ca89e36976d536f680e89849cd9077fd096e20bc")
	pk := _byteArray("a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972")
	domain := _byteArray32("0200000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459")
	sigByts := _byteArray("a2c156a4bc9439f1d85f922f2abaa96e830f1c526101211bdb7d16f4ad9490a0302fc5adb089c05b5f16fd465962f47c04fc2b81a94d135a07c1613db61511c17284b51fafab984e56d3411e16e45f5068f146d9412f91d31ab0f237eac3d745")

	// setup KeyVault
	store := inmemStorage()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	options.SetWalletType(core.NDWallet)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)
	wallet, err := vault.Wallet()
	require.NoError(t, err)
	k, err := core.NewHDKeyFromPrivateKey(sk, "")
	require.NoError(t, err)
	acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
	require.NoError(t, wallet.AddValidatorAccount(acc))

	// setup signer
	protector := prot.NewNormalProtection(store)
	signer := NewSimpleSigner(wallet, protector, core.PyrmontNetwork)

	// decode epoch
	epoch := binary.LittleEndian.Uint64(_byteArray("0000000000000000000000000000000000000000000000000000000000000000"))

	sig, err := signer.SignEpoch(epoch, domain, pk)
	require.NoError(t, err)
	require.EqualValues(t, sigByts, sig)
}

func TestEpochSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupNoSlashingProtection(seed)
	require.NoError(t, err)

	pk := &bls.PublicKey{}
	require.NoError(t, pk.Deserialize(_byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")))

	tests := []struct {
		name          string
		epoch         uint64
		pubKey        []byte
		domain        []byte
		expectedError error
		msg           string
	}{
		{
			name:          "simple sign",
			epoch:         1,
			pubKey:        _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: nil,
			msg:           "7920f65abe2efb506d0ec763e227ab58978b6e2dda41d4bc2ceb785b4084b0fa",
		},
		{
			name:          "unknown account, should error",
			epoch:         1,
			pubKey:        _byteArray("83e04069ed28b637f113d272a235af3e610401f252860ed2063d87d985931229458e3786e9b331cd73d9fc58863d9e4c"),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			msg:           "",
		},
		{
			name:          "nil account, should error",
			epoch:         1,
			pubKey:        nil,
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account was not supplied"),
			msg:           "",
		},
		{
			name:          "empty account, should error",
			epoch:         1,
			pubKey:        _byteArray(""),
			domain:        _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"),
			expectedError: errors.New("account not found"),
			msg:           "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := signer.SignEpoch(test.epoch, test.domain, test.pubKey)
			if test.expectedError != nil {
				if err != nil {
					require.Equal(t, test.expectedError.Error(), err.Error())
				} else {
					t.Errorf("no error returned, expected: %s", test.expectedError.Error())
				}
			} else {
				// check sign worked
				require.NoError(t, err)

				sig := &bls.Sign{}
				err := sig.Deserialize(res)
				require.NoError(t, err)
				require.True(t, sig.VerifyByte(pk, _byteArray(test.msg)))
			}
		})
	}
}
