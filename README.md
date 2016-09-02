# Transcribe

So this works for taking an **infinity foot pedal** (or others?). It defaults with maps to the buttons Previous, Play, and Next. This is for Linux and relies on a Linux utility `xte` (for Arch users, check out the xautomation package).

This `driver.go`, if it can be called a driver (it's quite simple, and hacky, and it isn't really a driver...), is in the Public Domain. It's more of a pedal input mapper using `xte`... `transcribe.go` is licensed under GPLv3.

It can either be run using sudo, or it a user can be added to a group with privilages to read `/dev/usb/hiddev0`, or one may chmod a+x `hiddev0` to grant reading privilages.

For **transcription**, the pedals ought to *just work* for playing, skipping, and returning on VLC (or other music players). So, for transcription, the default VLC hotkeys (and these can be set in the Preferences) are `Shift+Left/Alt+Left/Shift+Alt+Left` and `Shift+Right/Alt+Right/Shift+Alt+Right` for **short/medium/long jumps forward and backwards**, respectively. These keybindings won't be picked up by VLC unless the application window is active -- making the `driver.go` not so useful for actual transcription (that is, one should be typing on a word document elsewhere). See `transcribe.go` for the work around.

The `transcribe.go` package is a work in progress, using the `libvlc` bindings to play a recording; the footpedal then jumps forward and back according to user settings.
