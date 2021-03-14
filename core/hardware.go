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

type Processor uint

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

type Status uint

const (
	Supported Status = iota
	Unsupported
	Warning
	Depreated
	Unknown
)

type Result struct {
	Model  []string
	Status Status
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

func Compatibility(i interface{}) Result {
	var model []string
	status := Unknown

	switch hw := i.(type) {
	// CPU
	case ghw.CPUInfo:
		// Get cpu info from CPUID
		cpu := cpuid.CPU
		model = append(model, cpu.BrandName)
		if !cpu.Supports(cpuid.SSE3) {
			status = Unsupported
		}

		// If not virtual machine
		if cpu.VM() {
			log.Println("Virtual machine %s detected.", cpu.VendorString)
			status = Supported
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
						status = Supported
					// Xeon
					case 29: // Penryn (0x1d)
					case 46: // Nehalem (0x2e)
					case 44: // Westmere (0x2c)
					case 47: // Westmere (0x2f)
					case 45: // Sandy Bridge-E (0x2d)
					case 63: // Haswell-E (0x3f)
					case 85: // Skylake-X/W & Cascade Lake-X/W (0x55)
						status = Warning
					case 28: // Atom (0x1c)
					case 38: // Atom (0x26)
					case 54: // Atom (0x36)
					case 122: // Atom (0x7a)
					case 134: // Atom (0x86)
						status = Unsupported
					default:
						status = Unknown
					}
				} else {
					status = Unsupported
				}
			} else if cpu.IsVendor(cpuid.AMD) {
				// For vanilla, only 15h (Family 21+) and above are supported
				if cpu.Family >= 21 {
					status = Warning
				} else {
					status = Unsupported
				}
			} else {
				status = Unsupported
			}
		}
		break
	// Motherboard
	case ghw.ProductInfo:
		vendor := strings.Fields(hw.Vendor)[0]
		status = Supported
		model = append(model, vendor+hw.Name)
		break
	// Memory
	case ghw.MemoryInfo:
		phys := hw.TotalPhysicalBytes
		size := phys / 1024 / 1024
		if size < 2048 {
			status = Unsupported
		} else {
			status = Supported
		}
		model = append(model, fmt.Sprintf("%dMB\n", size))
		break
	// Hard Disk
	case ghw.BlockInfo:
		if len(hw.Disks) > 0 {
			for _, disk := range hw.Disks {
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
				model = append(model, fmt.Sprintf("%s %s (%s %s, %d%s)\n", disk.Vendor, disk.Model, disk.StorageController.String(), disk.DriveType.String(), size, unitStr))
			}
			status = Supported
		} else {
			err := fmt.Errorf("Can't find any hard disk\n")
			log.Fatal(err)
			status = Unknown
		}
		break
	// Graphics card
	case ghw.GPUInfo:
		for _, card := range hw.GraphicsCards {
			vendor := strings.Fields(card.DeviceInfo.Vendor.Name)[0]
			model = append(model, fmt.Sprintf("%s %s\n", vendor, card.DeviceInfo.Product.Name))
		}
		status = Supported
		break
	// PCI device
	case ghw.PCIInfo:
		devices := hw.ListDevices()
		if len(devices) > 0 {
			for _, device := range devices {
				// Class 03 = Display controller (GPU), skip
				var devStr string
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
					model = append(model, fmt.Sprintf("Chipest: "+device.String()+"\n"))
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
					model = append(model, fmt.Sprintf("%s %s\n", vendor, device.Product.Name))
				}
			}
		} else {
			err := fmt.Errorf("Could not retrieve pci device\n")
			log.Fatal(err)
			status = Unknown
		}
		break
	// USB
	case gousb.Context:
		ctx := hw

		defer ctx.Close()
		devs, _ := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
			var devStr string

			if desc.Class == 0x01 && desc.SubClass == 0x01 { // Audio
				devStr = "Audio"
			} else if desc.Class == 0x02 && desc.SubClass == 0x06 { // Ethernet network
				devStr = "Ethernet"
			} else if desc.Class == 0x0e && desc.SubClass == 0x01 { // Video
				devStr = "Video"
			} else if desc.Class == 0xe0 && desc.SubClass == 0x01 {
				if desc.Protocol == 0x01 { // Bluetooth
					devStr = "Bluetooth"
				} else if desc.Protocol == 0x02 { // Wireless
					devStr = "Wireless"
				}

				if devStr != "" {
					model = append(model, fmt.Sprintf("(USB) Wireless:\t%s", usbid.Classify(desc)))
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

		break
	}

	res := Result{model, status}
	return res
}
