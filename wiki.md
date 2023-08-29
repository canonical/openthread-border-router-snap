# A How-To Guide for Setting Up OTBR with the Nordic nRF52840 Dongle

Before proceeding with setting up OTBR (Openthread Brouder Router) on your machine with Nordic nRF52840 Dongle, which acts as a Radio Co-Processor (RCP), 
it is essential to build and flash the RCP firmware onto the dongle.  To do this, follow these steps:

### Build and flash RCP firmware on nRF52480 dongle:

Following the instructions [here](https://developer.nordicsemi.com/nRF_Connect_SDK/doc/1.9.99-dev1/matter/openthread_rcp_nrf_dongle.html):

If you encounter an error related to pip during the installation of nRF Util, update pip with the following commands:
```
python3 -m pip install -U pip
python3 -m pip install -U nrfutil
```

### Configure and run OTBR snap:

Now that the RCP firmware has successfully flashed, it's time to configure and run the OTBR Snap. 
Follow these steps in this [guide](https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap#install-and-configure-the-otbr-snap)
to:

1. Build, install, and configure the snap
2. Connect the RCP dongle (A) to a USB port
3. Start the OTBR snap


### Development 
During the setup process, verify if the dongle device is connected and recognized as expected at `/dev/ttyACM0`. 
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

### References:
- https://github.com/openthread/openthread/blob/main/src/posix/README.md#nordic-semiconductor-nrf52840
- https://developer.nordicsemi.com/nRF_Connect_SDK/doc/1.9.99-dev1/matter/openthread_rcp_nrf_dongle.html
- https://github.com/canonical/openthread-border-router-snap/wiki/Commission-and-control-a-Matter-Thread-device-via-the-OTBR-Snap
- https://github.com/openthread/ot-nrf528xx/blob/main/src/nrf52840/README.md#mass-storage-device-known-issue
- https://github.com/crownstone/bluenet/issues/81#issuecomment-561257090
- https://github.com/makerdiary/nrf52840-mdk-usb-dongle/issues/56#issuecomment-1322301257

This how-to guide will help you set up and configure OTBR with Nordic nRF52840 Dongle. 
For additional information and troubleshooting, please refer to the provided references.
