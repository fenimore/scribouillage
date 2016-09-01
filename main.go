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
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	file, err := os.Open("/dev/usb/hiddev0")
	if err != nil {
		fmt.Println(err)
	}
	data := make([]byte, 24) // Buffer for reading file

	// Flags
	// TODO: Have flag decide distance of Jump
	tranFlag := flag.Bool("t", false, "For transcriptions")

	flag.Parse()
	if *tranFlag {
		fmt.Println("Using for transcription")
	}

	for {
		_, err := file.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		switch byte(1) {
		case data[4]:
			fmt.Println("Left")
			prev()
			//jumpBack()
		case data[12]:
			fmt.Println("Center")
			pause()
		case data[20]:
			fmt.Println("Right")
			next()
		}
	}
}

func next() {
	cmd := exec.Command("xte", "key XF86AudioNext")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func prev() {
	cmd := exec.Command("xte", "key XF86AudioPrev")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func pause() {
	cmd := exec.Command("xte", "key XF86AudioPlay")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// This is in the right track,
// But it totally doesn't work
func jumpBack() {
	cmd := exec.Command("xte", "keydown Alt_L key Right keyup Alt_L")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

}
