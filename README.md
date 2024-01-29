# openthread-border-router-snap
[![openthread-border-router](https://snapcraft.io/openthread-border-router/badge.svg)](https://snapcraft.io/openthread-border-router)

OpenThread Border Router (OTBR) is a Thread border router for POSIX-based platforms.

The snap packaging makes it easy to setup and run the OTBR on Linux.
**It has been tested on Ubuntu (Desktop/Server/Core) as a Border Router in Matter systems.**
For a list of limitations, refer to [#16](https://github.com/canonical/openthread-border-router-snap/issues/16).

OpenThread Border Router source code: https://github.com/openthread/ot-br-posix

For issues related to this snap, refer [here](https://github.com/canonical/openthread-border-router-snap/issues).

Usage instructions are available below and on the **[wiki](https://github.com/canonical/openthread-border-router-snap/wiki)**.

## Install
The snap can be installed from the store:
```bash
sudo snap install openthread-border-router
```

To build locally and install, refer [here](#build).

## Configure

### Set application configurations
View default configurations:
```bash
$ sudo snap get openthread-border-router 
Key        Value
infra-if   wlan0
radio-url  spinel+hdlc+uart:///dev/ttyACM0
thread-if  wpan0
```

Change using `sudo snap set openthread-border-router key="value"`.

For technical reasons, it is currently not allowed to change the value of `thread-if`; see [#17](https://github.com/canonical/openthread-border-router-snap/issues/17).

> **Note**  
> By default, the services are disabled and not started.
> They can be started and enabled as described [here](#run).
> To start and enable via a [Gadget snap](https://snapcraft.io/docs/the-gadget-snap), set `autostart` to `true`.
>

### Grant access to resources

Connect interfaces to access desired resources:
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

> **Note**  
> On **Ubuntu Core**, the `avahi-control` and `bluez` interfaces are not provided by the system.
> These interfaces should be consumed from other snaps, such as the [Avahi](https://snapcraft.io/avahi) and [BlueZ](https://snapcraft.io/bluez) snaps.
> 
> To install the snaps, and establish connections for the `avahi-control` interface from the `avahi` snap, and the `service` interface from the `bluez` snap, run:
> ```bash
> sudo snap install avahi bluez
> sudo snap connect openthread-border-router:avahi-control avahi:avahi-control
> sudo snap connect openthread-border-router:bluez bluez:service
> ```

## Run
Start once:
```bash
sudo snap start openthread-border-router
```
Add `--enable` flag to make the service start on boot as well.

To query and follow the logs: `snap logs -n 10 -f openthread-border-router`

## Usage

### Control a Matter Thread device
To commission and control a Matter Thread device via the OTBR Snap, please refer to the [wiki](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap).

### Pre-Shared Key for the Commissioner (PSKc) generator

To authenticate an external Thread Commissioner to a Thread network:

```bash
openthread-border-router.pskc J01NME 1234AAAA1234BBBB MyOTBRNetwork
```

where:

- `J01NME` is the passphrase/commissioner credential. It is a user-defined string encoded in UTF-8 and should be between 6 and 255 characters long.
- `1234AAAA1234BBBB` is the Extended PAN ID that was used in the operational dataset when forming the OTBR network.
- `MyOTBRNetwork` is the Network Name that was used in the operational dataset when forming the OTBR network.

### Steering data generator

To generate a hash of the set of Joiners intended for commissioning:

```bash
openthread-border-router.steering-data 8 0000b57fffe15d68
```

where:

- `8` is the Byte length of steering data (optional, default is 16).
- `0000b57fffe15d68` is the Joiner ID (EUI-64).

## Build

Build locally for the same architecture as the host:
```bash
snapcraft -v
```

Build remotely for all supported architectures:
```
snapcraft remote-build
```

Given the snap package file with `.snap` extension, install:
```bash
sudo snap install --dangerous *.snap
```
