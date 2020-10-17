package controlflow

import (
	"fmt"
	"testing"

	"github.com/dannypsnl/extend"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func compileStmt(f *ir.Func, bb *ir.Block, stmt Stmt) {
	switch s := stmt.(type) {
	case *SIf:
		thenB := extend.Block(f.NewBlock(""))
		compileStmt(f, thenB.Block, s.Then)
		elseB := f.NewBlock("")
		compileStmt(f, elseB, s.Else)
		bb.NewCondBr(compileExpr(bb, s.Cond), thenB.Block, elseB)
		if thenB.HasTerminator() {
			leaveB := f.NewBlock("")
			thenB.NewBr(leaveB)
		}
	case *SRet:
		bb.NewRet(compileExpr(bb, s.Val))
	}
}

func TestParameterAttr(t *testing.T) {
	f := ir.NewFunc("foo", types.Void)
	bb := f.NewBlock("")

	compileStmt(f, bb, &SIf{
		Cond: &EBool{V: true},
		Then: nil,
		Else: &SRet{Val: &EVoid{}},
	})

	// whatever what we did in compileStmt, we use convention that a block leave in the end is empty.
	f.Blocks[len(f.Blocks)-1].NewRet(nil)

	fmt.Println(f.LLString())
}

// generated LLVM IR:
//
// ```
// define void @foo() {
// ; <label>:0
// br i1 true, label %1, label %2
//
// ; <label>:1
// br label %3
//
// ; <label>:2
// ret void
//
// ; <label>:3
// ret void
// }
// ```
