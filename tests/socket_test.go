package tests

import (
	"os"
	"testing"
	"time"

	"github.com/canonical/matter-snap-testing/utils"
)

func TestSocketFile(t *testing.T) {
	// Start clean
	utils.SnapStop(t, otbrSnap)

	t.Cleanup(func() {
		utils.SnapStop(t, otbrSnap)
	})

	utils.SnapStart(nil, otbrSnap)
	waitForFileCreation(t, "/run/snap.openthread-border-router/openthread-wpan0.sock", 10)
}

func waitForFileCreation(t *testing.T, filePath string, maxRetry int) {
	t.Helper()

	for i := 1; i <= maxRetry; i++ {
		time.Sleep(1 * time.Second)
		t.Logf("Retry %d/%d: Waiting for file creation: %s", i, maxRetry, filePath)

		if _, err := os.Stat(filePath); err == nil {
			t.Logf("Socket File exists in: %s", filePath)
			return
		} else if !os.IsNotExist(err) {
			t.Fatalf("Error checking file: %s\n", err)
			continue
		}
	}

	t.Fatalf("Timeout: File not created after %d retries.", maxRetry)
}
