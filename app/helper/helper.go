package helper

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/patrickhener/htb/config"
)

// List will walk the report dir and print every name of the subdirectories
func List(reportdir string) error {
	content, err := ioutil.ReadDir(reportdir)
	if err != nil {
		return err
	}

	for _, c := range content {
		if c.IsDir() {
			fmt.Println(c.Name())
		}
	}
	return nil
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("Reportdir %+v already exists", dst)
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// GrabYes will receive a question to ask and return true if answered with y/Y or Enter
func GrabYes(text string) (bool, error) {
	fmt.Println(text)
	choice := bufio.NewReader(os.Stdin)
	char, _, err := choice.ReadRune()

	if err != nil {
		return false, err
	}

	switch char {
	case 'Y', 'y', '\n':
		return true, nil
	}
	return false, nil
}

// CreateLootDir will create the loot dir for the box
// It will also create an 'nmap' folder in the loot dir
func CreateLootDir(lootpath string) error {
	if _, err := os.Stat(lootpath); os.IsNotExist(err) {
		yes, err := GrabYes(fmt.Sprintf("[*] Lootdir %+v does not exist. Create it? [Y/n]", lootpath))
		if err != nil {
			return err
		}
		if yes {
			// Create loot folder
			if err := os.MkdirAll(lootpath, 0755); err != nil {
				return err
			}
			// Create nmap directory
			if err := os.Mkdir(path.Join(lootpath, "nmap"), 0755); err != nil {
				return err
			}

			fmt.Printf("[+] Directory %+v created\n", lootpath)
		} else {
			// Do not create and exit
			fmt.Println("[-] Aborting operation. Exiting.")
			os.Exit(-1)
		}
	} else {
		fmt.Println("[-] Box already exists. Aborting.")
		os.Exit(-1)
	}

	return nil
}

// CreateReportDir will create the report directory for a box
// And copy over the template files from "writeup" style
// It will also fill in the details for box and author
func CreateReportDir(reportdir, boxname, basetexfile, preambletexfile string, cfg *config.Config) error {
	if _, err := os.Stat(reportdir); os.IsNotExist(err) {
		yes, err := GrabYes(fmt.Sprintf("[*] Reportdir %+v does not exist. Create it? [Y/n]", reportdir))
		if err != nil {
			return err
		}

		if yes {
			// Create directory and copy over template assets using helper function
			if err := CopyDir(path.Join(cfg.WriteupLatexPath, "template"), reportdir); err != nil {
				return err
			}
			fmt.Println("[+] Template files copied over")

			// Rename template report file
			oldName := path.Join(reportdir, "report.tex")
			if err := os.Rename(oldName, basetexfile); err != nil {
				return err
			}
			// Replace placeholders in preamble
			read, err := ioutil.ReadFile(preambletexfile)
			if err != nil {
				return err
			}

			newContent := strings.Replace(string(read), "++machinename++", strings.Title(boxname), -1)
			newContent = strings.Replace(string(newContent), "++authorname++", cfg.HTBAuthor, -1)
			if err := ioutil.WriteFile(preambletexfile, []byte(newContent), 0); err != nil {
				return err
			}
		} else {
			// Do not create and exit
			return fmt.Errorf("%s", "[*] Aborting operation")
		}
	} else {
		return fmt.Errorf("%s", "[-] Box already exists")
	}

	return nil
}

/*
func fetchBatch(box *Box) error {



		// Now use phantomjs and badge.js to convert html to badge.png
		cmd := exec.Command("phantomjs", box.badgejs)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			return err
		}
}
*/

// UpdateBadge will fetch the html code from hackthebox.eu and then
// generate a png out of it using phantomjs
// It will copy it into 'writeup' Latex style image directory as the style
// is using it in the front matter
func UpdateBadge(cfg *config.Config) error {
	if cfg.HTBProfileID != "" {
		// Make http request to fetch batch raw response
		resp, err := http.Get(fmt.Sprintf("https://www.hackthebox.eu/badge/%s", cfg.HTBProfileID))
		if err != nil {
			return err
		}

		// Read in body
		rawBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Trim left and right to extraxt pure base64 part
		b64Body := strings.TrimPrefix(string(rawBody), "document.write(window.atob(\"")
		b64Body = strings.TrimSuffix(b64Body, "\"))")

		// Decode base64 to html
		html, err := base64.StdEncoding.DecodeString(b64Body)
		if err != nil {
			return err
		}

		// Replace some things
		htmlStr := strings.Replace(string(html), "<div ", "<div class=\"wrapper\" ", 1)
		htmlStr = strings.ReplaceAll(htmlStr, "https://www.hackthebox.eu/images/screenshot.png", fmt.Sprintf("data:image/png;base64,%s", HTBCROSSHAIR))
		htmlStr = strings.ReplaceAll(htmlStr, "_thumb.png", ".png")
		htmlStr = strings.ReplaceAll(htmlStr, "https://www.hackthebox.eu/images/star.png", fmt.Sprintf("data:image/png;base64,%s", HTBSTAR))
		htmlStr = strings.ReplaceAll(htmlStr, "url(https://www.hackthebox.eu/images/icon20.png);", fmt.Sprintf("url('data:image/webp;base64,%s'; background-size: 20px;", HTBLOGO))

		// Make temp directory to operate in
		tmpDir, err := ioutil.TempDir("", "htb")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)

		// Write badge.html
		if err := ioutil.WriteFile(path.Join(tmpDir, "badge.html"), []byte(htmlStr), 0755); err != nil {
			return err
		}

		// Writing badge.js
		badgeJs := strings.ReplaceAll(BADGEJS, "%%tmpdirhere%%", tmpDir)

		if err := ioutil.WriteFile(path.Join(tmpDir, "badge.js"), []byte(badgeJs), 0755); err != nil {
			return err
		}

		// Now use phantomjs and badge.js to convert html to badge.png
		cmd := exec.Command("phantomjs", path.Join(tmpDir, "badge.js"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return err
		}

		// Finally copy badge.png to /texmf/tex/latex/writeup/images
		if err := CopyFile(path.Join(tmpDir, "badge.png"), path.Join(cfg.WriteupLatexPath, "images", "badge.png")); err != nil {
			return err
		}

		return nil
	}
	return fmt.Errorf("%s", "HTBPROFILEID not set - no banner update")
}
