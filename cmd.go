package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Args
var bootloader string
var check bool
var drive string
var img string
var kexts []string

var rootCmd = &cobra.Command{
	Use:   "paulared",
	Short: "Cross platform hackintosh installer",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Version = version

		// Check hardware and exit
		if check {
			Compatibility()
			return nil
		}

		// Check args is valid
		switch bootloader {
		case "clover":
			//InstallClover()
		case "opencore":
			//InstallOpenCore()
		case "none":
			// do nothings
		default:
			return fmt.Errorf("bootloader %s is invalid", bootloader)
		}

		if drive != "" || img != "" {
			//WriteToDisk(drive, img)
		}

		//InstallKexts(kexts)

		return nil
	},
}

func Execute() {
	// Flags
	rootCmd.Flags().StringVarP(&bootloader, "bootloader", "b", "none", "install bootloader (clover or opencore)")
	rootCmd.Flags().BoolVarP(&check, "check", "c", false, "check hardware compatibility")
	rootCmd.Flags().StringVarP(&drive, "drive", "d", "", "macOS install target device")
	rootCmd.Flags().StringVarP(&img, "image", "i", "", "download or select disk image")
	rootCmd.Flags().StringSliceVarP(&kexts, "kext", "k", []string{"fakesmc"}, "install kexts")
	rootCmd.Flags().BoolP("version", "v", false, "output version")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
