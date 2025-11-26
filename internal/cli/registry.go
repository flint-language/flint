package cli

import (
	"flag"
	"flint/internal/version"
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Description string
	Run         func(fs *flag.FlagSet)
}

var commands []Command

func init() {
	commands = []Command{
		{
			Name:        "run",
			Description: "Run a Flint program from a source file.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				runFile(os.Args[2])
			},
		},
		{
			Name:        "compile",
			Description: "Compile Flint code to a backend.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				compileFile(os.Args[2])
			},
		},
		{
			Name:        "check",
			Description: "Type-check Flint code without executing it.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				checkFile(os.Args[2])
			},
		},
		{
			Name:        "lsp",
			Description: "Start the Flint Language Server.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				startLsp()
			},
		},
		{
			Name:        "version",
			Description: "Print the Flint compiler version.",
			Run: func(fs *flag.FlagSet) {
				fmt.Println(version.FullVersion())
			},
		},
		{
			Name:        "help",
			Description: "Display general or command-specific help.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				printHelp()
			},
		},
	}
}
