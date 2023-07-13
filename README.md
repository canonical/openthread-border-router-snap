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

View default configurations:
```bash
$ sudo snap get openthread-border-router 
Key        Value
infra-if   wlp3s0
radio-url  spinel+hdlc+uart:///dev/ttyACM0
thread-if  wpan0
```

Change using `sudo snap openthread-border-router set key="value"`

## Usage

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

