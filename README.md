# Transcribe || Scribouillage

Scribouillage is **transcription** software intended for use on Linux, using **libvlc** and the **Infinity Foot Pedal**. 

Because the footpedal doesn't come with a Linux driver out of the box, I'm using [zserge](https://github.com/zserge)'s **HID** driver in order to access the device somewhere in `/dev/bus/usb/00?/???`. Therefore, one needs reading permissions on that file in order for the driver to work.

The `transcribe.go` package is a work in progress, it uses [libvlc](https://github.com/adrg/libvlc-go) bindings to play a recording and the footpedal then jumps forward and back. The GUI uses andlabs `ui` package, there aren't too many options for GUI and golang, and `ui` isn't *there* yet.
