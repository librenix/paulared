package core

import (
	"github.com/diskfs/go-diskfs"
	"github.com/fatih/color"
)

var CLOVER_REPO = "CloverHackyColor/CloverBootloader"
var OPENCORE_REPO = "acidanthera/OpenCorePkg"

func mkBootloader(device string) {
	disk, err := diskfs.Open(device)
	if err != nil {
		color.Red("Failed to open device %s", device)
	} else {
    }
}
