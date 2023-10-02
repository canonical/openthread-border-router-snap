## How to Set Up a Matter Thread Lighting Device on Nordic nRF52840 Dongle

This guide demonstrates how to build and flash a Matter Thread lighting application firmware onto a Nordic Semiconductor nRF52840 Dongle, 
thereby transforming the dongle into a Matter Thread device.

Complete the following steps to build and flash the firmware:

### Environment Setup
- Ubuntu 23.04 OS
- nRF52840 dongle

### Version tested in this guide
- Nordic Semiconductor's nRF Connect for Desktop: 4.2.0
- Nordic Semiconductor's nRF Connect SDK: 2.4.2
- [Matter](https://github.com/project-chip/connectedhomeip) SDK: 1.1.0.1

### 1. Setting up the Build Environment for nRF Connect SDK and Matter SDK

First, follow [Using native shell for setup](https://github.com/project-chip/connectedhomeip/tree/v1.1.0.1/examples/lighting-app/nrfconnect#using-native-shell-for-setup)
guide's steps 1 and 2 to install `nRF command-line tools`, `GN meta-build system`, and `nRF Connect SDK`.

Once nRF Desktop is installed, change its permissions, install dependencies, and run it:
```bash
chmod +x nrfconnect-4.2.0-x86_64.appimage
sudo apt install libfuse2
./nrfconnect-4.2.0-x86_64.appimage
```

In nRF Connect for Desktop, install the latest `toolchain manager` and `nRF Connect SDK v2.4.2`. 
After installations, click the down arrow next to the nRF Connect SDK version you installed and select `Open Terminal` or `Open Bash`.

### 2. Installing and Activating Matter SDK:
Install the Matter SDK with a shallow clone for the Linux platform:
```bash
cd ~
git clone https://github.com/project-chip/connectedhomeip.git --depth=1 --branch=v1.1.0.1
cd ~/connectedhomeip/
scripts/checkout_submodules.py --shallow --platform linux
git submodule update --init
```

Activate Matter build environment:
```bash
cd ~/connectedhomeip/

# Skip all non-core Python requirements:
sed -i '/^-r requirements/ s/./#&/' ./scripts/setup/requirements.txt
source scripts/activate.sh
```
### 3. Building the Matter Lighting Example Firmware:
By default, support for DFU (Device Firmware Upgrade) using Matter OTA is enabled.
To ensure successful firmware building for the dongle, utilize the following command, considering that nRF52840 dongle doesn't have the external flash required for DFU:
```bash
cd examples/lighting-app/nrfconnect
west build -b nrf52840dongle_nrf52840 -- -DCONF_FILE=./prj_no_dfu.conf
```

The output `zephyr.hex` file will be available at: `~/connectedhomeip/examples/lighting-app/nrfconnect/build/zephyr/zephyr.hex`.

### 4. Flashing the Firmware on the Dongle using nRF Desktop:

Refer to the "Programming the nRF52840 Dongle" section in [nRF Connect Programmer](https://infocenter.nordicsemi.com/index.jsp?topic=%2Fug_nc_programmer%2FUG%2Fnrf_connect_programmer%2Fncp_introduction.html) 
guide to flash the firmware on the dongle.

If you encounter the issue "Could not properly detect an unknown device. Please make sure you have nrf-udev installed," 
then install nrf-udev from [here](https://github.com/NordicSemiconductor/nrf-udev).

Upon successful flashing, logs will display:
```bash
Uploading image through SDFU: 100%
All DFU images have been written to the target device
Target device closed
```

### 5. Running the Matter Thread Lighting Device:
Connect the dongle to a USB port, then use `minicom` to run the lighting app:
```bash
sudo apt install minicom
sudo minicom -D /dev/ttyACM0 -b 115200
```
### Further Reading: 

- [nRF52840 dongle: Erasing persisted data of Matter Lighting App](https://github.com/canonical/openthread-border-router-snap/wiki/nRF52840-dongle:-Erasing-persisted-data-of-Matter-Lighting-App)
- [Commission and control a Matter Thread device via the OTBR Snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap)
- [Setup OpenThread Border Router with nRF52840 Dongle](https://github.com/canonical/openthread-border-router-snap/wiki/Setup-OpenThread-Border-Router-with-nRF52840-Dongle)
  
### References
- [Matter nRF Connect Lighting Example Application](https://github.com/project-chip/connectedhomeip/tree/v1.1.0.1/examples/lighting-app/nrfconnect)
- [nRF Connect Programmer](https://infocenter.nordicsemi.com/index.jsp?topic=%2Fug_nc_programmer%2FUG%2Fnrf_connect_programmer%2Fncp_introduction.html)
- [Using CLI in nRF Connect examples - minicom](https://github.com/project-chip/connectedhomeip/blob/master/docs/guides/nrfconnect_examples_cli.md#accessing-the-cli-console)
- [Disable DFU for nRF52840 dongle](https://github.com/project-chip/connectedhomeip/issues/24142#issuecomment-1359387829)
- [west config issue](https://github.com/project-chip/connectedhomeip/issues/6882#issuecomment-855862401)
