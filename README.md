# openthread-border-router-snap

OpenThread Border Router is a Thread border router for POSIX-based platforms.

This is a snap packaging of OpenThread Border Router. The snap packaging makes it easy to setup and run the OpenThread Border Router.

The snap is **NOT SUPPORTED** by the OpenThread Border Router community.

OpenThread Border Router source code: https://github.com/openthread/ot-br-posix

## Snap Build and Installation
Execute the following command from the top-level directory of this repo:
```bash
snapcraft -v
```

This will create a snap package file with .snap extension. It can be installed locally by setting the `--devmode` flag:
```bash
sudo snap install ./openthread-border-router_*.snap --devmode
```

## Configure
View default configurations:
```bash
$ sudo snap get openthread-border-router 
Key        Value
infra-if   wlp3s0
radio-url  spinel+hdlc+uart:///dev/ttyACM0
thread-if  wpan0
```

Change using `sudo snap openthread-border-router set key="value"`

## Grant access

Connect interfaces to access desired resources:
```bash
# Allow access to specific system files
sudo snap connect openthread-border-router:system-etc-iproute
sudo snap connect openthread-border-router:system-etc-sysctl
sudo snap connect openthread-border-router:system-run-openthread-wpan0

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

> **_NOTE:_**  If the thread interface is different from the default value (wpan0), instead of `system-run-openthread-wpan0`, connect:
> ```bash
> sudo snap connect openthread-border-router:system-run
> ```

## Run
```bash
sudo snap start openthread-border-router
```

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

## Viewing logs
```bash
snap logs -f openthread-border-router
```

