package signer

import (
	"errors"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// hashRootOnly implements only HashTreeRoot() — the single method the signer's
// HashRoot interface requires — and deliberately omits fastssz's GetTree(). It
// pins the narrowed contract: if HashRoot is ever widened back to ssz.HashRoot
// (or gains another method), the compile-time assertion below fails here rather
// than silently breaking dynamic-ssz-generated Gloas types downstream.
type hashRootOnly struct {
	root [32]byte
}

func (h hashRootOnly) HashTreeRoot() ([32]byte, error) { return h.root, nil }

// Compile-time guard: a type exposing only HashTreeRoot() must satisfy HashRoot.
var _ HashRoot = hashRootOnly{}

// hashRootErr satisfies HashRoot but always fails, to exercise the error path.
type hashRootErr struct{}

func (hashRootErr) HashTreeRoot() ([32]byte, error) {
	return [32]byte{}, errors.New("hash tree root failed")
}

// TestComputeETHSigningRoot_NarrowContract verifies that a value implementing
// only HashRoot (not fastssz's ssz.HashRoot) both satisfies the interface and
// is signed correctly end-to-end.
func TestComputeETHSigningRoot_NarrowContract(t *testing.T) {
	var objRoot [32]byte
	for i := range objRoot {
		objRoot[i] = byte(i + 1)
	}
	var domain phase0.Domain
	for i := range domain {
		domain[i] = byte(0xf0 - i)
	}

	got, err := ComputeETHSigningRoot(hashRootOnly{root: objRoot}, domain)
	require.NoError(t, err)

	// The object's HashTreeRoot() output must flow into the signing container,
	// so recomputing the signing root independently must match.
	want, err := (&phase0.SigningData{ObjectRoot: objRoot, Domain: domain}).HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, phase0.Root(want), got)

	// A different object root must yield a different signing root, i.e. the
	// object is actually mixed in rather than ignored.
	other, err := ComputeETHSigningRoot(hashRootOnly{root: [32]byte{0xbb}}, domain)
	require.NoError(t, err)
	require.NotEqual(t, got, other)
}

// TestComputeETHSigningRoot_PropagatesHashTreeRootError verifies that a failure
// from the object's HashTreeRoot() is surfaced rather than swallowed.
func TestComputeETHSigningRoot_PropagatesHashTreeRootError(t *testing.T) {
	_, err := ComputeETHSigningRoot(hashRootErr{}, phase0.Domain{})
	require.Error(t, err)
}
