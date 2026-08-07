package object

import (
	"bytes"
	"fmt"
	"strings"
)

type HashKey struct {
	Type  Type
	Value uint64
}

type Hashable interface {
	HashKey() (HashKey, error)
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() Type {
	return HASH_OBJ
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	var pairs []string

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
