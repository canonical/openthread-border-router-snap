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

## Viewing logs
```bash
snap logs -f openthread-border-router
```

