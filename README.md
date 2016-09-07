# Transcribe

So this works for taking an **infinity foot pedal** (or others?). It defaults with maps to the buttons Previous, Play, and Next. ~~This is for Linux and relies on a Linux utility `xte` (for Arch users, check out the xautomation package).~~ 

This `driver.go` in the Public Domain. ~~It's a pedal input mapper using `xte`.~~ `transcribe.go` is licensed under GPLv3. 

zserge's HID driver accesses the device somewhere in `/dev/bus/usb/00?/???`, one needs reading permissions in order for the 'driver' to work.

The `transcribe.go` package is a (quite doomed) work in progress, using the `libvlc` bindings to play a recording; the footpedal then jumps forward and back according to user settings.
