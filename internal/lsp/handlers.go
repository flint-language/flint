package lsp

import (
	"encoding/json"
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
			CodeLensProvider: &CodeLensOptions{
				ResolveProvider: false,
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
	word := ""
	start := min(params.Position.Character, len(line))
	i := start - 1
	for i >= 0 && (isIdentifierChar(line[i])) {
		i--
	}
	left := i + 1
	j := start
	for j < len(line) && isIdentifierChar(line[j]) {
		j++
	}
	right := j
	if left < right {
		word = line[left:right]
	}
	var hoverText string
	for _, sym := range symbols[params.TextDocument.URI] {
		if sym.Name == word {
			hoverText = "```flint\n" + sym.Type + "\n```"
			break
		}
	}
	send(ResponseMessage{
		Jsonrpc: "2.0",
		ID:      req.ID,
		Result: HoverResult{
			Contents: hoverText,
		},
	})
}

func handleCodeLens(req RequestMessage) {
	var params CodeLensParams
	json.Unmarshal(req.Params, &params)
	syms := symbols[params.TextDocument.URI]
	lenses := []CodeLens{}
	for _, sym := range syms {
		if sym.Kind != FunctionSymbol {
			continue
		}
		lenses = append(lenses, CodeLens{
			Range: Range{
				Start: Position{
					Line:      sym.Line - 1,
					Character: 0,
				},
				End: Position{
					Line:      sym.Line - 1,
					Character: 0,
				},
			},
			Command: &Command{
				Title: sym.CurriedSig,
			},
		})
	}
	send(ResponseMessage{
		Jsonrpc: "2.0",
		ID:      req.ID,
		Result:  lenses,
	})
}
