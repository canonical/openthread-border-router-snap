# OpenThread Border Router snap

[![openthread-border-router](https://snapcraft.io/openthread-border-router/badge.svg)](https://snapcraft.io/openthread-border-router)

OpenThread Border Router (OTBR) is a Thread border router for POSIX-based platforms.

The snap packaging makes it easy to setup and run the OTBR on Linux.
**It has been tested on Ubuntu (Desktop/Server/Core) as a Border Router in Matter systems.**
For a list of limitations, refer to [#16](https://github.com/canonical/openthread-border-router-snap/issues/16).

OpenThread Border Router source code: https://github.com/openthread/ot-br-posix

For issues related to this snap, refer [here](https://github.com/canonical/openthread-border-router-snap/issues).

Usage instructions are available **[here](https://canonical-matter.readthedocs-hosted.com/en/latest/how-to/otbr-on-ubuntu/)**.

## Build

Build locally for the same architecture as the host:

```bash
snapcraft -v
```

Build remotely for all supported architectures:

```bash
snapcraft remote-build
```

Given the snap package file with `.snap` extension, install:

```bash
sudo snap install --dangerous *.snap
```

## Advanced Usage

### Thread Interface

For technical reasons, it is currently not allowed to change the value of `thread-if`; see [#17](https://github.com/canonical/openthread-border-router-snap/issues/17).

### Starting the service

By default, the services are disabled and not started.
They can be started and enabled as described [here](https://canonical-matter.readthedocs-hosted.com/en/latest/how-to/otbr-on-ubuntu/#start-otbr).
To start and enable via a [Gadget snap](https://snapcraft.io/docs/the-gadget-snap), set `autostart` to `true`.

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
