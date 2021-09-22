package helper

import (
	"bufio"
	"fmt"
	"os"
)

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
