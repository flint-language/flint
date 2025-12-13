package cli

import "fmt"

func interpretFile(filename string) {
	fmt.Println("Running " + filename)
	_, _ = loadAndParse(filename)
	fmt.Println("Program executed (no backend yet)")
}
