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
	ghid "github.com/GeertJohan/go.hid"
	bhid "github.com/boombuler/hid"
	"github.com/flynn/hid"
	"strings"
)

func main() {
	devices, errs := hid.Devices()
	if errs != nil {
		fmt.Println(errs)
	}
	for elem := range devices {
		fmt.Println(elem)
	}
	//boombuler()
}

func boombuler() {
	devices := bhid.Devices()
	fmt.Printf("%T\n", devices)
	for elem := range devices {
		fmt.Println(elem.Path)
		if strings.HasPrefix(elem.Path, "05f3:00ff") {
			fmt.Println(elem.VendorId)
			fmt.Println("yes")
			d, err := elem.Open()
			// and then close
			if err != nil {
				fmt.Println(err)
				// Seems to always get -3 error
				// Insufficent Permission
			}
			fmt.Println(d)
		}
	}
}

func Geert() {
	ds, _ := ghid.Enumerate(0x0, 0x0)
	for _, d := range ds {
		fmt.Println(d.VendorId)
		fmt.Println(d.Path)
		dev, err := d.Device()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(dev)
	}
}
