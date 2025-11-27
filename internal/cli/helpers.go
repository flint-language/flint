package cli

import (
	"flint/internal/lexer"
	"flint/internal/parser"
	"flint/internal/typechecker"
	"fmt"
	"os"
)

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func loadAndParse(filename string) (*parser.Program, *typechecker.TypeChecker) {
	tc := typechecker.New()
	src, err := os.ReadFile(filename)
	if err != nil {
		fatal(fmt.Sprintf("error reading %s: %v", filename, err))
	}

	tokens, err := lexer.Tokenize(string(src), filename)
	if err != nil {
		fatal(err.Error())
	}

	prog, errs := parser.ParseProgram(tokens)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, e)
		}
		os.Exit(1)
	}

	for _, ex := range prog.Exprs {
		if _, err := tc.CheckExpr(ex); err != nil {
			fatal("Type error: " + err.Error())
		}
	}
	return prog, tc
}
