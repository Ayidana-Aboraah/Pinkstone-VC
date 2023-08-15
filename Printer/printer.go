package Printer

import (
	"strings"

	"github.com/karalabe/usb"
)

type Printer struct {
	device usb.Device
	Active bool
}

func FindUSBDevices() (devs []usb.DeviceInfo, out []int, err error) {
	devs, err = usb.Enumerate(0, 0)
	for i := range devs {
		out = append(out, i)
	}
	return
}

func Init(dev usb.DeviceInfo) (printer Printer, err error) {
	printer.device, err = dev.Open()
	printer.Active = true
	return
}

func (p Printer) Close() {
	p.device.Close()
}

func (p Printer) Print(out string) {
	p.device.Write([]byte(out))
}

func ConvertForPrinter(in string) string {
	in = string([]byte{0x1B, '@'}) + in // Adding intialisation to the beginnning

	// check for byte(183) and convert the line into Underline
	// in = strings.replaceAll(in, string([]byte{183}), string([]byte{0x1B, }))

	in = strings.ReplaceAll(in, "\n", string([]byte{0x1B, 'd', 1})) // replacing /n with line feed

	in = string([]byte{0x1B, 'V', 66}) + in // add cut command

	return in
}
