package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var packageDiffs = make(map[int]string)
var installsAndDeletes []string

func (o OS) getDiff() (string, error) {
	out, err := exec.Command("dnf", "history", "list").CombinedOutput()
	if err != nil {
		// Include the command output because dnf often explains the failure there
		return "", fmt.Errorf("dnf history list failed: %w\noutput:\n%s", err, string(out))
	}
	return string(out), nil
}

func (o OS) buildDiffMap() []string {
	allDiffs, err := o.getDiff()
	if err != nil {
		return nil
	}

	sc := bufio.NewScanner(strings.NewReader(allDiffs))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		//fmt.Printf("Processing line: %s\n", line) // Debugging output
		if strings.Contains(line, "ID") ||
			strings.Contains(line, "Command line") ||
			strings.Contains(line, "Date and time") ||
			strings.Contains(line, "Action(s)") ||
			strings.Contains(line, "-------") {
			//fmt.Printf("Skipping line: %s\n", line) // Debugging output
			continue
		}

		fields := strings.Split(line, "|")
		fmt.Printf("Fields: %v\n", fields) // Debugging output
		if len(fields) < 2 {
			continue // or return an error if you want strict parsing
		}

		transID := fields[0]
		transID = strings.Trim(transID, " ")
		fmt.Printf("Transaction ID: %v\n", transID) // Debugging output
		transIDInt, err := strconv.Atoi(transID)
		if err != nil {
			fmt.Printf("Error converting transaction ID to integer: %v\n", err)
			continue
		}
		action := fields[1]
		packageDiffs[transIDInt] = action

		if strings.Contains(action, "install") || strings.Contains(action, "remove") {
			actionandPkg := strings.Fields(action)
			installsAndDeletes = append(installsAndDeletes, actionandPkg[1])
		}

	}

	if err := sc.Err(); err != nil {
		//return fmt.Errorf("scan failed: %w", err)
		return nil
	}

	//fmt.Println(baselinePackages)
	if len(packageDiffs) == 0 {
		fmt.Println("No package differences found.")
		return nil
	}

	fmt.Println("The following is the diff map:")

	for k, v := range packageDiffs {
		fmt.Printf("%d: %s\n", k, v)
	}
	return installsAndDeletes
}
