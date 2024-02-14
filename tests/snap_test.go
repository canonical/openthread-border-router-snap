package tests

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/canonical/matter-snap-testing/utils"
	"github.com/stretchr/testify/require"
)

const (
	otbrSnap     = "openthread-border-router"
	otbrSetupApp = "openthread-border-router.otbr-setup"
	otbrAgentApp = "openthread-border-router.otbr-agent"
	otbrWebApp   = "openthread-border-router.otbr-web"

	defaultInfraInterfaceValue = "wlan0"
	infraInterfaceKey          = "infra-if"
	infraInterfaceEnv          = "INFRA_IF"

	defaultWebGUIPort = "80"
)

var infraInterfaceValue = defaultInfraInterfaceValue

func TestMain(m *testing.M) {
	teardown, err := setup()
	if err != nil {
		log.Fatalf("Failed to setup tests: %s", err)
	}

	code := m.Run()
	teardown()

	os.Exit(code)
}

// TestSetup verifies the configuration setup of OTBR based on otbr-setup.sh script:
// https://github.com/canonical/openthread-border-router-snap/blob/main/snap/local/stage/bin/otbr-setup.sh
func TestSetup(t *testing.T) {
	infraInterfaceFromSnap, _, _ := utils.Exec(t, "sudo snap get openthread-border-router "+infraInterfaceKey)
	trimmedInfraInterface := strings.TrimSpace(infraInterfaceFromSnap)

	// Start clean
	utils.SnapStop(t, otbrSnap)

	t.Cleanup(func() {
		utils.SnapStop(t, otbrSnap)
	})

	utils.SnapStart(nil, otbrSnap)
	t.Run("Firewall rule", func(t *testing.T) {
		forwardRules, _, _ := utils.Exec(t, "sudo iptables -S FORWARD | grep \"comment OTBR\"")
		require.NotEmpty(t, forwardRules)
	})

	t.Run("IP forwarding", func(t *testing.T) {
		ipv6Forwarding, _, _ := utils.Exec(t, "sudo sysctl net.ipv6.conf.all.forwarding")
		require.Equal(t, "net.ipv6.conf.all.forwarding = 1\n", ipv6Forwarding)
		ipv4Forwarding, _, _ := utils.Exec(t, "sudo sysctl net.ipv4.ip_forward")
		require.Equal(t, "net.ipv4.ip_forward = 1\n", ipv4Forwarding)
	})

	t.Run("RT tables for backbone router", func(t *testing.T) {
		socketBufferSize, _, _ := utils.Exec(t, "sudo sysctl net.core.optmem_max")
		require.Equal(t, "net.core.optmem_max = 65536\n", socketBufferSize)
	})

	t.Run("Random fwmark bits", func(t *testing.T) {
		mangleTablePreroutingChain, _, _ := utils.Exec(t, "sudo iptables -t mangle -L PREROUTING -n -v | grep OTBR")
		require.NotEmpty(t, mangleTablePreroutingChain)
		natTablePostroutingChain, _, _ := utils.Exec(t, "sudo iptables -t nat -L POSTROUTING -n -v | grep OTBR")
		require.NotEmpty(t, natTablePostroutingChain)
	})

	t.Run("Firewall rule setup for Infrastructure interface", func(t *testing.T) {
		forwardRule, _, _ := utils.Exec(t, "sudo iptables -t filter -L FORWARD -n -v | grep OTBR | grep "+trimmedInfraInterface)
		require.NotEmpty(t, forwardRule)
	})

	t.Run("Border routing", func(t *testing.T) {
		acceptRA, _, _ := utils.Exec(t, "sudo sysctl net.ipv6.conf."+trimmedInfraInterface+".accept_ra")
		require.Equal(t, "net.ipv6.conf."+trimmedInfraInterface+".accept_ra = 2\n", acceptRA)
		acceptRARTInfoMaxPlen, _, _ := utils.Exec(t, "sudo sysctl net.ipv6.conf."+trimmedInfraInterface+".accept_ra_rt_info_max_plen")
		require.Equal(t, "net.ipv6.conf."+trimmedInfraInterface+".accept_ra_rt_info_max_plen = 64\n", acceptRARTInfoMaxPlen)
	})
}

func TestSocketFile(t *testing.T) {
	// Start clean
	utils.SnapStop(t, otbrSnap)

	t.Cleanup(func() {
		utils.SnapStop(t, otbrSnap)
	})

	utils.SnapStart(nil, otbrSnap)
	waitForFileCreation(t, "/run/snap.openthread-border-router/openthread-wpan0.sock", 10)
}

func TestSnapServicesStatus(t *testing.T) {
	// Start clean
	utils.SnapStop(t, otbrSnap)

	t.Cleanup(func() {
		utils.SnapStop(t, otbrSnap)
	})

	utils.SnapStart(nil, otbrSnap)

	// Oneshot service
	require.False(t, utils.SnapServicesActive(t, otbrSetupApp))

	// Actice services
	require.True(t, utils.SnapServicesActive(t, otbrWebApp))
	require.True(t, utils.SnapServicesActive(t, otbrAgentApp))
}

