package helper

import (
	"fmt"
	"io/ioutil"
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
