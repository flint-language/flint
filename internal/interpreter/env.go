package interpreter

type Env struct {
	Parent *Env
	Vals   map[string]Value
}

func NewEnv(parent *Env) *Env {
	return &Env{
		Parent: parent,
		Vals:   make(map[string]Value),
	}
}

func (e *Env) Lookup(name string) (Value, bool) {
	for env := e; env != nil; env = env.Parent {
		if v, ok := env.Vals[name]; ok {
			return v, true
		}
	}
	return Value{}, false
}
