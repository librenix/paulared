package core

import (
	"log"

	"github.com/diskfs/go-diskfs"
)

var OPENCORE_REPO = "acidanthera/OpenCorePkg"
var CLOVER_REPO = "CloverHackyColor/CloverBootloader"

func mkBootloader(device string) {
	disk, err := diskfs.Open(device)
	if err != nil {
		log.Fatal("Failed to open device %s", device)
	} else {
	}
}
