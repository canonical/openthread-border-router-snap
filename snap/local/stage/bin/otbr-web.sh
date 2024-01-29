#!/bin/bash -ex

export THREAD_IF=$(snapctl get thread-if)
export WEBGUI_LISTEN_ADDRESS=$(snapctl get webgui-listen-address)
export WEBGUI_PORT=$(snapctl get webgui-port)

exec $SNAP/bin/otbr-web -I $THREAD_IF -a $WEBGUI_LISTEN_ADDRESS -p $WEBGUI_PORT

# Usage: otbr-web [-d DEBUG_LEVEL] [-I interfaceName] [-p port] [-a listenAddress] [-v]
