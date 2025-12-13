package lsp

import (
	"encoding/json"
	"flint/internal/lexer"
	"strings"
	"time"
)

func handleInitialize(req RequestMessage) {
	result := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: 1,
			CompletionProvider: &CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{"."},
			},
			HoverProvider: true,
		},
	}
	send(ResponseMessage{
		Jsonrpc: "2.0",
		ID:      req.ID,
		Result:  result,
	})
}

func handleDidOpen(params json.RawMessage) {
	var p DidOpenTextDocumentParams
	json.Unmarshal(params, &p)
	docs[p.TextDocument.URI] = p.TextDocument.Text
	parseMu.Lock()
	if timer, ok := parseTimers[p.TextDocument.URI]; ok {
		timer.Stop()
	}
	parseTimers[p.TextDocument.URI] = time.AfterFunc(150*time.Millisecond, func() {
		parseAndUpdateSymbols(p.TextDocument.URI)
	})
	parseMu.Unlock()
}

func handleDidChange(params json.RawMessage) {
	var p DidChangeTextDocumentParams
	json.Unmarshal(params, &p)
	if len(p.ContentChanges) == 0 {
		return
	}
	docs[p.TextDocument.URI] = p.ContentChanges[0].Text
	parseMu.Lock()
	if timer, ok := parseTimers[p.TextDocument.URI]; ok {
		timer.Stop()
	}
	parseTimers[p.TextDocument.URI] = time.AfterFunc(150*time.Millisecond, func() {
		parseAndUpdateSymbols(p.TextDocument.URI)
	})
	parseMu.Unlock()
}

func handleShutdown(req RequestMessage) {
	send(ResponseMessage{
		Jsonrpc: "2.0",
		ID:      req.ID,
		Result:  nil,
	})
}

func runDiagnostics(uri string) {
	text := docs[uri]
	diagnostics := []Diagnostic{}
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.Contains(line, "error") {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: 1,
				Message:  "Detected the word 'error'",
				Range: Range{
					Start: Position{Line: i, Character: 0},
					End:   Position{Line: i, Character: len(line)},
				},
			})
		}
	}
	send(map[string]any{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnostics,
		},
	})
}

func handleCompletion(req RequestMessage) {
	var params CompletionParams
	json.Unmarshal(req.Params, &params)
	text := docs[params.TextDocument.URI]
	lines := strings.Split(text, "\n")
	line := ""
	if params.Position.Line < len(lines) {
		line = lines[params.Position.Line]
	}
	prefix := ""
	fields := strings.Fields(line[:params.Position.Character])
	if len(fields) > 0 {
		prefix = fields[len(fields)-1]
	}
	suggestions := []CompletionItem{}
	for kw := range lexer.KeywordMap {
		if strings.HasPrefix(kw, prefix) {
			suggestions = append(suggestions, CompletionItem{
				Label: kw,
				Kind:  14,
			})
		}
	}
	for _, sym := range symbols[params.TextDocument.URI] {
		if strings.HasPrefix(sym.Name, prefix) {
			kind := 6
			if sym.Kind == FunctionSymbol {
				kind = 3
			}
			suggestions = append(suggestions, CompletionItem{
				Label: sym.Name,
				Kind:  kind,
			})
		}
	}
	send(ResponseMessage{
		Jsonrpc: "2.0",
		ID:      req.ID,
		Result:  suggestions,
	})
}

func handleHover(req RequestMessage) {
	var params HoverParams
	json.Unmarshal(req.Params, &params)
	text, ok := docs[params.TextDocument.URI]
	if !ok {
		return
	}
	lines := strings.Split(text, "\n")
	if params.Position.Line >= len(lines) {
		return
	}
	line := lines[params.Position.Line]
	pos := min(params.Position.Character, len(line))
	start := pos
	for start > 0 && isIdentifierChar(line[start-1]) {
		start--
	}
	end := pos
	for end < len(line) && isIdentifierChar(line[end]) {
		end++
	}
	if start >= end {
		return
	}
	word := line[start:end]
	syms, ok := symbols[params.TextDocument.URI]
	if !ok {
		return
	}
	for _, sym := range syms {
		if sym.Name == word {
			text := "```flint\n" + sym.Type + "\n```"
			send(ResponseMessage{
				Jsonrpc: "2.0",
				ID:      req.ID,
				Result: HoverResult{
					Contents: text,
				},
			})
			return
		}
	}
}
