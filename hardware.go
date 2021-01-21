package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/jaypipes/ghw"
	"github.com/muesli/termenv"
)

func ColorString(s string, c string) string {
	p := termenv.ColorProfile()
	style := termenv.String(s).Foreground(p.Color(c)).Bold()
	return fmt.Sprintln(style)
}

func Detect() *ghw.HostInfo {
	host, err := ghw.Host()
	if err != nil {
		log.Fatal("Failed to get hardware info", err)
	} else {
		log.Println(host.YAMLString())
	}
	// Save hardware info to json file
	b := []byte(host.JSONString(true))
	cache, err := os.UserConfigDir()
	path := cache + "/paulared/"
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir|0755)
	}
	f, err := os.OpenFile(path+"hardware.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	} else {
		_, err := f.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	return host
}

func Compatibility() (string, error) {
	var b strings.Builder
	hw := Detect()

	// Motherboard
	//baseboard := hw.Baseboard
	product := hw.Product
	vendor := strings.Fields(product.Vendor)[0]
	b.WriteString(ColorString(fmt.Sprintf("Host:\t %s %s\n", vendor, product.Name), "76"))

	// CPU
	cpu := hw.CPU
	b.WriteString(ColorString(fmt.Sprintf("CPU:\t %s\n", cpu.Processors[0].Model), "76"))

	// Memory
	mem := hw.Memory
	phys := mem.TotalPhysicalBytes
	size := phys / 1024 / 1024
	if size < 2048 {
		b.WriteString(ColorString(fmt.Sprintf("Memory:\t %dMB\n", size), "196"))
	} else {
		b.WriteString(ColorString(fmt.Sprintf("Memory:\t %dMB\n", size), "76"))
	}

	// Hard Disk
	block := hw.Block
	if len(block.Disks) > 0 {
		for _, disk := range block.Disks {
			if !disk.IsRemovable && disk.DriveType.String() == "HDD" || disk.DriveType.String() == "SSD" {
				size := disk.SizeBytes / 1024 / 1024 / 1024
				unitStr := ""
				if size > 1024 {
					size = size / 1024
					unitStr = "TB"
				} else {
					unitStr = "GB"
				}
				b.WriteString(ColorString(fmt.Sprintf("Disk:\t %s %s (%s %s, %d%s)\n", disk.Vendor, disk.Model, disk.StorageController.String(), disk.DriveType.String(), size, unitStr), "76"))
			}
		}
	} else {
		log.Fatal("Can't find any hard disk\n")
	}

	// GPU
	gpu := hw.GPU
	for _, card := range gpu.GraphicsCards {
		vendor := strings.Fields(card.DeviceInfo.Vendor.Name)[0]
		b.WriteString(ColorString(fmt.Sprintf("GPU:\t %s %s\n", vendor, card.DeviceInfo.Product.Name), "76"))
	}

	// PCI Devices
	pci, err := ghw.PCI()
	if err != nil {
		log.Fatal("Failed to get pci info", err)
	}

	devices := pci.ListDevices()
	if len(devices) > 0 {
		for _, device := range devices {
			// Class 03 = Display controller (GPU), skip
			devStr := ""
			vendor := strings.Fields(device.Vendor.Name)[0]
			if device.Class.ID == "02" { // Network controller
				if device.Subclass.ID == "00" {
					devStr = "Ethernet"
				} else if device.Subclass.ID == "80" {
					devStr = "Wireless"
				}
			} else if device.Class.ID == "04" && device.Subclass.ID == "03" { // Multimedia controller
				devStr = "Audio"
			} else if device.Class.ID == "06" && device.Subclass.ID == "01" { // Southbridge
				ColorString(fmt.Sprintf("Chipest: "+device.String()+"\n"), "184")
			} else if device.Class.ID == "09" { // Input controller
				if device.Subclass.ID == "00" {
					devStr = "Keyboard"
				} else if device.Subclass.ID == "02" {
					devStr = "Mouse"
				}
			} else if device.Class.ID == "0d" && device.Subclass.ID == "80" { // Network controller
				devStr = "Wireless"
			}
			if devStr != "" {
				b.WriteString(ColorString(fmt.Sprintf("%s:\t %s %s\n", devStr, vendor, device.Product.Name), "76"))
			}
		}
	} else {
		log.Fatal("Could not retrieve pci device\n")
	}

	// USB Devices
	ctx := gousb.NewContext()
	defer ctx.Close()
	devs, _ := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if desc.Class == 0x01 && desc.SubClass == 0x01 { // Audio
			ColorString(fmt.Sprintf("(USB) Audio:\t%s", usbid.Classify(desc)), "76")
		} else if desc.Class == 0x02 && desc.SubClass == 0x06 { // Ethernet network
			ColorString(fmt.Sprintf("(USB) Ethernet:\t%s", usbid.Classify(desc)), "76")
		} else if desc.Class == 0x0e && desc.SubClass == 0x01 { // Video
			ColorString(fmt.Sprintf("(USB) Video:\t%s", usbid.Classify(desc)), "76")
		} else if desc.Class == 0xe0 && desc.SubClass == 0x01 {
			if desc.Protocol == 0x01 { // Bluetooth
				b.WriteString(ColorString(fmt.Sprintf("(USB) Bluetooth:\t%s\n", usbid.Classify(desc)), "76"))
			} else if desc.Protocol == 0x02 { // Wireless
				b.WriteString(ColorString(fmt.Sprintf("(USB) Wireless:\t%s", usbid.Classify(desc)), "76"))
			}
		}
		return false
	})

	for _, dev := range devs {
		_ = dev
	}

	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	return b.String(), nil
}
