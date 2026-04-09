package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type OS struct {
	Name string
	Args []string
}

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("This program is not supported on Windows.")
		return
	} else {
		os := OS{Name: runtime.GOOS}
		//fmt.Printf("The program is running on %s \n", runtime.GOOS)
		os.buildBaselineMap()
		allpackages := os.buildDiffMap()

		tabley := table.New("Name", "Action", "Date").WithPadding(3)

		headerColor := color.New(color.FgCyan, color.Bold, color.Underline).SprintfFunc()
		tabley.WithHeaderFormatter(headerColor)

		addedColor := color.New(color.FgGreen).SprintfFunc()
		removedColor := color.New(color.FgRed).SprintfFunc()

		for _, pkg := range allpackages.Packages {
			name := pkg.Name

			switch {
			case strings.Contains(pkg.Action, "install"):
				name = addedColor("+%s", name)
			case strings.Contains(pkg.Action, "remove"):
				name = removedColor("-%s", name)
			default:
				name = pkg.Name
			}
			tabley.AddRow(name, pkg.Action, pkg.Date)

		}

		tabley.Print()
	}
}
