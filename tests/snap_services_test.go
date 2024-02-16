package tests

import (
	"testing"

	"github.com/canonical/matter-snap-testing/utils"
	"github.com/stretchr/testify/require"
)

func TestSnapServicesStatus(t *testing.T) {
	// Start clean
	utils.SnapStop(t, otbrSnap)

	t.Cleanup(func() {
		utils.SnapStop(t, otbrSnap)
	})

	utils.SnapStart(nil, otbrSnap)

	// Oneshot service
	require.False(t, utils.SnapServicesActive(t, otbrSetupApp))

	// Active services
	require.True(t, utils.SnapServicesActive(t, otbrWebApp))
	require.True(t, utils.SnapServicesActive(t, otbrAgentApp))
}
