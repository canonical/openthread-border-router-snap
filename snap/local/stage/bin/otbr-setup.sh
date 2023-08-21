#!/bin/bash -ex

#
#  Copyright (c) 2017, The OpenThread Authors.
#  Copyright (c) 2023, Canonical Ltd.
#  All rights reserved.
#
#  Redistribution and use in source and binary forms, with or without
#  modification, are permitted provided that the following conditions are met:
#  1. Redistributions of source code must retain the above copyright
#     notice, this list of conditions and the following disclaimer.
#  2. Redistributions in binary form must reproduce the above copyright
#     notice, this list of conditions and the following disclaimer in the
#     documentation and/or other materials provided with the distribution.
#  3. Neither the name of the copyright holder nor the
#     names of its contributors may be used to endorse or promote products
#     derived from this software without specific prior written permission.
#
#  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
#  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
#  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
#  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
#  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
#  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
#  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
#  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
#  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
#  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
#  POSSIBILITY OF SUCH DAMAGE.

echo "OTBR oneshot setup service"

# The operations performed here get reversed upon reboot and removal of the snap.
# The removal script is at snap/hooks/remove

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
# Skip setting "88 openthread" routing table mapping:
# https://github.com/canonical/openthread-border-router-snap/issues/14
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
