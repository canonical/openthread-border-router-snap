#!/bin/bash -x

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

