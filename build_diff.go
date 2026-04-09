package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Package struct {
	Name   string
	Action string
	Date   string
}

type PackageList struct {
	Packages []Package
}

var packageDiffs = make(map[int]string)

func (o OS) getDiff() (string, error) {
	out, err := exec.Command("dnf", "history", "list").CombinedOutput()
	if err != nil {
		// Include the command output because dnf often explains the failure there
		return "", fmt.Errorf("dnf history list failed: %w\noutput:\n%s", err, string(out))
	}
	return string(out), nil
}

func (o OS) getPkgDeets(transactionID int) (version string, installDate string) {
	pkgDetails, err := exec.Command("dnf", "history", "info", strconv.Itoa(transactionID)).CombinedOutput()
	if err != nil {
		// Include the command output because dnf often explains the failure there
		return "Hi", "Nooo"
	}
	scanny := bufio.NewScanner(strings.NewReader(string(pkgDetails)))
	for scanny.Scan() {
		line := strings.TrimSpace(scanny.Text())
		if strings.Contains(line, "Packages Altered:") {
			scanny.Scan()
			line = strings.TrimSpace(scanny.Text())
			version = strings.Fields(line)[1]
		} else if strings.Contains(line, "End time") {
			line = strings.TrimSpace(scanny.Text())
			installDate = strings.Split(line, " : ")[1]
		}
	}
	return version, installDate

}

func (o OS) buildDiffMap() PackageList {
	allDiffs, err := o.getDiff()
	if err != nil {
		return PackageList{}
	}

	var packagio Package
	var packagelisty PackageList

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
		if len(fields) < 2 {
			continue // or return an error if you want strict parsing
		}

		transID := fields[0]
		transID = strings.Trim(transID, " ")
		transIDInt, err := strconv.Atoi(transID)
		if err != nil {
			fmt.Printf("Error converting transaction ID to integer: %v\n", err)
			continue
		}
		action := fields[1]
		packageDiffs[transIDInt] = action
		version, installDate := o.getPkgDeets(transIDInt)

		packagio.Name = version
		packagio.Action = action
		packagio.Date = installDate

		packagelisty.Packages = append(packagelisty.Packages, packagio)

	}
	if len(packageDiffs) == 0 {
		fmt.Println("No package differences found.")
		return PackageList{}
	}

	return packagelisty
}
