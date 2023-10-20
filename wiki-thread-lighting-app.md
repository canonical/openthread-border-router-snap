## How to Set Up and Control a Matter Thread Lighting Device

This guide will walk you through the process of building, running, and controlling a Matter Thread lighting application over a Thread network. 

Here are step-by-step instructions:

#### Tested Hardware
- Machine A
  - nRF52840 dongle as [Radio Co-Processor (RCP)](https://openthread.io/platforms/co-processor#radio_co-processor_rcp)
  - PC running Ubuntu Desktop 23.04 OS
- Machine B
  - nRF52840 dongle as RCP
  - PC running Ubuntu Desktop 23.04 OS
 
#### Environment Setup
- Machine A running
  - [OpenThread Border Router (OTBR) snap](https://snapcraft.io/openthread-border-router)
  - Matter Thread ligting App
- Machine B running
  - OTBR snap
  - [Chip Tool snap](https://snapcraft.io/chip-tool)

#### Versions Tested in This Guide
- Matter SDK: [`6b01cb9`](https://github.com/project-chip/connectedhomeip/commit/6b01cb977127eb8547ce66d5b627061dc2dd6c90)
- OTBR SDK: [`thread-reference-20230119`](https://github.com/openthread/ot-br-posix/tree/thread-reference-20230119)
(API version: [6](https://github.com/openthread/openthread/blob/thread-reference-20230119/src/lib/spinel/spinel.h#L380))

- RCP firmware: [`00ac6cd`](https://github.com/openthread/ot-nrf528xx/tree/00ac6cd0137a4f09288b455bf8d7aa72d74062d1)
with openthread submodule commit [`9af0bfa`](https://github.com/openthread/openthread/tree/9af0bfa60e373d81a5576b298d6664045870a375)
(API version: [6](https://github.com/openthread/openthread/blob/9af0bfa60e373d81a5576b298d6664045870a375/src/lib/spinel/spinel.h#L420))

#### 1. Form a Thread Network on Machine A
On Machine A, following the guide to [install and configure the OTBR snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#install-and-configure-the-otbr-snap) 
and [form a Thread network](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#form-a-thread-network) 
using the nRF52840 Dongle as the RCP. For details on setting up the RCP, 
please refer to [Build and flash RCP firmware on nRF52480 dongle](https://github.com/canonical/openthread-border-router-snap/wiki/Setup-OpenThread-Border-Router-with-nRF52840-Dongle#build-and-flash-rcp-firmware-on-nrf52480-dongle).

#### 2. Run a Matter Thread Lighting App on Machine B

Build and run the [Matter Thread lighting app example](https://github.com/project-chip/connectedhomeip/tree/6b01cb977127eb8547ce66d5b627061dc2dd6c90/examples/lighting-app/linux) natively. By doing so, you will be able to observe the lighting status (on/off) from the logs later.

To get started, follow these steps:

1. Install the Matter SDK with a shallow clone for the Linux platform:
```bash
git clone https://github.com/project-chip/connectedhomeip.git --depth=1
cd connectedhomeip
git checkout 6b01cb977127eb8547ce66d5b627061dc2dd6c90
scripts/checkout_submodules.py --shallow --platform linux
```

2. Build the lighting app:
```bash
cd examples/lighting-app/linux
gn gen out/debug
ninja -C out/debug
```

3. Run the lighting app with Thread feature enabled:
```bash
./out/debug/chip-lighting-app --thread
```

#### 3. Run OpenThread Border Router Agent on Machine B

The lighting app communicates with OTBR agent via the dbus message bus. This is necessary to support the lighting app's Thread feature, as explained
[here](https://github.com/project-chip/connectedhomeip/tree/6b01cb977127eb8547ce66d5b627061dc2dd6c90/examples/lighting-app/linux#commandline-arguments).

Follow these steps to run OTBR agent:

1. Install and configure OTBR snap by referring [here](https://github.com/canonical/openthread-border-router-snap#readme)
2. Plug the nRF52840 dongle with the OTBR firmware. This is the same firmware used on the dongle connected to machine A. Please refer to [here](https://github.com/canonical/openthread-border-router-snap/wiki/Setup-OpenThread-Border-Router-with-nRF52840-Dongle#build-and-flash-rcp-firmware-on-nrf52480-dongle) for building and flashing the OpenThread RCP firmware on it.
3. Start only the OTBR agent (the GUI is not required):
```bash
sudo snap start openthread-border-router.otbr-agent
```

For more details on running the snap and logging, refer [here](https://github.com/canonical/openthread-border-router-snap#run).

#### 4. Pair the Matter Thread Lighting Device via Chip Tool on Machine A
Follow
[Commission and control a Matter Thread device via the OTBR Snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#pair-the-thread-lighting-device)
to pair the Matter Thread lighting device with the Thread network.

#### 5. Control the Lighting Device via Chip-Tool on Machine A
Read the "Control the Device" section on [Commission and control a Matter Thread device via the OTBR Snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#control-the-device) 
to control the lighting device.

Upon successful execution, the lighting status should be updated to "1/on" or "0/off" in the lighting app's logs:
```bash
$ ./out/debug/chip-lighting-app --thread
...
CHIP:DMG: Received command for Endpoint=1 Cluster=0x0000_0006 Command=0x0000_0002
CHIP:ZCL: Toggle ep1 on/off from state 0 to 1
```

### References
- [Commission and control a Matter Thread device via the OTBR Snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap)
- [Setup OpenThread Border Router with nRF52840 Dongle](https://github.com/canonical/openthread-border-router-snap/wiki/Setup-OpenThread-Border-Router-with-nRF52840-Dongle)
- [Matter Linux Lighting Example](https://github.com/project-chip/connectedhomeip/tree/6b01cb977127eb8547ce66d5b627061dc2dd6c90/examples/lighting-app/linux#chip-linux-lighting-example)
- https://github.com/project-chip/connectedhomeip/issues/29738
- [OpenThread Border Router Snap](https://github.com/canonical/openthread-border-router-snap#readme)
  