package signer

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/wrapper"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/timeutils"
	"github.com/stretchr/testify/require"

	eth2keymanager "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
	prot "github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/bloxapp/eth2-key-manager/wallets"
)

func testBlock(t *testing.T) *eth.BeaconBlock {
	blockByts := "7b22736c6f74223a312c2270726f706f7365725f696e646578223a38352c22706172656e745f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c2273746174655f726f6f74223a227264584c666d704c2f396a4f662b6c7065753152466d4747486a4571315562633955674257576d505236553d222c22626f6479223a7b2272616e64616f5f72657665616c223a226f734657704c79554f664859583549764b727170626f4d5048464a684153456232333057394b32556b4b41774c38577473496e41573138572f555a5a597652384250777267616c4e45316f48775745397468555277584b4574522b767135684f56744e424868626b5831426f3855625a51532b5230787177386a667177396446222c22657468315f64617461223a7b226465706f7369745f726f6f74223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d222c226465706f7369745f636f756e74223a3132382c22626c6f636b5f68617368223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d227d2c226772616666697469223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d222c2270726f706f7365725f736c617368696e6773223a6e756c6c2c2261747465737465725f736c617368696e6773223a6e756c6c2c226174746573746174696f6e73223a5b7b226167677265676174696f6e5f62697473223a2248773d3d222c2264617461223a7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c22736f75726365223a7b2265706f6368223a302c22726f6f74223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d227d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d227d7d2c227369676e6174757265223a226c37627963617732537751633147587a4c36662f6f5a39616752386562685278503550675a546676676e30344b367879384a6b4c68506738326276674269675641674347767965357a7446797a4772646936555a655a4850593030595a6d3964513939764352674d34676f31666b3046736e684543654d68522f45454b59626a227d5d2c226465706f73697473223a6e756c6c2c22766f6c756e746172795f6578697473223a6e756c6c7d7d"
	blk := &eth.BeaconBlock{}
	require.NoError(t, json.Unmarshal(_byteArray(blockByts), blk))
	return blk
}

// tested against a block and sig generated from https://github.com/prysmaticlabs/prysm/blob/hf1/shared/testutil/altair.go#L313-L438
func TestBenchmarkBlockProposalAltair(t *testing.T) {
	require.NoError(t, core.InitBLS())

	// fixture
	sk := "2799ceccbdaf1e36679b413193a363bfe6d2d35c8cf6ff6151165707461eaed7"
	pk := "b245d63d3f9d8ea1807a629fcb1b328cb4d542f35a3d5bc478be0df389dddd712fc4c816ba3fede9a96320ae6b24a7d8"
	domain := "000000004535ad2cc8e6d4f94cde8707ab5ef9a7d23f884df84c2531832a2b5e"
	blockSSZByts := "01000000000000001c00000000000000df7140ad4f8e394cab798d89fa7612a284de78aa004f6db9387a4269a4a0669c83387dd0abb441a3c16886c8144098cb4cac5e363516f329c368550094fd7ff754000000b1e2f27dfac80e4f1bce84adf11acf6cdbb0d8e59a575c9795020e614eb3aa29634108c0559c04ce02b93fc9a5a8daf60485ebac039864c79d51bef54915aa8c45cbcde3215f14962be196a6b8648851c35b4a804ce8d5fb6c5ff49800ef7740685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb732000000000000000685b310217aa8ed11d15f0ac1df629f3dc95b9e0e8fc550025cb18ae36f8fb7300000000000000000000000000000000000000000000000000000000000000007c0100007c0100007c0100006502000065020000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa2ec291dd5e91096ae48b3659a7ac59567a48c030bb6ac9435d6d44ef39f3f664742f35b38cd6e41ade9ed417183cc0c0b407dfea8627ccc2275fc82ab3d2182e58a037eb144811d741d18894698396efde2b7873c2db9b712e03dfcd03705ef04000000e400000000000000000000000000000000000000df7140ad4f8e394cab798d89fa7612a284de78aa004f6db9387a4269a4a0669c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000df7140ad4f8e394cab798d89fa7612a284de78aa004f6db9387a4269a4a0669cb62ce3f28e8731dce73d5761fdc5e30383d42a022d6e939974d0586d82270f79b38b86d17237e4241a761e239c594e7a0d4ef731470001be3b125ba515f8f215f9309a9ba12653bf9d704a4125865b9775c8a65223e3ca027781175200a2d24403"
	sigByts := "ae37c5b026490df77d4757ad278f1a5f2c46c8971eaa9b62ad09534d3626b1d1da1371edeb061ed6f2f327f79d3bf79101941225dfb97af6d15c8fd6ac7ae8a180392a19aecec50c938b1ea6faf8ba848a4dc9a01f46502e98caec4a01b06575"

	// setup KeyVault
	store := inmemStorage()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	options.SetWalletType(core.NDWallet)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)
	wallet, err := vault.Wallet()
	require.NoError(t, err)
	k, err := core.NewHDKeyFromPrivateKey(_byteArray(sk), "")
	require.NoError(t, err)
	acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
	require.NoError(t, wallet.AddValidatorAccount(acc))

	// setup signer
	signer := NewSimpleSigner(wallet, &prot.NoProtection{}, core.PraterNetwork)

	// decode block
	blk := &eth.BeaconBlockAltair{}
	require.NoError(t, blk.UnmarshalSSZ(_byteArray(blockSSZByts)))

	altairBeaconInterface, err := wrapper.WrappedAltairBeaconBlock(blk)
	require.NoError(t, err)
	sig, err := signer.SignBeaconBlock(altairBeaconInterface, _byteArray(domain), _byteArray(pk))
	require.NoError(t, err)
	require.EqualValues(t, _byteArray(sigByts), sig)
}

