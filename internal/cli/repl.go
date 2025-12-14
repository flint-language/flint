package cli

import (
	"bufio"
	"flint/internal/interpreter"
	"flint/internal/version"
	"fmt"
	"os"
	"strings"
)

func startRepl() {
	fmt.Printf("Flint version: %s\n", version.Version)
	fmt.Println("Enter :h for help.")
	env := interpreter.NewEnv(nil)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == ":h" {
			fmt.Println("")
			fmt.Println(":h\n   Prints a list of avaliable directives.")
			fmt.Println(":q\n   Exit the toplevel.")
			continue
		}
		if line == ":q" {
			return
		}
		interpreter.RunReplLine(line, env)
	}
}
