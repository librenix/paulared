package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/jaypipes/ghw"
)

type Hardware uint

type Result uint

const (
	CPU Hardware = iota
	Motherboard
	Memory
	Disk
	GPU
	Auido
	Ethernet
	Wireless
	Bluetooth
)

const (
	Supported Result = iota
	Unsupported
	Warning
	Depreated
	Unknown
)

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

func Compatibility(hw Hardware) (Result, error) {
	var pciDev string
	info := Detect()

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
				pciDev = fmt.Sprintf("Chipest: " + device.String() + "\n")
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
				pciDev = fmt.Sprintf("%s:\t %s %s\n", devStr, vendor, device.Product.Name)
			}
		}
	} else {
		err := fmt.Errorf("Could not retrieve pci device\n")
		log.Fatal(err)
		return Unknown, err
	}

	// USB Devices
	ctx := gousb.NewContext()
	defer ctx.Close()
	devs, _ := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if desc.Class == 0x01 && desc.SubClass == 0x01 { // Audio
			fmt.Sprintf("(USB) Audio:\t%s", usbid.Classify(desc))
		} else if desc.Class == 0x02 && desc.SubClass == 0x06 { // Ethernet network
			fmt.Sprintf("(USB) Ethernet:\t%s", usbid.Classify(desc))
		} else if desc.Class == 0x0e && desc.SubClass == 0x01 { // Video
			fmt.Sprintf("(USB) Video:\t%s", usbid.Classify(desc))
		} else if desc.Class == 0xe0 && desc.SubClass == 0x01 {
			if desc.Protocol == 0x01 { // Bluetooth
				fmt.Sprintf("(USB) Bluetooth:\t%s\n", usbid.Classify(desc))
			} else if desc.Protocol == 0x02 { // Wireless
				fmt.Sprintf("(USB) Wireless:\t%s", usbid.Classify(desc))
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

	// Motherboard
	//baseboard := hw.Baseboard
	switch hw {
	case Motherboard:
		product := info.Product
		vendor := strings.Fields(product.Vendor)[0]
		fmt.Sprintf("Host:\t %s %s\n", vendor, product.Name)
		return Supported, nil
		// CPU
	case CPU:
		cpu := info.CPU
		fmt.Sprintf("CPU:\t %s\n", cpu.Processors[0].Model)
		return Supported, nil
	// Memory
	case Memory:
		mem := info.Memory
		phys := mem.TotalPhysicalBytes
		size := phys / 1024 / 1024
		if size < 2048 {
			fmt.Sprintf("Memory:\t %dMB\n", size)
			return Unsupported, nil
		} else {
			fmt.Sprintf("Memory:\t %dMB\n", size)
			return Supported, nil
		}
	// Hard Disk
	case Disk:
		block := info.Block
		if len(block.Disks) > 0 {
			for _, disk := range block.Disks {
				size := disk.SizeBytes / 1024 / 1024 / 1024
				unitStr := ""
				if !disk.IsRemovable && disk.DriveType.String() == "HDD" || disk.DriveType.String() == "SSD" {
					if size > 1024 {
						size = size / 1024
						unitStr = "TB"
					} else {
						unitStr = "GB"
					}
				}
				fmt.Sprintf("Disk:\t %s %s (%s %s, %d%s)\n", disk.Vendor, disk.Model, disk.StorageController.String(), disk.DriveType.String(), size, unitStr)
			}
			return Supported, nil
		} else {
			err := fmt.Errorf("Can't find any hard disk\n")
			log.Fatal(err)
			return Unknown, err
		}
	// GPU
	case GPU:
		gpu := info.GPU
		for _, card := range gpu.GraphicsCards {
			vendor := strings.Fields(card.DeviceInfo.Vendor.Name)[0]
			fmt.Sprintf("GPU:\t %s %s\n", vendor, card.DeviceInfo.Product.Name)
		}
		return Supported, nil
	// Ethernet
	case Ethernet:
		fmt.Sprintln(pciDev)
		return Unknown, nil
	// Wireless
	case Wireless:
		fmt.Sprintln(pciDev)
		return Unknown, nil
	// Bluetooth
	case Bluetooth:
		fmt.Sprintln(pciDev)
		return Unknown, nil
	default:
		return Unknown, nil
	}
}
