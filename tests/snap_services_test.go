package tests

import (
	"testing"
	"time"

	"github.com/canonical/matter-snap-testing/utils"
	"github.com/stretchr/testify/require"
)

func TestSnapServicesStatus(t *testing.T) {
	// Start clean
	utils.SnapStop(t, otbrSnap)

	t.Cleanup(func() {
		utils.SnapStop(t, otbrSnap)
	})

	start := time.Now()
	utils.SnapStart(nil, otbrSnap)

	// Oneshot service
	waitForLogMessage(t, otbrSetupApp, "OTBR completed oneshot setup", start)
	require.False(t, utils.SnapServicesActive(t, otbrSetupApp))

	// Active services
	waitForLogMessage(t, otbrWebApp, "Border router web started", start)
	require.True(t, utils.SnapServicesActive(t, otbrWebApp))

	// [I] BbrManager----: Start Backbone TMF agent: OK
	waitForLogMessage(t, otbrAgentApp, "Start Backbone TMF agent: OK", start)
	require.True(t, utils.SnapServicesActive(t, otbrAgentApp))
}
