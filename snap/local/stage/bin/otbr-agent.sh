#!/bin/bash -ex

export OTBR_NO_AUTO_ATTACH=0

export INFRA_IF_NAME=$(snapctl get infra-if)
export THREAD_IF=$(snapctl get thread-if)
export RADIO_URL=$(snapctl get radio-url)

exec $SNAP/bin/otbr-agent -v -I $THREAD_IF -B $INFRA_IF_NAME $RADIO_URL

# /etc/default/otbr-agent:
# Default settings for otbr-agent. This file is sourced by systemd
# Options to pass to otbr-agent
# OTBR_AGENT_OPTS="-I wpan0 -B wlan0 spinel+hdlc+uart:///dev/ttyACM0 trel://wlan0"
# OTBR_NO_AUTO_ATTACH=0

# Usage: otbr-agent [-I interfaceName] [-B backboneIfName] [-d DEBUG_LEVEL] [-v] [--auto-attach[=0/1]] RADIO_URL [RADIO_URL]
#     --auto-attach defaults to 1
# RadioURL:
#     forkpty-arg[=argument string]  Command line arguments for subprocess, can be repeated.
#     spinel+hdlc+uart://${PATH_TO_UART_DEVICE}?${Parameters} for real uart device
#     spinel+hdlc+forkpty://${PATH_TO_UART_DEVICE}?${Parameters} for forking a pty subprocess.
# Parameters:
#     uart-parity[=even|odd]         Uart parity config, optional.
#     uart-stop[=number-of-bits]     Uart stop bit, default is 1.
#     uart-baudrate[=baudrate]       Uart baud rate, default is 115200.
#     uart-flow-control              Enable flow control, disabled by default.
#     uart-reset                     Reset connection after hard resetting RCP(USB CDC ACM).
#     region[=region-code]          Set the radio's region code. The region code must be an
#                                   ISO 3166 alpha-2 code.
#     cca-threshold[=dbm]           Set the radio's CCA ED threshold in dBm measured at antenna connector.
#     enable-coex[=1|0]             If not specified, RCP coex operates with its default configuration.
#                                   Disable coex with 0, and enable it with other values.
#     fem-lnagain[=dbm]             Set the Rx LNA gain in dBm of the external FEM.
#     ncp-dataset                   Retrieve dataset from ncp.
#     no-reset                      Do not send Spinel reset command to RCP on initialization.
#     skip-rcp-compatibility-check  Skip checking RCP API version and capabilities during initialization.
