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
			Description: "Compile and execute a Flint source file in a single step.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				runFile(os.Args[2])
			},
		},
		{
			Name:        "compile",
			Description: "Compile a Flint program to the selected backend (LLVM IR, object file, or executable).",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				compileFile(os.Args[2])
			},
		},
		{
			Name:        "check",
			Description: "Perform full type-checking on a Flint source file without generating code.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				checkFile(os.Args[2])
			},
		},
		{
			Name:        "interpret",
			Description: "Execute a Flint program using the bytecode virtual machine.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				interpretFile(os.Args[2])
			},
		},
		{
			Name:        "repl",
			Description: "Launch an interactive Read-Eval-Print Loop (REPL) for Flint.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				startRepl()
			},
		},
		{
			Name:        "lsp",
			Description: "Launch the Flint Language Server for editor integration and IDE features.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				startLsp()
			},
		},
		{
			Name:        "version",
			Description: "Print the currently installed Flint compiler version.",
			Run: func(fs *flag.FlagSet) {
				fmt.Println(version.FullVersion())
			},
		},
		{
			Name:        "help",
			Description: "Show help information for Flint or a specific subcommand.",
			Run: func(fs *flag.FlagSet) {
				fs.Parse(os.Args[2:])
				printHelp()
			},
		},
	}
}
