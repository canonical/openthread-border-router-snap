#!/bin/bash -ex

export OTBR_NO_AUTO_ATTACH=0

export INFRA_IF_NAME=$(snapctl get infra-if)
export THREAD_IF=$(snapctl get thread-if)
export RADIO_URL=$(snapctl get radio-url)
export VENDOR_NAME=$(snapctl get vendor-name)
export VENDOR_MODEL=$(snapctl get vendor-model)

exec $SNAP/bin/otbr-agent -v \
    -I $THREAD_IF \
    -B $INFRA_IF_NAME \
    --vendor-name "$VENDOR_NAME" \
    --model-name "$VENDOR_MODEL" \
    $RADIO_URL

# See https://github.com/openthread/ot-br-posix/blob/main/src/agent/otbr-agent.default.in for agent default settings.

# Usage: otbr-agent [-I interfaceName] [-B backboneIfName] [-d DEBUG_LEVEL] [-v] [-s] [--auto-attach[=0/1]] RADIO_URL [RADIO_URL]
#         --data-path        Path of directory to store data.
#     -I, --thread-ifname    Name of the Thread network interface (default: wpan0).
#     -B, --backbone-ifname  Name of the backbone network interfaces (can be specified multiple times).
#     -d, --debug-level      The log level (EMERG=0, ALERT=1, CRIT=2, ERR=3, WARNING=4, NOTICE=5, INFO=6, DEBUG=7).
#     -v, --verbose          Enable verbose logging.
#     -s, --syslog-disable   Disable syslog and print to standard error.
#     -h, --help             Show this help text.
#     -V, --version          Print the application's version and exit.
#     --radio-version        Print the radio coprocessor version and exit.
#     --auto-attach          Whether or not to automatically attach to the saved network (default: 1).
#     --rest-listen-address  Network address to listen on for the REST API (default: 127.0.0.1).
#     --rest-listen-port     Network port to listen on for the REST API (default: 8081).
#     --vendor-name          Vendor Name.
#     --model-name           Model Name.
#
# RadioURL:
# Radio Url format:    {Protocol}://${PATH_TO_DEVICE}?${Parameters}
#
# Protocol=[spinel+hdlc*]           Specify the Spinel interface as the Spinel HDLC interface
#     forkpty-arg[=argument string]  Command line arguments for subprocess, can be repeated.
#     spinel+hdlc+uart://${PATH_TO_UART_DEVICE}?${Parameters} for real uart device
#     spinel+hdlc+forkpty://${PATH_TO_UART_DEVICE}?${Parameters} for forking a pty subprocess.
# Parameters:
#     uart-parity[=even|odd]         Uart parity config, optional.
#     uart-stop[=number-of-bits]     Uart stop bit, default is 1.
#     uart-baudrate[=baudrate]       Uart baud rate, default is 460800.
#     uart-flow-control              Enable flow control, disabled by default.
#     uart-init-deassert             Deassert lines on init when flow control is disabled.
#     uart-reset                     Reset connection after hard resetting RCP(USB CDC ACM).
#     uart-exclusive                 Lock uart device using flock / TIOCEXCL.
#
#     region[=region-code]          Set the radio's region code. The region code must be an
#                                   ISO 3166 alpha-2 code.
#     cca-threshold[=dbm]           Set the radio's CCA ED threshold in dBm measured at antenna connector.
#     enable-coex[=1|0]             If not specified, RCP coex operates with its default configuration.
#                                   Disable coex with 0, and enable it with other values.
#     fem-lnagain[=dbm]             Set the Rx LNA gain in dBm of the external FEM.
#     no-reset                      Do not send Spinel reset command to RCP on initialization.
#     skip-rcp-compatibility-check  Skip checking RCP API version and capabilities during initialization.
#     bus-latency[=usec]            Communication latency in usec, default is 0.
#     product-config-file[=path]    Specify a custom path to the openthread.conf product configuration
#                                   file.
#                                   If not specified, the default path set at build time is used.
#     factory-config-file[=path]    Specify a custom path to the openthread.conf factory configuration
#                                   file.
#                                   If not specified, the default path set at build time is used.
