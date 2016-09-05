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
	"fmt"
	"github.com/GeertJohan/go.hid"
)

func main() {
	_, err := hid.Open(0x05f3, 0x00ff, "")
	if err != nil {
		fmt.Println(err)
	}
	//defer pedal.Close()
}
