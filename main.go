// Fenimore Love | Public Domain
// See Python2 (yuck) driver https://ubuntuforums.org/showthread.php?t=2232673
// Also, some info about Infinity pedal: Vendor id = 05f3
//                                       Product id = 00ff
// My infinity pedal outputs to a file in /dev/usb hiddev0
// This, however, could change (to say, hiddev1?)
// Notes on fmt Formating, %b is base 2 %o is base 8.
// Looks like the data is given in unsigned base 8 ints. uint8
// From 24 byte buffer, bytes 5, 13, and 21 are important.

package main

import "os"
import "fmt"

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
		//fmt.Printf("read %d bytes: %q\n\n", count, data[:count])
		fmt.Printf("Left %b, Center %b, Right %b", data[4], data[12], data[20])
	}
}
