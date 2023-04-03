package core

import (
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestNetworkMainnet(t *testing.T) {
	network, err := NetworkFromString(string(MainNetwork))
	require.NoError(t, err)

	secondsPassedSinceGenesis := time.Now().Unix() - 1606824023
	require.EqualValues(t, phase0.Epoch(secondsPassedSinceGenesis/(12*32)), network.EstimatedCurrentEpoch())
	require.EqualValues(t, phase0.Epoch(secondsPassedSinceGenesis/12), network.EstimatedCurrentSlot())
	require.EqualValues(t, phase0.Epoch(secondsPassedSinceGenesis/12), network.EstimatedSlotAtTime(time.Now().Unix()))
	require.EqualValues(t, phase0.Epoch(101010/32), network.EstimatedEpochAtSlot(phase0.Slot(101010)))
}