func TestConfig(t *testing.T) {
	serviceWaitTimeout := 10

	// Start clean
	utils.SnapStop(t, otbrSnap)

	t.Cleanup(func() {
		utils.SnapStop(t, otbrSnap)
	})

	t.Run("Set infra-if", func(t *testing.T) {
		configKey := infraInterfaceKey
		configValue := "wpan1"
		defaultConfigValue := infraInterfaceValue
		expectedLog := infraInterfaceEnv + "=" + configValue
		checkSnapOptions(t, configKey, configValue, defaultConfigValue, otbrSetupApp, expectedLog)
	})
	t.Run("Set radio-url", func(t *testing.T) {
		configKey := "radio-url"
		configValue := "spinel+hdlc+uart:///dev/ttyACM1"
		defaultConfigValue := "spinel+hdlc+uart:///dev/ttyACM0"
		expectedLog := "RADIO_URL=" + configValue
		checkSnapOptions(t, configKey, configValue, defaultConfigValue, otbrAgentApp, expectedLog)
	})

	t.Run("Set invalid thread interface", func(t *testing.T) {
		configKey := "thread-if"
		defaultConfigValue := "wpan0"
		invalidConfigValue := "wpan1"

		t.Cleanup(func() {
			utils.SnapSet(t, otbrSnap, configKey, defaultConfigValue)
			utils.SnapStop(t, otbrSnap)
		})

		command := "sudo snap set openthread-border-router thread-if=" + invalidConfigValue
		output, err := exec.Command("/bin/bash", "-c", command).CombinedOutput()
		t.Logf("[exec] %s", command)

		require.NotEmpty(t, output)
		require.Error(t, err, "Expected an error while setting an invalid thread interface")
	})

	t.Run("Set webgui-listen-address", func(t *testing.T) {
		configKey := "webgui-listen-address"
		configValue := "192.168.178.1"
		defaultConfigValue := "::"

		t.Cleanup(func() {
			utils.SnapSet(t, otbrSnap, configKey, defaultConfigValue)
			utils.SnapStop(t, otbrSnap)
		})

		utils.SnapSet(t, otbrSnap, configKey, configValue)
		utils.SnapStart(t, otbrSnap)

		stdout, _, _ := utils.Exec(t, "curl "+configValue+":"+defaultWebGUIPort)
		require.NotEmpty(t, stdout)
	})
	t.Run("Set webgui-port", func(t *testing.T) {
		configKey := "webgui-port"
		configValue := "90"
		defaultConfigValue := defaultWebGUIPort

		t.Cleanup(func() {
			utils.SnapSet(t, otbrSnap, configKey, defaultConfigValue)
			utils.SnapStop(t, otbrSnap)
		})

		utils.RequirePortAvailable(t, configValue)
		utils.SnapSet(t, otbrSnap, configKey, configValue)
		utils.SnapStart(nil, otbrSnap)
		utils.WaitServiceOnline(t, serviceWaitTimeout, configValue)

		stdout, _, _ := utils.Exec(t, "curl localhost:"+configValue)
		require.NotEmpty(t, stdout)
	})

	t.Run("Set autostart", func(t *testing.T) {
		t.Cleanup(func() {
			utils.SnapStop(t, otbrSnap)
		})

		require.False(t, utils.SnapServicesEnabled(t, otbrSnap))
		require.False(t, utils.SnapServicesActive(t, otbrSnap))

		utils.SnapSet(t, otbrSnap, "autostart", "true")
		require.True(t, utils.SnapServicesEnabled(t, otbrAgentApp))
		require.True(t, utils.SnapServicesActive(t, otbrAgentApp))
		require.True(t, utils.SnapServicesEnabled(t, otbrWebApp))
		require.True(t, utils.SnapServicesActive(t, otbrWebApp))
		require.True(t, utils.SnapServicesEnabled(t, otbrSetupApp))
		require.False(t, utils.SnapServicesActive(t, otbrSetupApp))
	})
}

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

func setup() (teardown func(), err error) {

	log.Println("[CLEAN]")
	utils.SnapRemove(nil, otbrSnap)

	start := time.Now()
	log.Println("[SETUP]")

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

	// Connect interfaces
	utils.SnapConnect(nil, otbrSnap+":avahi-control", "")
	utils.SnapConnect(nil, otbrSnap+":firewall-control", "")
	utils.SnapConnect(nil, otbrSnap+":raw-usb", "")
	utils.SnapConnect(nil, otbrSnap+":network-control", "")
	utils.SnapConnect(nil, otbrSnap+":bluetooth-control", "")
	utils.SnapConnect(nil, otbrSnap+":bluez", "")

	// Copy and set simulated RCP
	utils.Exec(nil, "sudo cp ot-rcp-simulator-thread-reference-20230119-amd64 /var/snap/openthread-border-router/common/")
	utils.SnapSet(nil, otbrSnap, "radio-url", "'spinel+hdlc+forkpty:///var/snap/openthread-border-router/common/ot-rcp-simulator-thread-reference-20230119-amd64?forkpty-arg=1'")

	// Get and set infrastructure interface
	if v := os.Getenv(infraInterfaceEnv); v != "" {
		infraInterfaceValue = v
		utils.SnapSet(nil, otbrSnap, infraInterfaceKey, infraInterfaceValue)
	} else {
		utils.SnapSet(nil, otbrSnap, infraInterfaceKey, defaultInfraInterfaceValue)
	}

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

func checkSnapOptions(t *testing.T, configKey, configValue, defaultConfigValue, otbrService, expectedLog string) {
	t.Helper()

	// Start clean
	utils.SnapStop(t, otbrSnap)
	start := time.Now()

	t.Cleanup(func() {
		utils.SnapSet(t, otbrSnap, configKey, defaultConfigValue)
		utils.SnapStop(t, otbrSnap)
	})

	utils.SnapSet(t, otbrSnap, configKey, configValue)
	command := "sudo snap start openthread-border-router"
	_, _ = exec.Command("/bin/bash", "-c", command).CombinedOutput()
	t.Logf("[exec] %s", command)
	utils.WaitForLogMessage(t, otbrService, expectedLog, start)
}
