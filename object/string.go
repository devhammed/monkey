package object

import "hash/fnv"

type String struct {
	Value string
}

func (s *String) Type() Type {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

// Clone creates a new copy
func (s *String) Clone() Object {
	return &String{Value: s.Value}
}

func (s *String) HashKey() (HashKey, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(s.Value))
	if err != nil {
		return HashKey{}, err
	}

	return HashKey{Type: s.Type(), Value: h.Sum64()}, nil
}
