package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGen) defineExit() {
	statusParam := ir.NewParam("statusCode", cg.platformIntType())
	flint_exit := cg.mod.NewFunc("flint_exit", types.Void, statusParam)
	flint_exit.FuncAttrs = append(flint_exit.FuncAttrs, enum.FuncAttrNoReturn)
	flint_exit_dispatch := cg.mod.NewFunc("flint_exit_dispatch", types.I32, statusParam)

	dispatch_entry := flint_exit_dispatch.NewBlock("entry")
	dispatch_entry.NewRet(dispatch_entry.NewTrunc(statusParam, types.I32))

	exit_entry := flint_exit.NewBlock("entry")
	exit_entry.NewCall(flint_exit_dispatch, statusParam).Tail = enum.TailMustTail
	exit_entry.NewUnreachable()

	cg.funcs["exit"] = flint_exit
}

// func (cg *CodeGen) defineAssert() {
// 	valueParam := ir.NewParam("v", types.I1)
// 	flint_assert := cg.mod.NewFunc("flint_assert", types.Void, valueParam)

// 	assert_entry := flint_assert.NewBlock("entry")
// 	// assert_true := flint_assert.NewBlock("assert.true")
// 	assert_false := flint_assert.NewBlock("assert.false")

// 	cmp := assert_entry.NewICmp(enum.IPredNE, valueParam, constant.NewInt(types.I1, 1))
// 	assert_entry.NewCondBr(cmp, cg.mod._, nil)
// 	assert_entry.NewRet(nil)

// 	assert_entry.NewEx
// 	assert_false.NewRet(nil)

// 	cg.funcs["assert"] = flint_assert
// }

func (cg *CodeGen) defineStrlen() {
	intType := cg.platformIntType()

	strParam := ir.NewParam("str", types.I8Ptr)
	flint_strlen := cg.mod.NewFunc("flint_strlen", intType, strParam)

	entry := flint_strlen.NewBlock("entry")
	loop := flint_strlen.NewBlock("loop")
	inc := flint_strlen.NewBlock("loop.inc")
	exit := flint_strlen.NewBlock("loop.exit")

	zeroInt := constant.NewInt(intType, 0)
	oneInt := constant.NewInt(intType, 1)

	zero8 := constant.NewInt(types.I8, 0)

	// entry_start
	strIdx := entry.NewAlloca(intType)
	intPtr := entry.NewAlloca(intType)

	entry.NewStore(zeroInt, strIdx)
	entry.NewStore(entry.NewPtrToInt(strParam, intType), intPtr)

	entry.NewBr(loop)
	// entry_end

	// loop_start
	intPtrAdd := loop.NewAdd(
		loop.NewLoad(intType, intPtr),
		loop.NewLoad(intType, strIdx),
	)
	strPtr := loop.NewIntToPtr(intPtrAdd, types.I8Ptr)

	chr := loop.NewLoad(types.I8, strPtr)
	cmp := loop.NewICmp(enum.IPredEQ, chr, zero8)
	loop.NewCondBr(cmp, exit, inc)
	// loop_end

	// exit_start
	idxLoad := exit.NewLoad(intType, strIdx)
	exit.NewRet(idxLoad)
	// exit_end

	// inc_start
	strPtrAdd := inc.NewAdd(inc.NewLoad(intType, strIdx), oneInt)
	inc.NewStore(strPtrAdd, strIdx)
	inc.NewBr(loop)
	// inc_end

	cg.funcs["strlen"] = flint_strlen
}
