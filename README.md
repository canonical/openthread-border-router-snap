# openthread-border-router-snap

OpenThread Border Router is a Thread border router for POSIX-based platforms.

This is a snap packaing of OpenThread Border Router. The snap packaging makes it easy to setup and run OpenThread Border Router.

The snap is **NOT SUPPORTED** by the OpenThread Border Router community.

OpenThread Border Router source code: https://github.com/openthread/ot-br-posix

## Snap Build and Installation
Execute the following command from the top-level directory of this repo:
```
snapcraft -v
```

This will create a snap package file with .snap extension. It can be installed locally by setting the `--devmode` flag:

```
snap install ./openthread-border-router_*.snap --devmode
```
Once the installation is complete, the following services should be running:
```
$ snap services openthread-border-router
Service                                 Startup  Current   Notes
openthread-border-router.ot-ctl         enabled  inactive  -
openthread-border-router.otbr-agent     enabled  inactive  -
openthread-border-router.otbr-firewall  enabled  inactive  -
openthread-border-router.otbr-nat44     enabled  inactive  -
openthread-border-router.otbr-web       enabled  active    -
```

## Viewing logs
```
snap logs -f openthread-border-router
```

