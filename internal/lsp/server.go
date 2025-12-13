package lsp

import (
	"bufio"
	"encoding/json"
	"os"
)

func StartLsp() {
	reader := bufio.NewReader(os.Stdin)
	for {
		body, err := readMessage(reader)
		if err != nil {
			return
		}
		var req RequestMessage
		if err := json.Unmarshal(body, &req); err != nil {
			continue
		}
		switch req.Method {
		case "initialize":
			handleInitialize(req)
		case "textDocument/didOpen":
			handleDidOpen(req.Params)
		case "textDocument/didChange":
			handleDidChange(req.Params)
		case "textDocument/completion":
			handleCompletion(req)
		case "textDocument/hover":
			handleHover(req)
		case "textDocument/codeLens":
			handleCodeLens(req)
		case "shutdown":
			handleShutdown(req)
		case "exit":
			os.Exit(0)
		}
	}
}
