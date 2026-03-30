package main

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
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
		fmt.Printf("The program is running on %s \n", runtime.GOOS)
		os.buildBaselineMap()
		installsAndDeletes := os.buildDiffMap()
		fmt.Println("The following are the packages that were installed or removed:")
		for _, pkg := range installsAndDeletes {
			if pkg[0] == '+' {
				color.Green(pkg) // Print installed packages in green
			} else if pkg[0] == '-' {
				color.Red(pkg) // Print removed packages in red
			} else {
				fmt.Println(pkg) // Print any other packages without color
			}
		}
	}
}
