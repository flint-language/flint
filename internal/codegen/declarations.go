package codegen

import (
	"flint/internal/parser"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitVarDecl(b *ir.Block, e *parser.VarDeclExpr) value.Value {
	expr := cg.emitExpr(b, e.Value, false)
	alloc := b.NewAlloca(expr.Type())
	cg.locals[e.Name.Lexeme] = alloc
	b.NewStore(expr, alloc)
	return alloc
}

func (cg *CodeGen) resolveType(t parser.Expr) types.Type {
	if t == nil {
		return types.Void
	}
	switch ty := t.(type) {
	case *parser.TypeExpr:
		switch ty.Name {
		case "Int":
			return cg.platformIntType()
		case "Float":
			return cg.platformFloatType()
		case "Bool":
			return types.I1
		case "Byte":
			return types.I8
		case "String":
			return types.I8Ptr
		case "Nil":
			return types.Void
		default:
			return nil
		}
	case *parser.TupleTypeExpr:
		subTypes := make([]types.Type, len(ty.Types))
		for i, e := range ty.Types {
			subTypes[i] = cg.resolveType(e)
		}
		return types.NewStruct(subTypes...)
	default:
		panic(fmt.Sprintf("unexpected type node: %T", t))
	}
}
