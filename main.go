package main

import "os"
import "fmt"

// import "reflect"

func main() {
	file, err := os.Open("/dev/usb/hiddev0")
	if err != nil {
		fmt.Println(err)
	}
	data := make([]byte, 100)
	for {
		count, err := file.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("read %d bytes: %q\n\n", count, data[:count])
		fmt.Printf("%d", reflect.ValueOf(data[2]))
	}
}
