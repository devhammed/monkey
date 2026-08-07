package object

import "unicode"

// Binding is an object that holds a bound object
type Binding struct {
	Value Object
	BindingOptions
}

// BindingOptions is an object that holds options for a binding
type BindingOptions struct {
	SuperGlobal bool
}

// Environment is an object that holds a mapping of names to bound objets
type Environment struct {
	store map[string]Binding
	outer *Environment
}

// NewEnvironment constructs a new Environment object to hold bindings
// of identifiers to their names
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Binding), outer: nil}
}

// NewEnclosedEnvironment returns a new Environment with the parent set to the current
// environment (enclosing environment)
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

// NewModuleEnvironment creates a new environment containing only superglobal bindings from the provided parent environment.
func NewModuleEnvironment(parent *Environment) *Environment {
	env := NewEnvironment()

	for current := parent; current != nil; current = current.outer {
		for name, binding := range current.store {
			if binding.SuperGlobal {
				env.Set(name, binding.Value, binding.BindingOptions)
			}
		}
	}

	return env
}

// ExportedHash returns a new Hash with the names and values of every publicly
// exported binding in the environment. That is every binding that starts with a
// capital letter. This is used by the module import system to wrap up the
// evaluated module into an object.
func (e *Environment) ExportedHash() *Hash {
	pairs := make(map[HashKey]HashPair)

	for k, v := range e.store {
		if unicode.IsUpper(rune(k[0])) && !v.SuperGlobal {
			s := &String{Value: k}

			hashKey, err := s.HashKey()
			if err != nil {
				continue
			}

			pairs[hashKey] = HashPair{Key: s, Value: v.Value}
		}
	}

	return &Hash{Pairs: pairs}
}

// Get returns the object bound by name
func (e *Environment) Get(name string) (Binding, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set stores the object with the given name
func (e *Environment) Set(name string, val Object, options BindingOptions) Binding {
	binding := Binding{Value: val, BindingOptions: options}

	e.store[name] = binding

	return binding
}
