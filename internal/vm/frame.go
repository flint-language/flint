package vm

import "flint/internal/bytecode"

type callFrame struct {
	ip         int
	chunk      *bytecode.Chunk
	funcIndex  int
	stackStart int
}