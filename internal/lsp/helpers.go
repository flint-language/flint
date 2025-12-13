package lsp

import (
	"bufio"
	"encoding/json"
	"flint/internal/lexer"
	"flint/internal/parser"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

var (
	parseTimers = map[string]*time.Timer{}
	docs        = map[string]string{}
	symbols     = map[string][]Symbol{}
	parseMu     sync.Mutex
)

func send(msg any) {
	data, _ := json.Marshal(msg)
	fmt.Printf("Content-Length: %d\r\n\r\n%s", len(data), data)
}

func readMessage(r *bufio.Reader) ([]byte, error) {
	headers := map[string]string{}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	length := 0
	fmt.Sscanf(headers["Content-Length"], "%d", &length)
	body := make([]byte, length)
	_, err := io.ReadFull(r, body)
	return body, err
}

func isIdentifierChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '_'
}

func parseAndUpdateSymbols(uri string) {
	text := docs[uri]
	tokens, err := lexer.Tokenize(text, uri)
	if err != nil {
		return
	}
	prog, _ := parser.ParseProgram(tokens)
	updateSymbols(uri, prog)
	runDiagnostics(uri)
	tokens = nil
	prog = nil
}

func formatType(e parser.Expr) string {
	if e == nil {
		return "_"
	}
	switch t := e.(type) {
	case *parser.TypeExpr:
		if t.Generic != nil {
			return t.Name + "<" + formatType(t.Generic) + ">"
		}
		return t.Name
	default:
		return "_"
	}
}
