package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

var baselinePackages = make(map[string]string)

func (o OS) getBaseline() (string, error) {
	out, err := exec.Command("dnf", "list", "installed").CombinedOutput()
	if err != nil {
		// Include the command output because dnf often explains the failure there
		return "", fmt.Errorf("dnf failed: %w\noutput:\n%s", err, string(out))
	}
	return string(out), nil
}

func (o OS) buildBaselineMap() map[string]string {
	allpkgs, err := o.getBaseline()
	if err != nil {
		return nil
	}

	sc := bufio.NewScanner(strings.NewReader(allpkgs))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || line == "Installed Packages" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue // or return an error if you want strict parsing
		}

		nameArch := fields[0]
		version := fields[1]
		baselinePackages[nameArch] = version

	}

	if err := sc.Err(); err != nil {
		//return fmt.Errorf("scan failed: %w", err)
		return nil
	}

	//fmt.Println("The following is the baseline map:")
	//fmt.Println(baselinePackages)
	//for k, v := range baselinePackages {
	//	fmt.Printf("%s: %s\n", k, v)
	//}
	return baselinePackages
}
