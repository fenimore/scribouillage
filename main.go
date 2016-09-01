// Fenimore Love | Public Domain
// See Python2 driver https://ubuntuforums.org/
//                           showthread.php?t=2232673
// Also, some info about Infinity pedal: Vendor id = 05f3
//                                       Product id = 00ff
// My infinity pedal outputs to a file in /dev/usb hiddev0
// This, however, could change (to say, hiddev1?)
// Notes on fmt Formating, %b is base 2 %o is base 8.
// Looks like the data is given in unsigned base 8 ints. uint8
// From 24 byte buffer, bytes 5, 13, and 21 are important.
// Like the python code, use the linux program xte to simulate
// keystrokes. For archlinux, this is found in the xautomation package.
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
	args := make([]string, 3)
	args[0] = "xte"
	args[1] = "key"
	args[2] = "XF86AudioPlay"

	for {
		_, err := file.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		switch byte(1) {
		case data[4]:
			fmt.Println("Left")
			cmd := exec.Command("xte", "key XF86AudioPrev")
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		case data[12]:
			fmt.Println("Center")
			cmd := exec.Command("xte", "key XF86AudioPlay")
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		case data[20]:
			fmt.Println("Right")
			cmd := exec.Command("xte", "key XF86AudioNext")
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
