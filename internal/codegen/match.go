package codegen

import (
	"flint/internal/parser"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
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
