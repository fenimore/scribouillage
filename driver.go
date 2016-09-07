// Public Domain | Fenimore Love
//
// Inspired by Python2 driver https://ubuntuforums.org/
//                           showthread.php?t=2232673
//
// Also, some info about Infinity pedal: Vendor id = 05f3
//                                       Product id = 00ff
//
// My infinity pedal outputs to a file in /dev/usb hiddev0
// This, however, could change (to say, hiddev1?)
//
// Like the python code, use the linux program xte to simulate
// keystrokes. For archlinux, this is found in the xautomation package.
package main

import (
	"bytes"
	"fmt"
	"github.com/zserge/hid"
	"time"
)

// Driver takes Infinity Foot Pedal and Input reads four different reports:
// 1. 0000 - Pedal released
// 2. 0100 - Left Pedal Pressed
// 3. 0200 - Middle Pedal Pressed
// 4. 0300 - Left + Middle Pressed
// 5. 0400 - Right Pedal Pressed
// 6. 0500 - Left + Right Pressed
// 7. 0600 - Middle + Right Pressed
// 8. 0700 - All Pedals Pressed

var (
	left    = []byte{1, 0} // "\x01\x00"
	right   = []byte{4, 0}
	middle  = []byte{2, 0}
	release = []byte{0, 0}
)

func main() {

	//REL_PED, _ = hex.DecodeString("0000")
	target := "05f3:00ff:0120"
	var dev hid.Device
	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		d := fmt.Sprintf("%04x:%04x:%04x", info.Vendor, info.Product, info.Revision)
		if d == target {
			dev = device
		}
	})
	err := dev.Open()
	// logs driver disconnect failed: -1 no data available ??
	if err != nil {
		fmt.Printf("Open Error: %s, Check Privileges\n", err)
	}
	defer dev.Close()

	for {
		buf, err := dev.Read(-1, 1*time.Second)
		if err == nil {
			// otherwise, get err 'connection timed out'
			if bytes.Equal(left, buf) {
				fmt.Println("Press: Left")
			} else if bytes.Equal(right, buf) {
				fmt.Println("Press: Right")
			} else if bytes.Equal(middle, buf) {
				fmt.Println("Press: Middle")
			} else if bytes.Equal(release, buf) {
				fmt.Println("Release")
			} else {
				// 0600, 0300, 0700 etc
				fmt.Println("Other Input")
			}
		}
	}
}
