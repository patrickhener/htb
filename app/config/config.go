package config

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type Config struct {
	TexmfPath        string
	WriteupLatexPath string
	BoxLatexPath     string
	HTBDir           string
	HTBAuthor        string
	HTBProfileID     string
	RequiredTools    []string
	RequiredPaths    []string
	RequiredEnv      []string
}

func New() *Config {
	// check if eisvogel template is there
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("[-] Cannot read users home dir: %+v", err)
		os.Exit(-1)
	}

	cfg := &Config{
		TexmfPath:     path.Join(homeDir, "texmf", "tex", "latex"),
		HTBDir:        os.Getenv("HTBDIR"),
		HTBAuthor:     os.Getenv("HTBAUTHOR"),
		HTBProfileID:  os.Getenv("HTBPROFILEID"),
		RequiredTools: []string{"lualatex", "phantomjs"},
		RequiredEnv:   []string{"HTBDIR", "HTBAUTHOR", "HTBPROFILEID"},
	}

	cfg.WriteupLatexPath = path.Join(cfg.TexmfPath, "writeup")
	cfg.RequiredPaths = []string{cfg.TexmfPath, cfg.WriteupLatexPath}

	return cfg
}

func (c *Config) Init() error {
	var list []string
	var err error

	// Sanity checks
	if len(os.Args) < 3 {
		if len(os.Args) == 1 || os.Args[1] != "list" && os.Args[1] != "badge" {
			fmt.Printf("Usage: %+v <mode> <box-name>\n", os.Args[0])
			fmt.Println("Valid modes are: create, edit, open, list or clear")
			os.Exit(0)
		}
	}

	// Check paths
	list, err = CheckPaths(c.RequiredPaths)
	if err != nil {
		fmt.Printf("[-] There was an error with the path check: %+v\n", err)
		fmt.Printf("[-] The list of missing items is: %+v\n", list)
		os.Exit(-1)
	}

	// Check tools
	list, err = CheckTools(c.RequiredTools)
	if err != nil {
		fmt.Printf("[-] There was an error with the tools check: %+v\n", err)
		fmt.Printf("[-] The list of missing items is: %+v\n", list)
		os.Exit(-1)
	}

	// Check env
	list, err = CheckEnv(c.RequiredEnv)
	if err != nil {
		fmt.Printf("[-] There was an error with the environment check: %+v\n", err)
		fmt.Printf("[-] The list of missing items is: %+v\n", list)
		os.Exit(-1)
	}
	return nil
}

// CheckEnv will check for the required env variables
func CheckEnv(requiredEnv []string) ([]string, error) {
	missingList := []string{}
	missing := false

	for _, e := range requiredEnv {
		if os.Getenv(e) == "" {
			missing = true
			missingList = append(missingList, e)
		}
	}

	if missing {
		return missingList, fmt.Errorf("%s", "missing environment variable(s)")
	}

	return nil, nil
}

// CheckPaths will check for the required paths existence
func CheckPaths(requiredPaths []string) ([]string, error) {
	missingPaths := []string{}
	missing := false

	for _, p := range requiredPaths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			missing = true
			missingPaths = append(missingPaths, p)
		}
	}

	if missing {
		return missingPaths, fmt.Errorf("%s", "missing path(s)")
	}

	return nil, nil
}

// CheckTools will check for the required tools
func CheckTools(requiredTools []string) ([]string, error) {
	missingTools := []string{}
	missing := false

	for _, t := range requiredTools {
		if m := CheckTool(t); m {
			missing = true
			missingTools = append(missingTools, t)
		}
	}

	if missing {
		return missingTools, fmt.Errorf("%s", "missing tool(s)")
	}

	return nil, nil
}

// CheckTool will check if the required tool is there and then return true
func CheckTool(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return true
	}
	return false
}
