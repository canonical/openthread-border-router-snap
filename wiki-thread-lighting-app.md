## How to Set Up and Control a Matter Thread Lighting Device

This guide will walk you through the process of building, running, and controlling a Matter Thread lighting application over a Thread network. 

Here are step-by-step instructions:

#### Hardware Setup
- Machine A
  - nRF52840 dongle as RCP
  - Choose one of the following options:
    - PC running Ubuntu Classic 23.04 OS
    - Raspberry Pi 4B running Ubuntu Server 22.04 OS or Ubuntu Core 22 and LED
- Machine B
  - nRF52840 dongle as RCP
  - PC running Ubuntu Classic 23.04 OS
 
#### Environment Setup
- Machine A running
  - OTBR snap
  - Matter Thread ligting App
- Machine B running
  - OTBR snap
  - chip-tool snap

#### Versions Tested in This Guide
- Matter SDK: `6b01cb977127eb8547ce66d5b627061dc2dd6c90`
- RCP: TBA

#### 1. Form a Thread Network on Machine A
on Machine A, following this guide to [install and configure](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#install-and-configure-the-otbr-snap) 
the OTBR snap and [form a Thread network](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#form-a-thread-network) 
using the nRF52840 Dongle as the RCP. For details on hsetting up the RCP, 
pelase refer to [Build and flash RCP firmware on nRF52480 dongle](https://github.com/canonical/openthread-border-router-snap/wiki/Setup-OpenThread-Border-Router-with-nRF52840-Dongle#build-and-flash-rcp-firmware-on-nrf52480-dongle).

#### 2. Run a Matter Thread Lighting App on Machine B
There are two options to run a Matter Thread lighting application. 

One is native build lighting app from Matter SDK runing on native shell, the lighting status of on/off could be observe from logs.
the second option is using a packaged Matter Pi GPIO Commander snap to turn your Raspberry Pi into a Matter Thread lighting device, and actual connected LED coiuld be rurn on/off. Please choose the one suit your need most.

**Option 1**: Native Build Lighting App
You can choose this option to build and run the Matter Thread lighting app from the Matter SDK in the native shell. With this option, you can observe the lighting status (on/off) from the application's logs.

To get started, follow these steps:

- 1. Install the Matter SDK with a shallow clone for the Linux platform:
```
git clone https://github.com/project-chip/connectedhomeip.git --depth=1
cd connectedhomeip
git checkout 6b01cb977127eb8547ce66d5b627061dc2dd6c90
scripts/checkout_submodules.py --shallow --platform linux
```
- 2. Build the lighting app:
```
cd examples/lighting-app/linux
gn gen out/debug
ninja -C out/debug
```
- 3. Run the lighting app with Thread feature enabled:
```
./out/debug/chip-lighting-app --thread
```

**Option 2**: GPIO Commander Snap
Alternatively, you can use this option to set up the Matter Thread lighting application by installing the Matter Pi GPIO Commander snap. 
This turns your Raspberry Pi into a Matter Thread lighting device, allowing you to control the connected LED, turning it on or off.

Follow the instructions for installing the matter-pi-gpio-commander snap, 
provided at [Setup OpenThread Border Router with nRF52840 Dongle](https://github.com/canonical/matter-pi-gpio-commander/wiki/Setup-and-control-a-lighting-device#installation).

TBA: interface connection

#### 3. Run OTBR D-Bus Daemon on Machine B
Check the [install and configure the OTBR snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#install-and-configure-the-otbr-snap) section
to run OTBR dbus daemon, which is needed to support lighting app's Thread feature.

#### 4. Pair the Matter Thread Lighting Device via Chip Tool on Machine A
follow the pairing steps in 
[Commission and control a Matter Thread device via the OTBR Snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#pair-the-thread-lighting-device)
to pair the Matter Thread lighting device with the Thread network.

#### 5. Control the Lighting Device via Chip-Tool on Machine A
Read the "Control the Device" section on [Commission and control a Matter Thread device via the OTBR Snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#control-the-device) 
to control the lighting device.

If you are running the lighting app with **Option 1** (native lighting app), upon successful execution, the lighting status should be updated to "on" or "off" in the lighting app's logs. 

If you are using **Option 2** (GPIO Commander Snap), upon successful execution, the connected LED on the Raspberry Pi will turn on or off.
### References
- [Commission and control a Matter Thread device via the OTBR Snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap)
- [Setup OpenThread Border Router with nRF52840 Dongle](https://github.com/canonical/openthread-border-router-snap/wiki/Setup-OpenThread-Border-Router-with-nRF52840-Dongle)
- https://github.com/project-chip/connectedhomeip/issues/29738




