package main

import (
	"strings"

	"github.com/fatih/color"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/jaypipes/ghw"
)

func errMsg(msg string) {
	color.Red("Error: Failed to get %s info!\n", msg)
}

func check() {
	ok := color.New(color.FgGreen, color.Bold).PrintfFunc()

	// Motherboard
	product, err := ghw.Product()
	if err != nil {
		errMsg("motherboard")
	} else {
		vendor := strings.Fields(product.Vendor)[0]
		ok("Host:\t %s %s\n", vendor, product.Name)
	}

	// CPU
	cpu, err := ghw.CPU()
	if err != nil {
		errMsg("cpu")
	} else {
		ok("CPU:\t %s\n", cpu.Processors[0].Model)
	}

	// Memory
	mem, err := ghw.Memory()
	if err != nil {
		errMsg("memory")
	} else {
		phys := mem.TotalPhysicalBytes
		size := phys / 1024 / 1024
		if size < 2048 {
			color.Red("Memory:\t %dMB\n", size)
		} else {
			ok("Memory:\t %dMB\n", size)
		}
	}

	// Hard Disk
	block, err := ghw.Block()
	if err != nil {
		errMsg("hard disk")
	} else {
		if len(block.Disks) > 0 {
			for _, disk := range block.Disks {
				if disk.DriveType.String() == "HDD" || disk.DriveType.String() == "SSD" {
					size := disk.SizeBytes / 1024 / 1024 / 1024
					unitStr := ""
					if size > 1024 {
						size = size / 1024
						unitStr = "TB"
					} else {
						unitStr = "GB"
					}
					ok("Disk:\t %s %s (%s %s, %d%s)\n", disk.Vendor, disk.Model, disk.StorageController.String(), disk.DriveType.String(), size, unitStr)
				}
			}
		} else {
			color.Red("Can't found any hard disk\n")
		}
	}

	// GPU
	gpu, err := ghw.GPU()
	if err != nil {
		errMsg("gpu")
	} else {
		for _, card := range gpu.GraphicsCards {
			vendor := strings.Fields(card.DeviceInfo.Vendor.Name)[0]
			ok("GPU:\t %s %s\n", vendor, card.DeviceInfo.Product.Name)
		}
	}

	// PCI Devices
	pci, err := ghw.PCI()
	if err != nil {
		errMsg("pci devices")
	} else {
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
				} else if device.Class.ID == "0d" && device.Subclass.ID == "80" { // Network controller
					devStr = "Wireless"
				}
				if devStr != "" {
					ok("%s:\t %s %s [%s]\n", devStr, vendor, device.Product.Name, device.Subclass.ID)
				}
			}
		} else {
			color.Red("Could not retrieve pci device\n")
		}
	}

	// USB Devices
	ctx := gousb.NewContext()
	defer ctx.Close()
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if desc.Class == 0x01 && desc.SubClass == 0x01 { // Audio
			ok("(USB) Audio: \t%s", usbid.Classify(desc))
		} else if desc.Class == 0x02 && desc.SubClass == 0x06 { // Ethernet network
			ok("(USB) Ethernet: \t%s", usbid.Classify(desc))
		} else if desc.Class == 0x0e && desc.SubClass == 0x01 { // Video
			ok("(USB) Video: \t%s", usbid.Classify(desc))
		} else if desc.Class == 0xe0 && desc.SubClass == 0x01 {
			if desc.Protocol == 0x01 { // Bluetooth
				ok("(USB) Bluetooth: \t%s\n", usbid.Classify(desc))
			} else if desc.Protocol == 0x02 { // Wireless
				ok("(USB) Wireless: \t%s", usbid.Classify(desc))
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
}
