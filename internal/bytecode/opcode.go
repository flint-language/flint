package bytecode

import "strconv"

type OpCode byte

const (
	// Stack / Constants
	OP_CONST OpCode = iota
	OP_LOAD
	OP_STORE
	OP_POP

	// Integer Math
	OP_ADD // +
	OP_SUB // -
	OP_MUL // *
	OP_DIV // /
	OP_MOD // %

	// Integer Logic
	OP_LT // <
	OP_GT // >
	OP_LE // <=
	OP_GE // >=

	// Float Math
	OP_FADD // +.
	OP_FSUB // -.
	OP_FMUL // *.
	OP_FDIV // /.

	// Float Logic
	OP_FLT // <.
	OP_FGT // >.
	OP_FLE // <=.
	OP_FGE // >=.

	// String Logic
	OP_LTGT // <>

	// Common Logic
	OP_EQ      // ==
	OP_NEQ     // !=
	OP_OR      // |
	OP_OR_OR   // ||
	OP_AND     // &
	OP_AND_AND // &&
	OP_NOT     // !

	// Control Flow
	OP_JUMP
	OP_JUMP_IF_FALSE
	OP_CALL
	OP_RETURN

	// termination
	OP_PRINT
	OP_HALT // Program must end with HALT
)

func (op OpCode) String() string {
	switch op {

	// Stack / Constants
	case OP_CONST:
		return "CONST"
	case OP_LOAD:
		return "LOAD"
	case OP_STORE:
		return "STORE"
	case OP_POP:
		return "POP"

	// Integer Math
	case OP_ADD:
		return "ADD"
	case OP_SUB:
		return "SUB"
	case OP_MUL:
		return "MUL"
	case OP_DIV:
		return "DIV"
	case OP_MOD:
		return "MOD"

	// Integer Logic
	case OP_LT:
		return "LT"
	case OP_GT:
		return "GT"
	case OP_LE:
		return "LE"
	case OP_GE:
		return "GE"

	// Float Math
	case OP_FADD:
		return "FADD"
	case OP_FSUB:
		return "FSUB"
	case OP_FMUL:
		return "FMUL"
	case OP_FDIV:
		return "FDIV"

	// Float Logic
	case OP_FLT:
		return "FLT"
	case OP_FGT:
		return "FGT"
	case OP_FLE:
		return "FLE"
	case OP_FGE:
		return "FGE"

	// String Logic
	case OP_LTGT:
		return "LTGT"

	// Common Logic
	case OP_EQ:
		return "EQ"
	case OP_NEQ:
		return "NEQ"
	case OP_OR:
		return "OR"
	case OP_OR_OR:
		return "OR_OR"
	case OP_AND:
		return "AND"
	case OP_AND_AND:
		return "AND_AND"
	case OP_NOT:
		return "NOT"

	// Control Flow
	case OP_JUMP:
		return "JUMP"
	case OP_JUMP_IF_FALSE:
		return "JUMP_IF_FALSE"
	case OP_CALL:
		return "CALL"
	case OP_RETURN:
		return "RETURN"

	// Termination
	case OP_PRINT:
		return "PRINT"
	case OP_HALT:
		return "HALT"

	default:
		return "OP<" + strconv.Itoa(int(op)) + ">"
	}
}
