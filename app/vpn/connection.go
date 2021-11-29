package vpn

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Handle() error {
	// Fetch active vpns
	var activeVPNs []string
	var activeErr bytes.Buffer
	var activeOut bytes.Buffer
	activeCmd := exec.Command("nmcli", "connection", "show", "--active")
	activeCmd.Stdout = &activeOut
	activeCmd.Stderr = &activeErr
	err := activeCmd.Run()
	if err != nil {
		return err
	}

	if activeErr.String() != "" {
		return fmt.Errorf("%s", &activeErr)
	}

	scanner := bufio.NewScanner(&activeOut)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "NAME") {
			continue
		}
		if strings.Contains(line, "vpn") {
			activeVPNs = append(activeVPNs, strings.Fields(line)[0])
		}
	}

	// Fetch inactive vpns
	var inactiveVPNs []string
	var inactiveErr bytes.Buffer
	var inactiveOut bytes.Buffer
	// activeConnections := []string{}
	inactiveCmd := exec.Command("nmcli", "connection", "show")
	inactiveCmd.Stdout = &inactiveOut
	inactiveCmd.Stderr = &inactiveErr
	err = inactiveCmd.Run()
	if err != nil {
		return err
	}

	if activeErr.String() != "" {
		return fmt.Errorf("%s", &activeErr)
	}

	scanner = bufio.NewScanner(&inactiveOut)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "NAME") {
			continue
		}
		if strings.Contains(line, "--") && strings.Contains(line, "vpn") {
			inactiveVPNs = append(inactiveVPNs, strings.Fields(line)[0])
		}
	}

	// Print results
	fmt.Println("Active VPNs:")
	for i, a := range activeVPNs {
		fmt.Printf("[%d] %s\n", i, a)
	}

	fmt.Println("")

	fmt.Println("Inactive VPNs:")
	for i, u := range inactiveVPNs {
		fmt.Printf("[%d] %s\n", i+len(activeVPNs), u)
	}

	// Read choice
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please choose a number. Active VPNs will be disconnected, inactive VPNs will be connected: ")

	choice, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	choice = strings.ReplaceAll(choice, "\n", "")

	// Convert choice to int
	iChoice, err := strconv.Atoi(choice)
	if err != nil {
		fmt.Println("You need to choose a numeric value")
		return err
	}

	fmt.Println("")

	// Determine which was chosen and either connect a thing or disconnect a thing
	switch {
	case iChoice <= len(activeVPNs)-1:
		fmt.Printf("Disconnecting '%s'\n", activeVPNs[iChoice])
		dis := exec.Command("nmcli", "connection", "down", activeVPNs[iChoice])
		out, err := dis.CombinedOutput()
		if err != nil {
			return err
		}

		fmt.Println(string(out))

	default:
		fmt.Printf("Connecting '%s'\n", inactiveVPNs[iChoice-len(activeVPNs)])

		con := exec.Command("nmcli", "connection", "up", inactiveVPNs[iChoice-len(activeVPNs)])
		out, err := con.CombinedOutput()
		if err != nil {
			return err
		}

		fmt.Println(string(out))
	}

	return nil
}
