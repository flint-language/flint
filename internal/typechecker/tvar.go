package typechecker

import "sync"

var (
	nextVarID int
	varLock   sync.Mutex
)

func newVarID() int {
	varLock.Lock()
	defer varLock.Unlock()
	nextVarID++
	return nextVarID
}

func NewTypeVar(fam FamilyKind) *Type {
	return &Type{
		TKind:  TyVar,
		VarID:  newVarID(),
		Family: fam,
	}
}
