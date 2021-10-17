package core

import (
	"testing"

	prysmTime "github.com/prysmaticlabs/prysm/time"

	types "github.com/prysmaticlabs/eth2-types"
	"github.com/stretchr/testify/require"
)

func TestNetworkMainnet(t *testing.T) {
	net := NetworkFromString(string(MainNetwork))

	secondsPassedSinceGenesis := prysmTime.Now().Unix() - 1606824023
	require.EqualValues(t, types.Epoch(secondsPassedSinceGenesis/(12*32)), net.EstimatedCurrentEpoch())
	require.EqualValues(t, types.Epoch(secondsPassedSinceGenesis/12), net.EstimatedCurrentSlot())
	require.EqualValues(t, types.Epoch(secondsPassedSinceGenesis/12), net.EstimatedSlotAtTime(prysmTime.Now().Unix()))
	require.EqualValues(t, types.Epoch(101010/32), net.EstimatedEpochAtSlot(types.Slot(101010)))
}
