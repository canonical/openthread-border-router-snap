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

	/*
		Becoming the Leader:
		The node successfully transitions from having no role to establishing the network.
		[N] Mle-----------: Role detached -> leader
	*/
	waitForLogMessage(t, otbrSnap, "Role detached -> leader", start)

	/*
		Forming a Partition:
		The network creates a unique Partition ID, meaning a Thread partition is actively running.
		[N] Mle-----------: Partition ID 0x6f4f1b4c
	*/
	waitForLogMessage(t, otbrSnap, "Partition ID 0x", start)
	/*
		Populating the Router Table:
		The node registers itself correctly in the routing table as the leader.
		[I] RouterTable---:      2 0x0800 - me - leader
	*/
	waitForLogMessage(t, otbrSnap, " - me - leader", start)

	state, _, _ := utils.Exec(t, "sudo openthread-border-router.ot-ctl state | head -n 1")
	state = strings.TrimRight(state, "\n")
	require.Equal(t, "leader", state)
}
