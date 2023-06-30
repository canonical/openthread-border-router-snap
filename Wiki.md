# Commission and control a Matter Thread device via the OTBR Snap

### Setup environment:
- Ubuntu Desktop 23.04

- nRF52840 dongle (A) with RCP firmware as an RCP

- nRF52840 dongle (B) with the Thread lighting app running on it

## Installation and setup of OTBR Snap:
1. Install the OpenThread Border Router (OTBR) snap:
```bash
snap install openthread-border-router
```

2. Configure the snap by updating the following settings based on your system:
```bash
$ sudo snap get openthread-border-router
Key        Value
infra-if   eno1
radio-url  spinel+hdlc+uart:///dev/ttyACM0
thread-if  wpan0

$ sudo snap set openthread-border-router infra-if=eno1
```

3. Start the OTBR snap:
```bash
snap start openthread-border-router
```

## Forming the OTBR Network:
1. Connect the RCP dongle (A) to a USB port
2. Use the CTL tool to initialize the Thread network:
```bash
sudo openthread-border-router.ot-ctl
> dataset init new
Done
> dataset commit active
Done
> ifconfig up
Done
> thread start
Done
```
Alternatively, these steps could also be performed in the GUI at [https://localhost](https://localhost). 
Please refer to the instructions [here](https://openthread.io/guides/border-router/web-gui.md) to configure and form, join, or check the status of a Thread network using GUI.

## Discovering and pairing the Thread lighting device into the OTBR network
Make sure the Thread lighting app is running on dongle (B). Flashing and running the Thread lighting app on the device is beyond the scope of this guide. 
Please refer to [here](https://github.com/project-chip/connectedhomeip/tree/master/examples/lighting-app/nrfconnect). 

1. Obtaining the OTBR operational dataset (OTBR network's credentials):
```bash
sudo openthread-border-router.ot-ctl
> dataset active -x
0e08...f7f8
Done
```

2. Install and setup the chip-tool snap:
```bash
sudo snap install chip-tool
# Connect the avahi-observe interface to allow DNS-SD based discovery
sudo snap connect chip-tool:avahi-observe
# Connect the bluez interface for device discovery over Bluetooth Low Energy (BLE)
sudo snap connect chip-tool:bluez
```

3. Discovering and pairing the Thread lighting device into the OTBR network over Bluetooth LE:
```bash
sudo chip-tool pairing ble-thread 110 hex:0e08...f7f8 20202021 3840
```
where:
- `110` is the assigned node ID for the app.
- `0e08...f7f8` is the truncated Thread network credential operational dataset for readability.
- `20202021` is the PIN code set on the app.
- `3840` is the discriminator ID.

## Controlling the Thread lighting device
Toggle the device on or off using the following command:
```bash
sudo chip-tool onoff toggle 110 1
```
where:

-   `onoff` is the matter cluster name
-   `on`/`off`/`toggle` is the command name.
-   `110` is the node id of the app assigned during the commissioning
-   `1` is the endpoint of the configured device

Upon successful execution, the green LED on the dongle (B) will turn on or off. 
Additionally, two connected Thread nodes (one leader and one child) on the OpenThread GUI could be observed.

## Reference

- https://developers.home.google.com/matter/vendors/nordic-semiconductor
- https://github.com/project-chip/connectedhomeip/tree/master/examples/lighting-app/nrfconnect
- https://github.com/project-chip/connectedhomeip/blob/master/docs/guides/chip_tool_guide.md#using-chip-tool-for-matter-device-testing
