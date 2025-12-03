package lsp

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func resetState() {
	docs = map[string]string{}
	symbols = map[string][]Symbol{}
}

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestInitialize(t *testing.T) {
	resetState()

	id := json.RawMessage(`1`)
	req := RequestMessage{
		Jsonrpc: "2.0",
		ID:      &id,
		Method:  "initialize",
	}

	out := captureStdout(func() {
		handleInitialize(req)
	})

	if !strings.Contains(out, `"textDocumentSync":1`) {
		t.Fatalf("expected textDocumentSync in response, got:\n%s", out)
	}

	if !strings.Contains(out, `"completionProvider"`) {
		t.Fatalf("expected completionProvider in response, got:\n%s", out)
	}
}

func TestDidOpenStoresDocument(t *testing.T) {
	resetState()

	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:  "file:///test.flint",
			Text: "val x = 10",
		},
	}

	data, _ := json.Marshal(params)

	captureStdout(func() {
		handleDidOpen(data)
	})

	if docs["file:///test.flint"] != "val x = 10" {
		t.Fatalf("document not stored")
	}
}

func TestDiagnosticsDetectsErrorWord(t *testing.T) {
	resetState()

	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:  "file:///err.flint",
			Text: "this is an error here",
		},
	}

	data, _ := json.Marshal(params)

	out := captureStdout(func() {
		handleDidOpen(data)
	})

	if !strings.Contains(out, `"textDocument/publishDiagnostics"`) {
		t.Fatalf("did not publish diagnostics")
	}

	if !strings.Contains(out, "Detected the word 'error'") {
		t.Fatalf("missing diagnostic message, got:\n%s", out)
	}
}

func TestSymbolDetection(t *testing.T) {
	resetState()

	uri := "file:///symbols.flint"
	text := `
fn add(x: Int) Int { x }
val hello = 10
mut world = 20
`

	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{URI: uri, Text: text},
	}

	data, _ := json.Marshal(params)
	handleDidOpen(data)

	syms := symbols[uri]

	if len(syms) != 3 {
		t.Fatalf("expected 3 symbols, got %d", len(syms))
	}

	if syms[0].Name != "add" || syms[0].Kind != FunctionSymbol {
		t.Fatal("function symbol not detected")
	}

	if syms[1].Name != "hello" {
		t.Fatal("val symbol not detected")
	}
}

func TestDidChangeUpdatesDocument(t *testing.T) {
	resetState()

	uri := "file:///change.flint"
	docs[uri] = "old"

	params := DidChangeTextDocumentParams{
		TextDocument: VersionedTextDocumentIdentifier{URI: uri},
		ContentChanges: []TextDocumentContentChange{
			{Text: "new content"},
		},
	}

	data, _ := json.Marshal(params)
	handleDidChange(data)

	if docs[uri] != "new content" {
		t.Fatalf("document did not update")
	}
}
