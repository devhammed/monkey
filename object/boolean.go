package object

import "fmt"

type Boolean struct {
	Value bool
}

// Type returns the type of the object
func (b *Boolean) Type() Type {
	return BOOLEAN_OBJ
}

// Clone creates a new copy
func (b *Boolean) Clone() Object {
	return &Boolean{Value: b.Value}
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}
