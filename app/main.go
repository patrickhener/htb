package main

import (
	"fmt"
	"os"
	"path"

	"github.com/patrickhener/htb/box"
	"github.com/patrickhener/htb/config"
	"github.com/patrickhener/htb/helper"
	"github.com/patrickhener/htb/vpn"
)

func main() {
	// Init config
	cfg := config.New()
	if err := cfg.Init(); err != nil {
		fmt.Printf("[-] Error when initializing the config: %+v\n", err)
		os.Exit(-1)
	}

	var mode string = os.Args[1]
	var reportdir string = path.Join(os.Getenv("HTBDIR"), "writeup")

	switch mode {
	case "badge":
		// if mode is badge just update the badge.png
		if err := helper.UpdateBadge(cfg); err != nil {
			fmt.Printf("[-] Error updating badge: %+v\n", err)
		}
		os.Exit(0)
	case "list":
		// if mode is list just list and exit
		if err := helper.List(reportdir); err != nil {
			fmt.Printf("[-] Error when listing existing boxes: %+v\n", err)
		}
		os.Exit(0)
	case "vpn":
		if err := vpn.Handle(); err != nil {
			fmt.Printf("[-] Error when handling vpn connection: %+v\n", err)
		}
		os.Exit(0)
	}

	// This hits if mode is none of the above
	// In this case switch again over mode
	var boxname string = os.Args[2]
	box := box.New(boxname, cfg)

	switch mode {
	case "create":
		if err := box.Create(); err != nil {
			fmt.Printf("[-] Error when creating box assets: %+v\n", err)
		}
	case "edit":
		if err := box.Edit(); err != nil {
			fmt.Printf("[-] Error when editing box report: %+v\n", err)
		}
	case "open":
		if err := box.Open(); err != nil {
			fmt.Printf("[-] Error when opening box report: %+v\n", err)
		}
	case "clear":
		if err := box.Clear(); err != nil {
			fmt.Printf("[-] Error when clearing box report: %+v\n", err)
		}
	default:
		fmt.Println("Valid modes are: create, edit, open, list, clear, badge and vpn")
	}
}
