#!/bin/bash -ex

echo "OTBR oneshot setup service"

# INFRA_IF_NAME=$(snapctl get infra-if-name)
# export INFRA_IF_NAME="${INFRA_IF_NAME:-wlp3s0}"
# echo "INFRA_IF_NAME=$INFRA_IF_NAME"
export INFRA_IF_NAME=$(snapctl get infra-if)

###############################################################################
echo "Setup the firewall"
$SNAP/bin/script/otbr-firewall start

# Custom _initrc without staging and build logic
export PLATFORM=ubuntu
source $SNAP/bin/_initrc_install

###############################################################################
echo "Setup IP forwarding"
source $SNAP/bin/script/_ipforward
ipforward_install

###############################################################################
echo "Setup RT Tables for the Backbone Router"
export BACKBONE_ROUTER=1
source $SNAP/bin/script/_rt_tables
rt_tables_install

###############################################################################
echo "Setup NAT44"
# The nat44_install function in scripts/_nat64 creates a service file and sets 
# firewall rules inside. We are only interested in the firewall rules:
$SNAP/bin/nat44_install

###############################################################################
echo "Setup Border Routing"
export BORDER_ROUTING=1
source $SNAP/bin/script/_border_routing
border_routing_install

echo "✔️ OTBR completed oneshot setup"
