# Setup OpenThread Border Router with nRF52840 Dongle

This guide will help you set up and configure OTBR with Nordic nRF52840 Dongle. 
For additional information and troubleshooting, please refer to the provided references.

Before proceeding with setting up OTBR (Openthread Brouder Router) on your machine with Nordic Semiconductor nRF52840 Dongle, which acts as a Radio Co-Processor (RCP), 
it is essential to build and flash the RCP firmware onto the dongle.  To do this, follow these steps:

### Build and flash RCP firmware on nRF52480 dongle:

Following the instructions [here](https://developer.nordicsemi.com/nRF_Connect_SDK/doc/1.9.99-dev1/matter/openthread_rcp_nrf_dongle.html):

#### Debugging tips 

##### Finding the device name

During the setup process, ensure that the dongle device is connected and recognized as expected at `/dev/ttyACM0`. 
The `dmesg` command can help with this verification. 
Here's an example of how to use it: 
```
$ sudo dmesg -W
[23930.322598] usb 1-8: USB disconnect, device number 24
[23930.971312] usb 1-8: new full-speed USB device number 25 using xhci_hcd
[23931.130085] usb 1-8: New USB device found, idVendor=1915, idProduct=cafe, bcdDevice= 1.00
[23931.130099] usb 1-8: New USB device strings: Mfr=1, Product=2, SerialNumber=3
[23931.130104] usb 1-8: Product: nRF528xx OpenThread Device
[23931.130108] usb 1-8: Manufacturer: Nordic Semiconductor
[23931.130111] usb 1-8: SerialNumber: F54FE6398B9E
[23931.133681] cdc_acm 1-8:1.0: ttyACM0: USB ACM device
```

##### PIP version error

If you encounter an error related to pip version during the installation of nRF Util, such as:

```bash
  Could not find a version that satisfies the requirement pc_ble_driver_py>=0.14.0 (from nrfutil) (from versions: 0.1.0, 0.2.0, 0.3.0, 0.4.0, 0.5.0, 0.6.0, 0.6.1, 0.6.2, 0.8.0, 0.8.1, 0.9.0, 0.9.1, 0.10.0, 0.11.0, 0.11.1, 0.11.2, 0.11.3, 0.11.4)
No matching distribution found for pc_ble_driver_py>=0.14.0 (from nrfutil)
```

Update pip with the following commands:
```bash
python3 -m pip install -U pip
python3 -m pip install -U nrfutil
```

##### AttributeError
If you encounter an error related to AttributeError during the installation of nRF Util, such as:
```bash
AttributeError: 'dict' object has no attribute 'iteritems'
```

Install missing packages with the following commands:

```bash
 sudo apt-get -y install libusb-1.0-0-dev sed
 pip3 install click crcmod ecdsa intelhex libusb1 piccata protobuf pyserial pyyaml tqdm pc_ble_driver_py
 pip3 install -U --no-dependencies nrfutil
```

##### Unexpected Executed OP_CODE error

If you encounter an error during flashing, such as:

```bash
pc_ble_driver_py.exceptions.NordicSemiException: Unexpected Executed OP_CODE.
Expected: 0x02 Received: 0x7E
```

Verify the environment as described [here](https://github.com/NordicSemiconductor/pc-ble-driver-py/issues/29#issuecomment-317967858), 
or unplug the dongle from current USB port and plug it into another new USB port as documented [here](https://devzone.nordicsemi.com/f/nordic-q-a/82759/serial-dfu-is-failing-with-unexpected-executed-op_code).



### Configure and run OTBR snap:

Now that the RCP firmware has successfully flashed, it's time to configure and run the OTBR Snap. 
Follow these steps in this [guide](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#install-and-configure-the-otbr-snap)
to:

1. Build, install, and configure the snap
2. Connect the RCP dongle (A) to a USB port
3. Start the OTBR snap




### References:
- [Build and flash Nordic Semiconductor nRF52840](https://github.com/openthread/openthread/blob/main/src/posix/README.md#nordic-semiconductor-nrf52840)
- [Configuring OpenThread Radio Co-processor on nRF52840 Dongle](https://developer.nordicsemi.com/nRF_Connect_SDK/doc/1.9.99-dev1/matter/openthread_rcp_nrf_dongle.html)
- [Commission and control a Matter Thread device via the OTBR Snap](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap)
- [Mass Storage Device known issue](https://github.com/openthread/ot-nrf528xx/blob/main/src/nrf52840/README.md#mass-storage-device-known-issue)
- [pip version issue](https://github.com/crownstone/bluenet/issues/81#issuecomment-561257090)
- [AttributeError issue](https://github.com/makerdiary/nrf52840-mdk-usb-dongle/issues/56#issuecomment-1322301257)
- [Unexpected Executed OP_CODE issue](https://devzone.nordicsemi.com/f/nordic-q-a/82759/serial-dfu-is-failing-with-unexpected-executed-op_code)
- [Flashing environment setup](https://github.com/NordicSemiconductor/pc-ble-driver-py/issues/29#issuecomment-317967858)

