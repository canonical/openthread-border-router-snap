# Run Tests

```bash
go test -v -failfast -count 1
```

where:
- `-v` is to enable verbose output
- `-failfast` makes the test stop after first failure
- `-count 1` is to avoid Go test caching for example when testing a rebuilt snap

# Override behavior
Use environment variables, as defined in [utils/env.go](https://github.com/canonical/matter-snap-testing/blob/main/utils/env.go).

For the infrastructure interface name, the default value is "wlan0". 
To override it during the Go test, run:

```bash
INFRA_IF="eth0" go test -v -failfast -count 1
```

# Build a Radio Co-Processor (RCP) simulator for testing

The tests use an RCP simulator for the host architecture. Simulators for amd64
and arm64 are included in the repository via Git LFS. Build replacement
simulators on Ubuntu 24.04 for the target architecture so they are compatible
with the snap's core24 base.
A Multipass VM can be used for this. Install `g++`, `cmake`, and `ninja-build`
before building.

```bash
git clone https://github.com/openthread/openthread.git --branch=thread-reference-20250612
cd openthread
OT_CMAKE_NINJA_TARGET=ot-rcp ./script/cmake-build simulation
```

Once built, copy the simulator binary into this directory with the architecture
suffix used by Go (`amd64` or `arm64`). For example, on arm64:
```bash
cp build/simulation/examples/apps/ncp/ot-rcp ~/<host-home-dir>/openthread-border-router-snap/tests/ot-rcp-simulator-thread-reference-20250612-arm64
```

During testing the tests will copy the simulator into the `$SNAP_COMMON` directory,
and the `radio-url` will be set to its path.
```bash
snap set openthread-border-router radio-url='spinel+hdlc+forkpty:///var/snap/openthread-border-router/common/ot-rcp-simulator-thread-reference-20250612-amd64?forkpty-arg=1'
```

For additional information regarding RCP simulation, please refer to the [openthread simulation posix](https://openthread.io/codelabs/openthread-simulation-posix#3).
