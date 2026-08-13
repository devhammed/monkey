package object

import "fmt"

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Clone() Object {
	return &Integer{Value: i.Value}
}

func (i *Integer) HashKey() (HashKey, error) {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}, nil
}
