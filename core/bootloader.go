package core

import (
	"log"

	"github.com/diskfs/go-diskfs"
)

var CLOVER_REPO = "CloverHackyColor/CloverBootloader"
var OPENCORE_REPO = "acidanthera/OpenCorePkg"

func mkBootloader(device string) {
	disk, err := diskfs.Open(device)
	if err != nil {
		log.Fatal("Failed to open device %s", device)
	} else {
	}
}
