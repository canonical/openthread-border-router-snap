#!/bin/bash -x

#
#  Copyright (c) 2021, The OpenThread Authors.
#  Copyright (c) 2023, Canonical Ltd
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
#

ipset_destroy_if_exist()
{
    if ipset list "$1"; then
        ipset destroy "$1"
    fi
}

# 1. Prepare environment variables for the firewall
THREAD_IF="wpan0"
OTBR_FORWARD_INGRESS_CHAIN="OTBR_FORWARD_INGRESS"

. /lib/lsb/init-functions
. /lib/init/vars.sh

set -euxo pipefail
modprobe ip6table_filter || true

# 3. Stop exsiting firewall
if ip6tables -C FORWARD -o "$THREAD_IF" -j "$OTBR_FORWARD_INGRESS_CHAIN"; then
    ip6tables -D FORWARD -o "$THREAD_IF" -j "$OTBR_FORWARD_INGRESS_CHAIN" || true
fi


if ip6tables -L $OTBR_FORWARD_INGRESS_CHAIN; then
    ip6tables -w -F $OTBR_FORWARD_INGRESS_CHAIN
    ip6tables -w -X $OTBR_FORWARD_INGRESS_CHAIN
fi

ipset_destroy_if_exist otbr-ingress-deny-src
ipset_destroy_if_exist otbr-ingress-deny-src-swap
ipset_destroy_if_exist otbr-ingress-allow-dst
ipset_destroy_if_exist otbr-ingress-allow-dst-swap

# 4. Create OTBR_FORWARD_INGRESS chain
ip6tables -N "$OTBR_FORWARD_INGRESS_CHAIN"

# 5. Start firewall configuration
ipset create -exist otbr-ingress-deny-src hash:net family inet6
ipset create -exist otbr-ingress-deny-src-swap hash:net family inet6
ipset create -exist otbr-ingress-allow-dst hash:net family inet6
ipset create -exist otbr-ingress-allow-dst-swap hash:net family inet6

ip6tables -I FORWARD 1 -o $THREAD_IF -j $OTBR_FORWARD_INGRESS_CHAIN

ip6tables -A "$OTBR_FORWARD_INGRESS_CHAIN" -m pkttype --pkt-type unicast -i "$THREAD_IF" -j DROP
ip6tables -A "$OTBR_FORWARD_INGRESS_CHAIN" -m set --match-set otbr-ingress-deny-src src -j DROP
ip6tables -A "$OTBR_FORWARD_INGRESS_CHAIN" -m set --match-set otbr-ingress-allow-dst dst -j ACCEPT
ip6tables -A "$OTBR_FORWARD_INGRESS_CHAIN" -m pkttype --pkt-type unicast -j DROP
ip6tables -A "$OTBR_FORWARD_INGRESS_CHAIN" -j ACCEPT

