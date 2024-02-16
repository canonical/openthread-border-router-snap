package tests

import (
	"strings"
	"testing"

	"github.com/canonical/matter-snap-testing/utils"
	"github.com/stretchr/testify/require"
)

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
