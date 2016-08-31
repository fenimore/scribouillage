// Fenimore Love | Public Domain
// See Python2 (yuck) driver https://ubuntuforums.org/showthread.php?t=2232673
// Also, some info about Infinity pedal: Vendor id = 05f3
//                                       Product id = 00ff
// My infinity pedal outputs to a file in /dev/usb hiddev0
// This, however, could change (to say, hiddev1?)
package main

import "os"
import "fmt"

// import "reflect"

func main() {
	file, err := os.Open("/dev/usb/hiddev0")
	if err != nil {
		fmt.Println(err)
	}
	data := make([]byte, 24)
	for {
		count, err := file.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("read %d bytes: %q\n\n", count, data[:count])
		//fmt.Printf("%d", reflect.ValueOf(data[2]))
		// TODO: should I use %b, %c or %o or %q
		fmt.Printf("Left %b, Center %b, Right %b", data[4], data[12], data[20])
	}
}
