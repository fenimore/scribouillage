// Public Domain | Fenimore Love
//
// Also, some info about Infinity pedal: Vendor id = 05f3
//                                       Product id = 00ff
package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/zserge/hid"
)

// Driver takes Infinity Foot Pedal and Input reads four different reports:
// Reports are given in two byte slices, such as \x02\x00
// 1. 0000 - Pedal released
// 2. 0100 - Left Pedal Pressed
// 3. 0200 - Middle Pedal Pressed
// 4. 0300 - Left + Middle Pressed
// 5. 0400 - Right Pedal Pressed
// 6. 0500 - Left + Right Pressed
// 7. 0600 - Middle + Right Pressed
// 8. 0700 - All Pedals Pressed

const (
	left    = 1
	right   = 4
	middle  = 2
	release = 0
)

func main() {
	target := "05f3:00ff" // Vendor and Product ID for Infinity FootPedal
	var device hid.Device
	hid.UsbWalk(func(dev hid.Device) {
		info := dev.Info()
		d := fmt.Sprintf("%04x:%04x", info.Vendor, info.Product)
		if d == target {
			device = dev
		}
	})
	err := device.Open()
	// logs driver disconnect failed: -1 no data available ??
	if err != nil {
		fmt.Printf("Open Error: %s, Check Privileges\n", err)
	}
	defer device.Close()

	for {
		buf, err := device.Read(-1, 1*time.Second)
		if err == nil {
			// otherwise, get err 'connection timed out'
			switch binary.LittleEndian.Uint16(buf) {
			case left:
				fmt.Println("Press: Left")
			case right:
				fmt.Println("Press: Right")
			case middle:
				fmt.Println("Press: Middle")
			case release:
				fmt.Println("Release")
			default:
				// 0600, 0300, 0700 etc
				fmt.Println("Other Input")
			}
		}
	}
}
