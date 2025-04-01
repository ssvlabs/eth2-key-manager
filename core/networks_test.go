package core

import (
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestNetworkMainnet(t *testing.T) {
	net, err := NetworkFromString(string(MainNetwork))
	require.NoError(t, err)
	require.Equal(t, MainNetwork, net)

	secondsPassedSinceGenesis := time.Now().Unix() - 1606824023
	require.EqualValues(t, phase0.Epoch(secondsPassedSinceGenesis/(12*32)), net.EstimatedCurrentEpoch())
	require.EqualValues(t, phase0.Epoch(secondsPassedSinceGenesis/12), net.EstimatedCurrentSlot())
	require.EqualValues(t, phase0.Epoch(secondsPassedSinceGenesis/12), net.EstimatedSlotAtTime(time.Now().Unix()))
	require.EqualValues(t, phase0.Epoch(101010/32), net.EstimatedEpochAtSlot(phase0.Slot(101010)))
}
