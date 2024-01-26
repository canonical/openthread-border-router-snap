#!/bin/bash -ex

export LISTEN_ADDRESS=$(snapctl get listen-address)
export PORT=$(snapctl get port)

exec $SNAP/bin/otbr-web -a "$LISTEN_ADDRESS" -p $PORT

# Usage: otbr-web [-d DEBUG_LEVEL] [-I interfaceName] [-p port] [-a listenAddress] [-v]
