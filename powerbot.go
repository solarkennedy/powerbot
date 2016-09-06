package main

import (
	"encoding/binary"
	"fmt"
	"github.com/deadsy/libusb"
	_ "github.com/thoj/go-ircevent"
	"log"
	"os"
	"time"
)

const USBRQ_HID_SET_REPORT = 0x09
const USB_HID_REPORT_TYPE_FEATURE = 0x03

type ep_info struct {
	itf int
	ep  *libusb.Endpoint_Descriptor
}

func printUSBInfo(handle libusb.Device_Handle) {
	dev := libusb.Get_Device(handle)
	path := make([]byte, 8)
	path, err := libusb.Get_Port_Numbers(dev, path)
	if err != nil {
		log.Fatal(err)
	}
	dd, err := libusb.Get_Device_Descriptor(dev)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bus %03d Device %03d: ID %04x:%04x\n", libusb.Get_Bus_Number(dev), libusb.Get_Device_Address(dev), dd.IdVendor, dd.IdProduct)
	fmt.Printf("%v %d\n", path, libusb.Get_Port_Number(dev))
}

func codeToBytes(code int) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(code))
	fmt.Printf("DEBUG: code: %v, bytes: %v\n", code, bs)
	return bs
}

func writeCode(handle libusb.Device_Handle, code int) {
	sendByte(handle, 0x73)
	sendByte(handle, 0x10)
	bs := codeToBytes(code)
	for i := range bs {
		if i != 0 {
			sendByte(handle, bs[i])
		}
	}
	sendByte(handle, 0xa)
}

func sendByte(handle libusb.Device_Handle, thebyte byte) {
	fmt.Printf("Sending %v\n", thebyte)
	timeout := uint(6000)
	request_type := uint8(0x20)
	request := uint8(0x09)
	wvalue := uint16(0x300)
	windex := uint16(thebyte)
	//data := []byte{thebyte}
	//fmt.Printf("Sending index %#x\n", data)
	data := []byte{0}
	bytes, err := libusb.Control_Transfer(handle, request_type, request, wvalue, windex, data, timeout)
	fmt.Printf("Err: %v\n", err)
	fmt.Printf("Sent bytes %v\n", bytes)
}

func main() {
	var ctx libusb.Context
	err := libusb.Init(&ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer libusb.Exit(ctx)

	vid := uint16(0x16c0)
	pid := uint16(0x05df)
	fmt.Printf("Opening device %04X:%04X...\n", vid, pid)
	handle := libusb.Open_Device_With_VID_PID(nil, vid, pid)
	if handle == nil {
		fmt.Fprintf(os.Stderr, "  Failed.\n")
		return
	}
	defer libusb.Close(handle)
	printUSBInfo(handle)
	writeCode(handle, 95500)
	time.Sleep(time.Duration(1) * time.Second)
	//	fmt.Printf(" Done. Next!\n")
	writeCode(handle, 95491)
}
