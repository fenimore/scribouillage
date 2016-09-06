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

// #cgo pkg-config: libusb-1.0
// #include <libusb-1.0/libusb.h>
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

func init() {
	C.libusb_init(nil)
}

func main() {
	devs := Devices()
	for _, dev := range devs {
		fmt.Println(dev)
		s := C.libusb_get_device_address(dev)
		fmt.Println(s)
	}
}

func Devices() []*C.libusb_device {
	//result := make(C.struct_libusb_device, 0, 10)
	var devices **C.struct_libusb_device // a list?
	count := C.libusb_get_device_list(nil, &devices)
	if count < 0 {
		return nil
	}
	defer C.libusb_free_device_list(devices, 1)

	return getDeviceList(devices, count)
}

func getDeviceList(devices **C.struct_libusb_device, cnt C.ssize_t) []*C.libusb_device {
	var list []*C.libusb_device
	*(*reflect.SliceHeader)(unsafe.Pointer(&list)) = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(devices)),
		Len:  int(cnt),
		Cap:  int(cnt),
	}
	return list
}
