package object

// Environment is where the variables are stored
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment returns a new environment ref
func NewEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

// NewEnclosedEnvironment returns an extension of the passed environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get returns the solicited object
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets a new value into the environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
