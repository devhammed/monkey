package object

import "unicode"

// Environment is an object that holds a mapping of names to bound objets
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment constructs a new Environment object to hold bindings
// of identifiers to their names
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), outer: nil}
}

// NewEnclosedEnvironment returns a new Environment with the parent set to the current
// environment (enclosing environment)
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

// ExportedHash returns a new Hash with the names and values of every publically
// exported binding in the environment. That is every binding that starts with a
// capital letter. This is used by the module import system to wrap up the
// evaluated module into an object.
func (e *Environment) ExportedHash() *Hash {
	pairs := make(map[HashKey]HashPair)

	for k, v := range e.store {
		if unicode.IsUpper(rune(k[0])) {
			s := &String{Value: k}

			pairs[s.HashKey()] = HashPair{Key: s, Value: v}
		}
	}

	return &Hash{Pairs: pairs}
}

// Get returns the object bound by name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set stores the object with the given name
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val

	return val
}
