package tests

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/canonical/matter-snap-testing/utils"
	"github.com/stretchr/testify/require"
)

var start = time.Now()

const (
	otbrSnap     = "openthread-border-router"
	otbrSetupApp = "openthread-border-router.otbr-setup"
	otbrAgentApp = "openthread-border-router.otbr-agent"
	otbrWebApp   = "openthread-border-router.otbr-web"
)

func TestMain(m *testing.M) {
	teardown, err := setup()
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestDeviceOperations(t *testing.T) {
	t.Cleanup(func() {
		// TODO
	})

	t.Run("Setup", func(t *testing.T) {
		//Firewall
		//IP forwarding
		//RT tables
		//NAT44
		//Border routing
	})

	t.Run("Socket File", func(t *testing.T) {
		waitForFileCreation(t, "/run/snap.openthread-border-router/openthread-wpan0.sock", 10*time.Second)
	})

	t.Run("Snap services status", func(t *testing.T) {
		// oneshot service
		require.False(t, utils.SnapServicesActive(t, otbrSetupApp))

		// actice services
		require.True(t, utils.SnapServicesActive(t, otbrWebApp))
		require.True(t, utils.SnapServicesActive(t, otbrAgentApp))
	})

	t.Run("Thread Network Formation", func(t *testing.T) {
		utils.Exec(t, "sudo openthread-border-router.ot-ctl dataset init new")
		utils.Exec(t, "sudo openthread-border-router.ot-ctl dataset commit active")
		utils.Exec(t, "sudo openthread-border-router.ot-ctl ifconfig up")
		utils.Exec(t, "sudo openthread-border-router.ot-ctl thread start")

		utils.WaitForLogMessage(t, otbrSnap, "Thread Network", start)

		state, _, _ := utils.Exec(t, "sudo openthread-border-router.ot-ctl state | head -n 1")
		state = strings.TrimRight(state, "\n")
		require.Equal(t, "leader", state)
	})
}

func setup() (teardown func(), err error) {

	log.Println("[CLEAN]")
	utils.SnapRemove(nil, otbrSnap)

	log.Println("[SETUP]")
	start := time.Now()

	teardown = func() {
		log.Println("[TEARDOWN]")
		utils.SnapDumpLogs(nil, start, otbrSnap)

		log.Println("Removing installed snap:", !utils.SkipTeardownRemoval)
		if !utils.SkipTeardownRemoval {
			utils.SnapRemove(nil, otbrSnap)
		}
	}

	if utils.LocalServiceSnap() {
		err = utils.SnapInstallFromFile(nil, utils.LocalServiceSnapPath)
	} else {
		err = utils.SnapInstallFromStore(nil, otbrSnap, utils.ServiceChannel)
	}
	if err != nil {
		teardown()
		return
	}

	// connect interfaces
	utils.SnapConnect(nil, otbrSnap+":avahi-control", "")
	utils.SnapConnect(nil, otbrSnap+":firewall-control", "")
	utils.SnapConnect(nil, otbrSnap+":raw-usb", "")
	utils.SnapConnect(nil, otbrSnap+":network-control", "")
	utils.SnapConnect(nil, otbrSnap+":bluetooth-control", "")
	utils.SnapConnect(nil, otbrSnap+":bluez", "")

	// copy and set simulated RCP
	utils.Exec(nil, "sudo cp ot-rcp-simulator-thread-reference-20230119-amd64 /var/snap/openthread-border-router/common/")
	utils.SnapSet(nil, otbrSnap, "radio-url", "'spinel+hdlc+forkpty:///var/snap/openthread-border-router/common/ot-rcp-simulator-thread-reference-20230119-amd64?forkpty-arg=1'")

	utils.SnapStart(nil, otbrSnap)

	return
}

func waitForFileCreation(t *testing.T, filePath string, timeout time.Duration) {
	t.Helper()

	const maxRetry = 10

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

func readLogFile(t *testing.T, filePath string) (string, error) {
	t.Helper()

	text, err := os.ReadFile(filePath)

	if err != nil {
		return "", err
	}
	return string(text), nil
}
