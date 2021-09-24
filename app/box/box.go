package box

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/patrickhener/htb/config"
	"github.com/patrickhener/htb/helper"
)

// Box will hold the box object
type Box struct {
	name            string
	htbdir          string
	lootdir         string
	reportdir       string
	writeupTexFile  string
	baseTexFile     string
	preambleTexFile string
	config          *config.Config
	pdf             string
}

// New is a convenience method to create new box object
func New(boxname string, cfg *config.Config) *Box {
	box := &Box{
		name:   boxname,
		htbdir: cfg.HTBDir,
		config: cfg,
	}

	box.lootdir = path.Join(box.htbdir, "loot", box.name)
	box.reportdir = path.Join(box.htbdir, "writeup", box.name)
	box.baseTexFile = path.Join(box.reportdir, fmt.Sprintf("%s-writeup.tex", box.name))
	box.preambleTexFile = path.Join(box.reportdir, "files", "preamble.tex")
	box.writeupTexFile = path.Join(box.reportdir, "files", "writeup.tex")
	box.pdf = path.Join(box.reportdir, fmt.Sprintf("%s-writeup.pdf", box.name))

	return box
}

// Create will create folder structure and copy template
func (box *Box) Create() error {
	// Check and handle lootdir creation
	if err := helper.CreateLootDir(box.lootdir); err != nil {
		return err
	}

	// Check and handle reportdir creation
	if err := helper.CreateReportDir(box.reportdir, box.name, box.baseTexFile, box.preambleTexFile, box.config); err != nil {
		return err
	}

	// Fetch Batch image
	if err := helper.UpdateBadge(box.config); err != nil {
		return err
	}

	yes, err := helper.GrabYes("[*] Open the box to edit right away? [Y/n]")
	if err != nil {
		return err
	}
	if yes {
		if err := box.Edit(); err != nil {
			return err
		}
	}

	return nil
}

// Edit will open the markdown file in an editor using xdg-open
func (box *Box) Edit() error {
	// Check if box was already created
	if _, err := os.Stat(box.reportdir); os.IsNotExist(err) {
		return fmt.Errorf("%s", "box is not there")
	}

	// Open the md file
	cmd := exec.Command("code", box.reportdir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("cannot open the file to edit: %+v", err)
	}

	return nil
}

// Open will open the report
func (box *Box) Open() error {
	// Check if box was already created
	if _, err := os.Stat(box.reportdir); os.IsNotExist(err) {
		return fmt.Errorf("%s", "box is not there")
	}
	// Check if pdf was compiled already
	if _, err := os.Stat(box.pdf); os.IsNotExist(err) {
		return fmt.Errorf("%s", "the writeup was not compiled yet")
	}
	// Open the pdf
	cmd := exec.Command("xdg-open", box.pdf)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

// Clear will remove box folders and files
func (box *Box) Clear() error {
	fmt.Println("[*] Do you really want to delete the box content? Then type all upper: YES")
	choice := bufio.NewReader(os.Stdin)
	answer, err := choice.ReadString('\n')

	if err != nil {
		return err
	}

	if answer == "YES\n" {
		if err := os.RemoveAll(box.lootdir); err != nil {
			return err
		}
		if err := os.RemoveAll(box.reportdir); err != nil {
			return err
		}
		fmt.Println("[+] Box was cleared successfully.")

	} else {
		fmt.Println("[*] Nothing was deleted. Aborting.")
	}

	return nil
}
