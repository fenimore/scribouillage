# Transcribe

So this works for taking an **infinity foot pedal** (or others?). It defaults with maps to the buttons Previous, Play, and Next. This is for Linux and relies on a Linux utility `xte` (for Arch users, check out the xautomation package).

This `driver.go`, if it can be called a driver (it has nothing to do with drivers...), is in the Public Domain. It's a pedal input mapper using `xte`... `transcribe.go` is licensed under GPLv3.

It can either be run using sudo, or it a user can be added to a group with privilages to read `/dev/usb/hiddev0`, or one may chmod a+x `hiddev0` to grant reading privilages.

The `transcribe.go` package is a (quite doomed) work in progress, using the `libvlc` bindings to play a recording; the footpedal then jumps forward and back according to user settings.
