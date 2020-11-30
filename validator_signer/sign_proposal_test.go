package validator_signer

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	bls2 "github.com/prysmaticlabs/prysm/shared/bls"

	e2types "github.com/wealdtech/go-eth2-types/v2"

	eth "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1_gateway"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/prysmaticlabs/prysm/shared/timeutils"

	"github.com/stretchr/testify/require"
	v1 "github.com/wealdtech/eth2-signer-api/pb/v1"
)

// tested against a block and sig generated from https://github.com/prysmaticlabs/prysm/blob/master/shared/testutil/block.go#L86
func TestBenchmarkBlockProposal(t *testing.T) {
	require.NoError(t, e2types.InitBLS())

	sk := "5470813f7deef638dc531188ca89e36976d536f680e89849cd9077fd096e20bc"
	pk := "a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972"
	domain := "0000000081509579e35e84020ad8751eca180b44df470332d3ad17fc6fd52459"
	blockByts := "7b22736c6f74223a312c2270726f706f7365725f696e646578223a38352c22706172656e745f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c2273746174655f726f6f74223a227264584c666d704c2f396a4f662b6c7065753152466d4747486a4571315562633955674257576d505236553d222c22626f6479223a7b2272616e64616f5f72657665616c223a226f734657704c79554f664859583549764b727170626f4d5048464a684153456232333057394b32556b4b41774c38577473496e41573138572f555a5a597652384250777267616c4e45316f48775745397468555277584b4574522b767135684f56744e424868626b5831426f3855625a51532b5230787177386a667177396446222c22657468315f64617461223a7b226465706f7369745f726f6f74223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d222c226465706f7369745f636f756e74223a3132382c22626c6f636b5f68617368223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d227d2c226772616666697469223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d222c2270726f706f7365725f736c617368696e6773223a6e756c6c2c2261747465737465725f736c617368696e6773223a6e756c6c2c226174746573746174696f6e73223a5b7b226167677265676174696f6e5f62697473223a2248773d3d222c2264617461223a7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c22736f75726365223a7b2265706f6368223a302c22726f6f74223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d227d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d227d7d2c227369676e6174757265223a226c37627963617732537751633147587a4c36662f6f5a39616752386562685278503550675a546676676e30344b367879384a6b4c68506738326276674269675641674347767965357a7446797a4772646936555a655a4850593030595a6d3964513939764352674d34676f31666b3046736e684543654d68522f45454b59626a227d5d2c226465706f73697473223a6e756c6c2c22766f6c756e746172795f6578697473223a6e756c6c7d7d"
	blockBodyRoot := "c7f94e8fd9d2bb7713ab5f8d50ae6d6530e61ae5ba0b8d2c585bb1a831ab796c"
	rootForSigByts := "1c5add8f95f400de68a735793b4578c555c98725978e2bc739fdd219c941d3ad"
	sigByts := "911ac2f6d74039279f16eee4cc46f4c6eea0ef9d18f0d9739b407c150c07ccb104c1c4b034ad46b25719bafc22fad05205975393000ea09636f5ce427814e2fe12ea72041099cc7f6ec249e504992dbf65e968ab448ddf4e124cbcbc722829b5"

	block := &eth.BeaconBlock{}
	require.NoError(t, json.Unmarshal(_byteArray(blockByts), block))

	req := &v1.SignBeaconProposalRequest{
		Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("a3862121db5914d7272b0b705e6e3c5336b79e316735661873566245207329c30f9a33d4fb5f5857fc6fd0a368186972")},
		Domain: _byteArray(domain),
		Data: &v1.BeaconBlockHeader{
			Slot:          block.Slot,
			ProposerIndex: block.ProposerIndex,
			ParentRoot:    block.ParentRoot,
			StateRoot:     block.StateRoot,
			BodyRoot:      _byteArray(blockBodyRoot),
		},
	}

	rootForSig, err := PrepareProposalReqForSigning(req)
	require.NoError(t, err)
	require.EqualValues(t, _byteArray(rootForSigByts), rootForSig)

	priv, err := bls2.SecretKeyFromBytes(_byteArray(sk))
	require.NoError(t, err)
	pub, err := bls2.PublicKeyFromBytes(_byteArray(pk))
	require.NoError(t, err)

	sig := priv.Sign(rootForSig)
	require.EqualValues(t, _byteArray(sigByts), sig.Marshal())

	sig.Verify(pub, rootForSig)
}

func TestProposalSlashingSignatures(t *testing.T) {
	seed, _ := hex.DecodeString("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	signer, err := setupWithSlashingProtection(seed, true)
	require.NoError(t, err)

	t.Run("valid proposal", func(t *testing.T) {
		_, err = signer.SignBeaconBlock(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NoError(t, err)
	})

	t.Run("valid proposal, sign using account name. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconBlock(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_Account{Account: "1"},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.Error(t, err, "account was not supplied")
	})

	t.Run("double proposal, different state root. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconBlock(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("A"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})

	t.Run("double proposal, different body root. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconBlock(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("A"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})

	t.Run("double proposal, different parent root. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconBlock(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 2,
				ParentRoot:    []byte("A"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})

	t.Run("double proposal, different proposer index. Should error", func(t *testing.T) {
		_, err = signer.SignBeaconBlock(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          99,
				ProposerIndex: 3,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NotNil(t, err)
		require.EqualError(t, err, "err, slashable proposal: DoubleProposal")
	})
}

func TestFarFutureProposalSignature(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	network := core.PyrmontNetwork
	maxValidSlot := network.EstimatedSlotAtTime(timeutils.Now().Unix() + FarFutureMaxValidEpoch)

	t.Run("max valid source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)
		_, err = signer.SignBeaconBlock(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          maxValidSlot,
				ProposerIndex: 3,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.NoError(t, err)
	})
	t.Run("too far into the future source", func(tt *testing.T) {
		signer, err := setupWithSlashingProtection(seed, true)
		require.NoError(t, err)
		_, err = signer.SignBeaconBlock(&v1.SignBeaconProposalRequest{
			Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: _byteArray("95087182937f6982ae99f9b06bd116f463f414513032e33a3d175d9662eddf162101fcf6ca2a9fedaded74b8047c5dcf")},
			Domain: []byte("domain"),
			Data: &v1.BeaconBlockHeader{
				Slot:          maxValidSlot + 1,
				ProposerIndex: 3,
				ParentRoot:    []byte("Z"),
				StateRoot:     []byte("Z"),
				BodyRoot:      []byte("Z"),
			},
		})
		require.EqualError(t, err, "proposed block slot too far into the future")
	})
}
