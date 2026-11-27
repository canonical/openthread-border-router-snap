package tests

import (
	"context"
	"os/exec"
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

	waitForOTCtl(t)
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

func waitForOTCtl(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var lastOutput []byte
	var lastErr error
	for {
		lastOutput, lastErr = exec.CommandContext(ctx, "sudo", "openthread-border-router.ot-ctl", "state").CombinedOutput()
		if lastErr == nil {
			return
		}

		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for ot-ctl: %v; output: %s", lastErr, lastOutput)
		case <-time.After(time.Second):
		}
	}
}
