package object

import (
	"fmt"
	"math"
)

type Float struct {
	Value float64
}

func (f *Float) Type() Type {
	return FLOAT_OBJ
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%g", f.Value)
}

func (f *Float) Clone() Object {
	return &Float{Value: f.Value}
}

func (f *Float) HashKey() (HashKey, error) {
	value := f.Value
	// Keep equal positive and negative zero values in the same hash bucket.
	if value == 0 {
		value = 0
	}
	return HashKey{Type: f.Type(), Value: math.Float64bits(value)}, nil
}
