# OpenThread Border Router snap

[![openthread-border-router](https://snapcraft.io/openthread-border-router/badge.svg)](https://snapcraft.io/openthread-border-router)

OpenThread Border Router (OTBR) is a Thread border router for POSIX-based platforms.

The snap packaging makes it easy to setup and run the OTBR on Linux.
**It has been tested on Ubuntu (Desktop/Server/Core) as a Border Router in Matter systems.**
For a list of limitations, refer to [#16](https://github.com/canonical/openthread-border-router-snap/issues/16).

OpenThread Border Router source code: https://github.com/openthread/ot-br-posix

For issues related to this snap, refer [here](https://github.com/canonical/openthread-border-router-snap/issues).

Usage instructions are available **[here](https://canonical-matter.readthedocs-hosted.com/en/latest/how-to/otbr-on-ubuntu/)**.

## Development

### Build the snap

Build locally for the same architecture as the host:

```bash
snapcraft -v
```

Build remotely for all supported architectures:

```bash
snapcraft remote-build
```

### Install the snap

Given the snap package file with `.snap` extension, install:

```bash
sudo snap install --dangerous *.snap
```

### Connect interfaces

When installing this snap from the store, some interfaces are automatically connected.
If you install the locally built snap, you need to connect these interfaces manually:

```bash
# Allow DNS-SD registration and discovery
sudo snap connect openthread-border-router:avahi-control
# Allow setting up the firewall
sudo snap connect openthread-border-router:firewall-control
# Allow access to USB Thread Radio Co-Processor (RCP)
sudo snap connect openthread-border-router:raw-usb
# Allow setting up the networking
sudo snap connect openthread-border-router:network-control
# Allow controlling the Bluetooth devices
sudo snap connect openthread-border-router:bluetooth-control
# Allow device discovery over Bluetooth Low Energy
sudo snap connect openthread-border-router:bluez
```

On Ubuntu Core the `avahi-control` and `bluez` interfaces are not provided by the system.
These interfaces should be consumed from other snaps, such as the [Avahi](https://snapcraft.io/avahi) and [BlueZ](https://snapcraft.io/bluez) snaps.
To install these snaps, and establish connections for the `avahi-control` interface from the `avahi` snap, and the `service` interface from the `bluez` snap, run:

```bash
sudo snap install avahi bluez
sudo snap connect openthread-border-router:avahi-control avahi:avahi-control
sudo snap connect openthread-border-router:bluez bluez:service
```
