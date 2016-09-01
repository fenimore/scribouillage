# pedal-driver

So this works for taking an **infinity foot pedal** (or others?). It defaults with maps to the buttons Previous, Play, and Next. This is for Linux and relies on a Linux utility `xte` (for Arch users, check out the xautomation package).

This driver, if it can be called that (it's quite simple, and hacky, and it isn't really a driver...), is in the Public Domain. It's more of a pedal input mapper using `xte`...

It can either be run using sudo, or it a user can be added to a group with privilages to read `/dev/usb/hiddev0`.

For **transcription**, the pedals ought to *just work* for playing, skipping, and returning on VLC (or other music players). So, for transcription, the default VLC hotkeys (and these can be set in the Preferences) are `Shift+Left/Alt+Left/Shift+Alt+Left` and `Shift+Right/Alt+Right/Shift+Alt+Right` for **short/medium/long jumps forward and backwards**, respectively.
