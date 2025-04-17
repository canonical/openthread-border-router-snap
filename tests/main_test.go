package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/canonical/matter-snap-testing/env"
	"github.com/canonical/matter-snap-testing/utils"
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

	defaultRadioURL = "spinel+hdlc+forkpty:///var/snap/openthread-border-router/common/ot-rcp-simulator-thread-reference-20230119-amd64?forkpty-arg=1"
	radioUrlKey     = "radio-url"
	radioUrlEnv     = "RADIO_URL"
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

func setup() (teardown func(), err error) {
	log.Println("[CLEAN]")
	utils.SnapRemove(nil, otbrSnap)

	start := time.Now()
	log.Println("[SETUP]")

	teardown = func() {
		log.Println("[TEARDOWN]")
		utils.SnapDumpLogs(nil, start, otbrSnap)

		log.Println("Removing installed snap:", env.Teardown())
		if env.Teardown() {
			utils.SnapRemove(nil, otbrSnap)
		}
	}

	if env.SnapPath() != "" {
		err = utils.SnapInstallFromFile(nil, env.SnapPath())
	} else {
		err = utils.SnapInstallFromStore(nil, otbrSnap, env.SnapChannel())
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
	utils.Exec(nil, "sudo cp ot-rcp-* /var/snap/openthread-border-router/common/")
	if v := os.Getenv(radioUrlEnv); v != "" {
		utils.SnapSet(nil, otbrSnap, radioUrlKey, v)
	} else {
		utils.SnapSet(nil, otbrSnap, radioUrlKey, defaultRadioURL)
	}

	// Get and set infrastructure interface
	if v := os.Getenv(infraInterfaceEnv); v != "" {
		infraInterfaceValue = v
		utils.SnapSet(nil, otbrSnap, infraInterfaceKey, infraInterfaceValue)
	} else {
		utils.SnapSet(nil, otbrSnap, infraInterfaceKey, defaultInfraInterfaceValue)
	}

	return
}