// tested against a block and sig generated from https://github.com/prysmaticlabs/prysm/blob/master/shared/testutil/block.go#L86
func TestBenchmarkBlockProposal(t *testing.T) {
	require.NoError(t, core.InitBLS())

	// fixture
	sk := "5470813f7deef638dc531188ca89e36976d536f680e89849cd9077fd096e20bc"
	pk := "a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972"
	domain := "0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"
	blockByts := "7b22736c6f74223a312c2270726f706f7365725f696e646578223a38352c22706172656e745f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c2273746174655f726f6f74223a227264584c666d704c2f396a4f662b6c7065753152466d4747486a4571315562633955674257576d505236553d222c22626f6479223a7b2272616e64616f5f72657665616c223a226f734657704c79554f664859583549764b727170626f4d5048464a684153456232333057394b32556b4b41774c38577473496e41573138572f555a5a597652384250777267616c4e45316f48775745397468555277584b4574522b767135684f56744e424868626b5831426f3855625a51532b5230787177386a667177396446222c22657468315f64617461223a7b226465706f7369745f726f6f74223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d222c226465706f7369745f636f756e74223a3132382c22626c6f636b5f68617368223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d227d2c226772616666697469223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d222c2270726f706f7365725f736c617368696e6773223a6e756c6c2c2261747465737465725f736c617368696e6773223a6e756c6c2c226174746573746174696f6e73223a5b7b226167677265676174696f6e5f62697473223a2248773d3d222c2264617461223a7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c22736f75726365223a7b2265706f6368223a302c22726f6f74223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d227d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d227d7d2c227369676e6174757265223a226c37627963617732537751633147587a4c36662f6f5a39616752386562685278503550675a546676676e30344b367879384a6b4c68506738326276674269675641674347767965357a7446797a4772646936555a655a4850593030595a6d3964513939764352674d34676f31666b3046736e684543654d68522f45454b59626a227d5d2c226465706f73697473223a6e756c6c2c22766f6c756e746172795f6578697473223a6e756c6c7d7d"
	sigByts := "911ac2f6d74039279f16eee4cc46f4c6eea0ef9d18f0d9739b407c150c07ccb104c1c4b034ad46b25719bafc22fad05205975393000ea09636f5ce427814e2fe12ea72041099cc7f6ec249e504992dbf65e968ab448ddf4e124cbcbc722829b5"

	// setup KeyVault
	store := inmemStorage()
	options := &eth2keymanager.KeyVaultOptions{}
	options.SetStorage(store)
	options.SetWalletType(core.NDWallet)
	vault, err := eth2keymanager.NewKeyVault(options)
	require.NoError(t, err)
	wallet, err := vault.Wallet()
	require.NoError(t, err)
	k, err := core.NewHDKeyFromPrivateKey(_byteArray(sk), "")
	require.NoError(t, err)
	acc := wallets.NewValidatorAccount("1", k, nil, "", vault.Context)
	require.NoError(t, wallet.AddValidatorAccount(acc))

	// setup signer
	signer := NewSimpleSigner(wallet, &prot.NoProtection{}, core.PraterNetwork)

	// decode block
	blk := &eth.BeaconBlock{}
	require.NoError(t, json.Unmarshal(_byteArray(blockByts), blk))

	sig, err := signer.SignBeaconBlock(wrapper.WrappedPhase0BeaconBlock(blk), _byteArray(domain), _byteArray(pk))
	require.NoError(t, err)
	require.EqualValues(t, _byteArray(sigByts), sig)
}

func TestProposalSlashingSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupWithSlashingProtection(nil, seed, true, true)
	require.NoError(t, err)

	t.Run("valid proposal", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99

		_, err = signer.SignBeaconBlock(wrapper.WrappedPhase0BeaconBlock(blk), _byteArray("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)
	})

	t.Run("valid proposal, sign using nil pk. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		_, err = signer.SignBeaconBlock(wrapper.WrappedPhase0BeaconBlock(blk), _byteArray("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), nil)
		require.NotNil(t, err)
		require.Error(t, err, "account was not supplied")
	})

	t.Run("double proposal, different state root. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		blk.StateRoot = _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459")
		_, err = signer.SignBeaconBlock(wrapper.WrappedPhase0BeaconBlock(blk), _byteArray("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
	})

	t.Run("double proposal, different body root. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		blk.Body.Graffiti = []byte("different body root")
		_, err = signer.SignBeaconBlock(wrapper.WrappedPhase0BeaconBlock(blk), []byte("domain"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
	})

	t.Run("double proposal, different parent root. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		blk.ParentRoot = _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52458")
		_, err = signer.SignBeaconBlock(wrapper.WrappedPhase0BeaconBlock(blk), []byte("domain"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
	})

	t.Run("double proposal, different proposer index. Should error", func(t *testing.T) {
		blk := testBlock(t)
		blk.Slot = 99
		blk.ProposerIndex = 3
		_, err = signer.SignBeaconBlock(wrapper.WrappedPhase0BeaconBlock(blk), []byte("domain"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NotNil(t, err)
		require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
	})
}

func TestFarFutureProposalSignature(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	network := core.PraterNetwork
	maxValidSlot := network.EstimatedSlotAtTime(timeutils.Now().Unix() + FarFutureMaxValidEpoch)

	t.Run("max valid source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(nil, seed, true, true)
		require.NoError(t, err)

		blk := testBlock(t)
		blk.Slot = maxValidSlot

		_, err = signer.SignBeaconBlock(wrapper.WrappedPhase0BeaconBlock(blk), _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.NoError(t, err)
	})
	t.Run("too far into the future source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(nil, seed, true, true)
		require.NoError(t, err)

		blk := testBlock(t)
		blk.Slot = maxValidSlot + 1

		_, err = signer.SignBeaconBlock(wrapper.WrappedPhase0BeaconBlock(blk), _byteArray32("0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"), _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf"))
		require.EqualError(t, err, "proposed block slot too far into the future")
	})
}
