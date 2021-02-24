package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/jaypipes/ghw"
	"github.com/klauspost/cpuid/v2"
)

type Hardware uint
type Processor uint
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
	// Intel
	Penryn Processor = iota
	Clarksfield
	Nehalem
	SandyBridge
	SandyBridge_E
	IvyBridge
	Haswell
	Haswell_E
	Boardwell
	Skylake
	Skylake_X
	KabyLake
	CoffeeLake
	IceLake
	TigerLake
	TigerLake_L
	CometLake
	CometLake_L

	// AMD
	Bulldozer
	Jaguar
	Zen
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

func Compatibility(info ghw.HostInfo, hw Hardware) (Result, error) {
	var pciDev string

	// PCI Devices
	pci := info.PCI
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
		// Get cpu info from CPUID
		cpu := cpuid.CPU
		fmt.Sprintf("CPU:\t %s\n", cpu.BrandName)
		if !cpu.Supports(cpuid.SSE3) {
			return Unsupported, nil
		}

		// If not virtual machine
		if cpu.VM() {
			log.Println("Virtual machine %s detected.", cpu.VendorString)
			return Supported, nil
		} else {
			if cpu.IsVendor(cpuid.Intel) {
				if cpu.Family == 6 {
					switch cpu.Model {
					// Core
					case 15: // Conroe (0x0f)
					case 23: // Penryn (0x17)
					case 26: // Clarksfield (0x1a)
					case 30: // Clarksfield (0x1e)
					case 37: // Sandy Bridge (0x25)
					case 42: // Sandy Bridge (0x2a)
					case 58: // Ivy Bridge (0x3a)
					case 60: // Haswell (0x3c)
					case 61: // Boardwell (0x3d)
					case 78: // Skylake (Laptop) (0x4e)
					case 94: // Skylake (Desktop) (0x5e)
					case 142: // Kaby Lake (Laptop) (0x8e)
					case 158: // Kaby Lake (Desktop) / Coffee Lake (0x9e)
					case 102: // Cannon Lake (0x66)
					case 126: // Ice Lake (0x7e)
					case 140: // Tiger Lake-L (0x8c)
					case 141: // Tiger Lake (0x8d)
					case 165: // Comet Lake (0xa5)
					case 166: // Comet Lake-L (0xa6)
						return Supported, nil
					// Xeon
					case 29: // Penryn (0x1d)
					case 46: // Nehalem (0x2e)
					case 44: // Westmere (0x2c)
					case 47: // Westmere (0x2f)
					case 45: // Sandy Bridge-E (0x2d)
					case 63: // Haswell-E (0x3f)
					case 85: // Skylake-X/W & Cascade Lake-X/W (0x55)
						return Warning, nil
					case 28: // Atom (0x1c)
					case 38: // Atom (0x26)
					case 54: // Atom (0x36)
					case 122: // Atom (0x7a)
					case 134: // Atom (0x86)
						return Unsupported, nil
					default:
						return Unknown, nil
					}
				} else {
					return Unsupported, nil
				}
			} else if cpu.IsVendor(cpuid.AMD) {
				// For vanilla, only 15h (Family 21+) and above are supported
				if cpu.Family >= 21 {
					return Warning, nil
				} else {
					return Unsupported, nil
				}
			} else {
				return Unsupported, nil
			}
		}
		return Unknown, nil
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
