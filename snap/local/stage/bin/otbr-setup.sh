#!/bin/bash -ex

echo "OTBR oneshot setup service"

INFRA_IF=$(snapctl get infra-if)
THREAD_IF=$(snapctl get thread-if)

###############################################################################
echo "Setup the firewall"
$SNAP/bin/script/otbr-firewall start

# Custom _initrc without staging and build logic
export PLATFORM=ubuntu
source $SNAP/bin/_initrc_install

###############################################################################
echo "Setup IP forwarding"
# Upstream equivalent:
# https://github.com/openthread/ot-br-posix/blob/thread-reference-20230119/script/_ipforward
sysctl -w net.ipv6.conf.all.forwarding=1
sysctl -w net.ipv4.ip_forward=1

###############################################################################
echo "Setup RT Tables for the Backbone Router"
# Upstream equivalent:
# https://github.com/openthread/ot-br-posix/blob/thread-reference-20230119/script/_rt_tables
# TODO: replace with ip command
sh -c 'echo "88 openthread" >>/etc/iproute2/rt_tables'
sysctl net.core.optmem_max=65536

###############################################################################
echo "Setup NAT44"
# The nat44_install function in scripts/_nat64 creates a service file and sets 
# firewall rules inside. We are only interested in the firewall rules.
# Upstream source: 
# https://github.com/openthread/ot-br-posix/blob/thread-reference-20230119/script/_nat64
echo "Set random fwmark bits"
iptables -t mangle -A PREROUTING -i $THREAD_IF -j MARK --set-mark 0x1001 -m comment --comment "OTBR"
iptables -t nat -A POSTROUTING -m mark --mark 0x1001 -j MASQUERADE -m comment --comment "OTBR"

echo "Setup NAT44: Setting firewall rule for $INFRA_IF"
iptables -t filter -A FORWARD -o $INFRA_IF -j ACCEPT -m comment --comment "OTBR"
iptables -t filter -A FORWARD -i $INFRA_IF -j ACCEPT -m comment --comment "OTBR"

###############################################################################
echo "Setup Border Routing"
# Upstream equivalent:
# https://github.com/openthread/ot-br-posix/blob/thread-reference-20230119/script/_border_routing
sysctl -w net.ipv6.conf.$INFRA_IF.accept_ra=2
sysctl -w net.ipv6.conf.$INFRA_IF.accept_ra_rt_info_max_plen=64

echo "✔️ OTBR completed oneshot setup"
