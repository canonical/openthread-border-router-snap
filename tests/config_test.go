package tests

import (
	"os/exec"
	"testing"
	"time"

	"github.com/canonical/matter-snap-testing/utils"
	"github.com/stretchr/testify/require"
)

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
		testSettingSnapOption(t, configKey, configValue, defaultConfigValue, otbrSetupApp, expectedLog)
	})
	t.Run("Set radio-url", func(t *testing.T) {
		configKey := "radio-url"
		configValue := "spinel+hdlc+uart:///dev/ttyACM1"
		defaultConfigValue := defaultRadioURL
		expectedLog := "RADIO_URL=" + configValue
		testSettingSnapOption(t, configKey, configValue, defaultConfigValue, otbrAgentApp, expectedLog)
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
		configValue := "127.0.0.1"
		defaultConfigValue := "::"

		t.Cleanup(func() {
			utils.SnapSet(t, otbrSnap, configKey, defaultConfigValue)
			utils.SnapStop(t, otbrSnap)
		})

		utils.RequirePortAvailable(t, defaultWebGUIPort)
		utils.SnapSet(t, otbrSnap, configKey, configValue)
		utils.SnapStart(t, otbrSnap)
		utils.WaitServiceOnline(t, 10, defaultWebGUIPort)
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

		utils.WaitServiceOnline(t, 10, configValue)
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

func testSettingSnapOption(t *testing.T, configKey, configValue, defaultConfigValue, otbrService, expectedLog string) {
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
