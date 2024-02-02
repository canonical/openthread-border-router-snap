# Run Tests

```bash
go test -v -failfast -count 1
```

where:
- `-v` is to enable verbose output
- `-failfast` makes the test stop after first failure
- `-count 1` is to avoid Go test caching for example when testing a rebuilt snap

# Build a Radio Co-Processor (RCP) simulator for testing

In this test, a pre-built RCP simulator will be utilized. 
The RCP simulator can be built by running the following commands:
```bash
git clone https://github.com/openthread/openthread.git --branch=thread-reference-20230119
cd openthread
./script/bootstrap
./script/cmake-build simulation
```

Once built, it needs to be relocated to a directory visible to OTBR snap 
and subsequently passed to OTBR snap using snap options during testing:
```bash
sudo cp openthread/build/simulation/examples/apps/ncp/ot-rcp/ot-rcp-simulator-thread-reference-20230119-amd64 \
    /var/snap/openthread-border-router/common/
snap set openthread-border-router radio-url='spinel+hdlc+forkpty:///var/snap/openthread-border-router/common/ot-rcp-simulator-thread-reference-20230119-amd64?forkpty-arg=1''
```

For additional information regarding RCP simulation, please refer to the [openthread simulation posix](https://openthread.io/codelabs/openthread-simulation-posix#3).

