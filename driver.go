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
	"encoding/hex"
	"fmt"
	"github.com/zserge/hid"
	"time"
)

func main() {
	target := "05f3:00ff:0120"
	var dev hid.Device
	hid.UsbWalk(func(device hid.Device) {
		info := device.Info()
		d := fmt.Sprintf("%04x:%04x:%04x", info.Vendor, info.Product, info.Revision)
		if d == target {
			dev = device
		}
	})
	fmt.Println(dev)
	err := dev.Open()
	if err != nil {
		fmt.Printf("Open Error: %s, Check Privileges\n", err)
	}
	defer dev.Close()

	if report, err := dev.HIDReport(); err != nil {
		fmt.Println("HID report error:", err)
		return
	} else {
		fmt.Println("HID report", hex.EncodeToString(report))
	}

	for {
		if buf, err := dev.Read(-1, 1*time.Second); err == nil {
			fmt.Println("Input report:  ", hex.EncodeToString(buf))
			decoded, err := hex.DecodeString("0000")
			if err != nil {
				fmt.Println(err)
			}
			if bytes.Equal(decoded, buf) {
				fmt.Println("Yes")
			}
		}
	}
}
