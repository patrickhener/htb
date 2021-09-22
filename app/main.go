package main

import (
	"fmt"
	"os"
	"path"

	"git.hener.eu/patrick/htb/box"
	"git.hener.eu/patrick/htb/helper"
)

func init() {
	// Check for arguments
	if len(os.Args) < 3 {
		if len(os.Args) == 1 || os.Args[1] != "list" {
			fmt.Printf("Usage: %+v <mode> <box-name>\n", os.Args[0])
			fmt.Println("Valid modes are: create, edit, convert, open, list or clear")
			os.Exit(1)
		}
	}

	// Check for env variable
	htbdir := os.Getenv("HTBDIR")
	if htbdir == "" {
		fmt.Println("For the app to work you will have to add an environment variable 'HTBDIR' pointing to the htb root folder!")
		os.Exit(2)
	}

	// check if software is installed
	if yes := helper.CheckTool("xelatex"); !yes {
		fmt.Print("You need to install xelatex (texlive) for the framework to work.\nCheck the README.md\n")
		os.Exit(0)
	}
	if yes := helper.CheckTool("pandoc"); !yes {
		fmt.Print("You need to install pandoc for the framework to work.\nCheck the README.md\n")
		os.Exit(0)
	}
	if yes := helper.CheckTool("panrun"); !yes {
		fmt.Print("You need to install panrun for the framework to work.\nCheck the README.md\n")
		os.Exit(0)
	}
	if yes := helper.CheckTool("phantomjs"); !yes {
		fmt.Println("Phantomjs is missing. Badge creation will not work")
	}

	if yes := helper.CheckTool("obsidian"); !yes {
		fmt.Println("Obsidian is missing. You cannot use mode <edit>")
	}

	// check if eisvogel template is there
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Cannot read users home dir: %+v", err)
		os.Exit(0)
	}
	checkEisvogelPaths := []string{path.Join(homeDir, ".local", "share", "pandoc", "templates", "eisvogel.latex"), path.Join(homeDir, ".pandoc", "templates", "eisvogel.latex")}

	found := false
	for _, p := range checkEisvogelPaths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			continue
		}
		found = true
	}
	if !found {
		fmt.Print("You need to install eisvogel latex template for the framework to work.\nCheck the README.md\n")
		os.Exit(0)
	}
}

func main() {
	var mode string = os.Args[1]
	if mode != "list" {
		var (
			boxname string = os.Args[2]
			htbdir  string = os.Getenv("HTBDIR")
		)
		box := box.New(boxname, htbdir)
		switch mode {
		case "create":
			if err := box.Create(); err != nil {
				fmt.Printf("Error when creating box assets: %+v\n", err)
			}
		case "edit":
			if err := box.Edit(); err != nil {
				fmt.Printf("Error when editing box report: %+v\n", err)
			}
		case "convert":
			if err := box.Convert(); err != nil {
				fmt.Printf("Error when converting report to pdf: %+v\n", err)
			}
		case "open":
			if err := box.Open(); err != nil {
				fmt.Printf("Error when opening box report: %+v\n", err)
			}
		case "clear":
			if err := box.Clear(); err != nil {
				fmt.Printf("Error when clearing box report: %+v\n", err)
			}
		default:
			fmt.Println("Valid modes are: create, edit, convert, open, list or clear")
			os.Exit(1)
		}
	} else {
		reportdir := path.Join(os.Getenv("HTBDIR"), "report")
		if err := helper.List(reportdir); err != nil {
			fmt.Printf("Error when listing existing boxes: %+v\n", err)
		}
	}

}
