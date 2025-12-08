package typechecker

import "maps"

type Env struct {
	vars    map[string]VarInfo
	parent  *Env
	modules map[string]*Env
}

type VarInfo struct {
	Ty      *Type
	Mutable bool
}

func NewEnv(parent *Env) *Env {
	modules := make(map[string]*Env)
	if parent != nil && parent.modules != nil {
		maps.Copy(modules, parent.modules)
	}
	return &Env{
		vars:    make(map[string]VarInfo),
		parent:  parent,
		modules: modules,
	}
}

func (e *Env) GetVar(name string) (VarInfo, bool) {
	if v, ok := e.vars[name]; ok {
		return v, true
	}
	if e.parent != nil {
		return e.parent.GetVar(name)
	}
	return VarInfo{}, false
}

func (e *Env) SetVar(name string, ty *Type, mutable bool) {
	e.vars[name] = VarInfo{Ty: ty, Mutable: mutable}
}

func (e *Env) Get(name string) (*Type, bool) {
	if v, ok := e.GetVar(name); ok {
		return v.Ty, true
	}
	return nil, false
}

func (e *Env) Set(name string, ty *Type) {
	e.SetVar(name, ty, true)
}
