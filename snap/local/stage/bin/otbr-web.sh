#!/bin/bash -ex

export INTERFACE_NAME=$(snapctl get interface-name)
export LISTEN_ADDRESS=$(snapctl get listen-address)
export PORT=$(snapctl get port)

exec $SNAP/bin/otbr-web -I $INTERFACE_NAME -a $LISTEN_ADDRESS -p $PORT

# Usage: otbr-web [-d DEBUG_LEVEL] [-I interfaceName] [-p port] [-a listenAddress] [-v]
