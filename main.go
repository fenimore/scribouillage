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
import "os/exec"

func main() {
	file, err := os.Open("/dev/usb/hiddev0")
	if err != nil {
		fmt.Println(err)
	}
	data := make([]byte, 24) // Buffer for reading file

	cmd := exec.Command("xte", "key XF86AudioPlay")

	for {
		_, err := file.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		switch byte(1) {
		case data[4]:
			fmt.Println("Left")
		case data[12]:
			fmt.Println("Center")
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		case data[20]:
			fmt.Println("Right")
		}
	}
}
