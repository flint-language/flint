package codegen

import (
	"flint/internal/lexer"
	"flint/internal/parser"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGen) emitMatch(b *ir.Block, m *parser.MatchExpr, isTail bool) value.Value {
	parent := b.Parent
	matchId := cg.globalMatchCount
	cg.globalMatchCount++
	armNames := make([]string, 0)
	armMap := make(map[string]*parser.MatchArm)
	armBlockMap := make(map[string]*ir.Block)
	armCheckMap := make(map[string]*ir.Block)
	armBodyMap := make(map[string]value.Value)
	checkList := make([]*ir.Block, 0)
	nextList := make([]*ir.Block, 0)
	checkList = append(checkList, b)
	wildCardName := fmt.Sprintf("match.%d.wild", matchId)
	for caseId, arm := range m.Arms {
		var name string
		if arm.IsWildCardArm() {
			name = wildCardName
		} else {
			name = fmt.Sprintf("match.%d.arm.%d", matchId, caseId)
		}
		armNames = append(armNames, name)
		armBlock := parent.NewBlock(name)
		if caseId != 0 && !arm.IsWildCardArm() {
			checkName := fmt.Sprintf("match.%d.check.%d", matchId, caseId)
			checkList = append(checkList, parent.NewBlock(checkName))
			nextList = append(nextList, checkList[len(checkList)-1])
			armCheckMap[checkName] = checkList[len(checkList)-1]
		} else if caseId != 0 {
			checkList = append(checkList, nil)
			nextList = append(nextList, armBlock)
		}
		armMap[name] = arm
		armBlockMap[name] = armBlock
		armBodyMap[name] = cg.emitMatchBody(armBlock, arm.Body, isTail)
	}
	scrutinee := cg.emitExpr(b, m.Value, false)
	var incomings []*ir.Incoming
	var phiType types.Type
	mergeBlock := parent.NewBlock(fmt.Sprintf("match.%d", matchId))
	nextList = append(nextList, mergeBlock)
	current := b
	armId := 0
	for _, name := range armNames {
		armBlock := armBlockMap[name]
		armBody := armBodyMap[name]
		if armBody != nil {
			if phiType == nil {
				phiType = armBody.Type()
			} else if !armBody.Type().Equal(phiType) {
				panic(fmt.Sprintf("match arm type mismatch: %v vs %v (arm %d)", armBody.Type(), phiType, armId))
			}
		}
		if !armMap[name].IsWildCardArm() {
			current = checkList[armId]
			if current != nil && current.Term == nil {
				arm := armMap[name]
				cond := cg.emitMatchCond(current, scrutinee, arm.Pattern, arm.Guard)
				current.NewCondBr(cond, armBlock, nextList[armId])
			}
		}
		if armBlock.Term == nil {
			armBlock.NewBr(mergeBlock)
		}
		if phiType != nil {
			if armBody != nil {
				incomings = append(incomings, &ir.Incoming{X: armBody, Pred: armBlock})
			} else {
				incomings = append(incomings, &ir.Incoming{X: constant.NewUndef(phiType), Pred: armBlock})
			}
		}
		armId++
	}
	if current != nil && current.Term == nil {
		current.NewBr(mergeBlock)
	}
	if phiType == nil || phiType.Equal(types.Void) {
		if mergeBlock.Term == nil {
			mergeBlock.NewRet(nil)
		}
		return mergeBlock
	}
	phi := mergeBlock.NewPhi(incomings...)
	if mergeBlock.Term == nil {
		mergeBlock.NewRet(phi)
	}
	return phi
}

func (cg *CodeGen) emitMatchBody(b *ir.Block, body parser.Expr, isTail bool) value.Value {
	switch x := body.(type) {
	case *parser.BlockExpr:
		return cg.emitBlock(b, x, isTail)
	case *parser.MatchExpr:
		return cg.emitMatch(b, x, isTail)
	default:
		return cg.emitExpr(b, body, isTail)
	}
}

func (cg *CodeGen) emitMatchCond(b *ir.Block, scr value.Value, pat parser.Expr, guard parser.Expr) value.Value {
	var baseCond value.Value
	switch p := pat.(type) {
	case *parser.IntLiteral:
		if _, ok := scr.Type().(*types.IntType); !ok {
			panic("match scrutinee is not an integer for integer pattern")
		}
		lit := constant.NewInt(types.I64, p.Value)
		baseCond = b.NewICmp(enum.IPredEQ, scr, lit)
	case *parser.BoolLiteral:
		if _, ok := scr.Type().(*types.IntType); !ok {
			panic("match scrutinee is not an integer for bool pattern")
		}
		var intVal int64
		if p.Value {
			intVal = 1
		} else {
			intVal = 0
		}
		lit := constant.NewInt(types.I1, intVal)
		baseCond = b.NewICmp(enum.IPredEQ, scr, lit)
	case *parser.Identifier:
		if p.Pos.Kind == lexer.Underscore {
			baseCond = constant.True
		} else {
			alloc := b.NewAlloca(scr.Type())
			b.NewStore(scr, alloc)
			cg.locals[p.Name] = alloc
			baseCond = constant.True
		}
	default:
		panic("unsupported match pattern: " + pat.NodeType())
	}
	if guard != nil {
		guardVal := cg.emitExpr(b, guard, false)
		if _, ok := guardVal.Type().(*types.IntType); ok {
			if guardVal.Type() != types.I1 {
				guardVal = b.NewTrunc(guardVal, types.I1)
			}
		} else {
			panic("guard expression is not a boolean")
		}
		baseCond = b.NewAnd(baseCond, guardVal)
	}
	return baseCond
}
