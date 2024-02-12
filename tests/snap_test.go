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

func TestSetUp(t *testing.T) {
	// TestSetUp verifies the configuration setup of OTBR based on otbr-setup.sh script:
	// https://github.com/canonical/openthread-border-router-snap/blob/main/snap/local/stage/bin/otbr-setup.sh

	t.Run("Check Setup", func(t *testing.T) {
		INFRA_IF, _, _ := utils.Exec(t, "sudo snap get openthread-border-router infra-if")
		t.Run("firewall", func(t *testing.T) {
			forwardRules, _, _ := utils.Exec(t, "sudo iptables -S FORWARD | grep \"comment OTBR\"")
			require.NotEmpty(t, forwardRules)
		})

		t.Run("IP forwarding", func(t *testing.T) {
			ipv6_forwarding, _, _ := utils.Exec(t, "sudo sysctl net.ipv6.conf.all.forwarding")
			require.Equal(t, "net.ipv6.conf.all.forwarding = 1\n", ipv6_forwarding)
			ipv4_forwarding, _, _ := utils.Exec(t, "sudo sysctl net.ipv4.ip_forward")
			require.Equal(t, "net.ipv4.ip_forward = 1\n", ipv4_forwarding)
		})

		t.Run("RT tables for backbone router", func(t *testing.T) {
			socket_buffer_size, _, _ := utils.Exec(t, "sudo sysctl net.core.optmem_max")
			require.Equal(t, "net.core.optmem_max = 65536\n", socket_buffer_size)
		})

		t.Run("random fwmark bits", func(t *testing.T) {
			mangle_table_prerouting_chain, _, _ := utils.Exec(t, "sudo iptables -t mangle -L PREROUTING -n -v | grep OTBR")
			require.NotEmpty(t, mangle_table_prerouting_chain)
			nat_table_postrouting_chain, _, _ := utils.Exec(t, "sudo iptables -t nat -L POSTROUTING -n -v | grep OTBR")
			require.NotEmpty(t, nat_table_postrouting_chain)
		})

		t.Run("firewall rule setup for INFRA_IF", func(t *testing.T) {
			forward_rule, _, _ := utils.Exec(t, "sudo iptables -t filter -L FORWARD -n -v | grep OTBR | grep "+INFRA_IF)
			require.NotEmpty(t, forward_rule)
		})

		t.Run("border routing", func(t *testing.T) {
			accept_ra, _, _ := utils.Exec(t, "sudo sysctl net.ipv6.conf."+strings.TrimSpace(INFRA_IF)+".accept_ra")
			require.Equal(t, "net.ipv6.conf."+strings.TrimSpace(INFRA_IF)+".accept_ra = 2\n", accept_ra)
			accept_ra_rt_info_max_plen, _, _ := utils.Exec(t, "sudo sysctl net.ipv6.conf."+strings.TrimSpace(INFRA_IF)+".accept_ra_rt_info_max_plen")
			require.Equal(t, "net.ipv6.conf."+strings.TrimSpace(INFRA_IF)+".accept_ra_rt_info_max_plen = 64\n", accept_ra_rt_info_max_plen)
		})
	})
}

func TestSocketFile(t *testing.T) {
	t.Run("Socket File", func(t *testing.T) {
		waitForFileCreation(t, "/run/snap.openthread-border-router/openthread-wpan0.sock", 10)
	})
}

func TestSnapServicesStatus(t *testing.T) {
	t.Run("Snap services status", func(t *testing.T) {
		// oneshot service
		require.False(t, utils.SnapServicesActive(t, otbrSetupApp))

		// actice services
		require.True(t, utils.SnapServicesActive(t, otbrWebApp))
		require.True(t, utils.SnapServicesActive(t, otbrAgentApp))
	})
}

func TestThreadNetworkFormation(t *testing.T) {
	t.Cleanup(func() {
		utils.Exec(t, "sudo openthread-border-router.ot-ctl thread stop")
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

	// set GitHub Action network interface
	utils.SnapSet(nil, otbrSnap, "infra-if", "eth0")

	utils.SnapStart(nil, otbrSnap)

	return
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

func readLogFile(t *testing.T, filePath string) (string, error) {
	t.Helper()

	text, err := os.ReadFile(filePath)

	if err != nil {
		return "", err
	}
	return string(text), nil
}
