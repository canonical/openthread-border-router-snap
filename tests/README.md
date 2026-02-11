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

In this test, a pre-built RCP simulator will be utilized. 
The simulator needs to be built on Ubuntu 22.04 to be compatible with the glibc that will be available to this snap that uses the core22 base.
A multipass VM can be used for this.

```bash
git clone https://github.com/openthread/openthread.git --branch=thread-reference-20250612
cd openthread
./script/bootstrap
./script/cmake-build simulation
```

Once built, copy the simulator binary from the Multipass VM into this directory, and rename it to `ot-rcp-simulator-thread-reference-20250612-amd64`.
```bash
cp build/simulation/examples/apps/ncp/ot-rcp ~/<host-home-dir>/openthread-border-router-snap/tests/ot-rcp-simulator-thread-reference-20250612-amd64
```

During testing the tests will copy the simulator into the `$SNAP_COMMON` directory,
and the `radio-url` will be set to its path.
```bash
snap set openthread-border-router radio-url='spinel+hdlc+forkpty:///var/snap/openthread-border-router/common/ot-rcp-simulator-thread-reference-20230119-amd64?forkpty-arg=1''
```

For additional information regarding RCP simulation, please refer to the [openthread simulation posix](https://openthread.io/codelabs/openthread-simulation-posix#3).
