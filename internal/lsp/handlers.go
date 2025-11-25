package lsp

import (
	"encoding/json"
	"strings"
)

func handleInitialize(req RequestMessage) {
	result := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: 1,
			CompletionProvider: &CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{"."},
			},
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
	updateSymbols(p.TextDocument.URI, p.TextDocument.Text)
	runDiagnostics(p.TextDocument.URI)
}

func handleDidChange(params json.RawMessage) {
	var p DidChangeTextDocumentParams
	json.Unmarshal(params, &p)
	if len(p.ContentChanges) > 0 {
		docs[p.TextDocument.URI] = p.ContentChanges[0].Text
		updateSymbols(p.TextDocument.URI, p.ContentChanges[0].Text)
		runDiagnostics(p.TextDocument.URI)
	}
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

var keywords = []string{
	"as", "assert", "Bool", "Byte", "else", "Float", "fn", "for", "if", "in", "Int", "List", "match", "mut", "Nil", "panic", "pub", "String", "type", "use", "val", "where",
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

	for _, kw := range keywords {
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
