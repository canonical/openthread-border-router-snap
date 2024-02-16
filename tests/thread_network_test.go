package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/canonical/matter-snap-testing/utils"
	"github.com/stretchr/testify/require"
)

func TestThreadNetworkFormation(t *testing.T) {
	// Start clean
	utils.SnapStop(t, otbrSnap)

	t.Cleanup(func() {
		utils.Exec(t, "sudo openthread-border-router.ot-ctl thread stop")
		utils.SnapStop(t, otbrSnap)
	})

	start := time.Now()
	utils.SnapStart(nil, otbrSnap)

	utils.Exec(t, "sudo openthread-border-router.ot-ctl dataset init new")
	utils.Exec(t, "sudo openthread-border-router.ot-ctl dataset commit active")
	utils.Exec(t, "sudo openthread-border-router.ot-ctl ifconfig up")
	utils.Exec(t, "sudo openthread-border-router.ot-ctl thread start")

	utils.WaitForLogMessage(t, otbrSnap, "Thread Network", start)

	state, _, _ := utils.Exec(t, "sudo openthread-border-router.ot-ctl state | head -n 1")
	state = strings.TrimRight(state, "\n")
	require.Equal(t, "leader", state)
}
