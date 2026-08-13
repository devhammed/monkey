package object

// Type represents a type of object
type Type string

// Object Types
const (
	INTEGER_OBJ      Type = "INTEGER"
	FLOAT_OBJ        Type = "FLOAT"
	BOOLEAN_OBJ      Type = "BOOLEAN"
	NULL_OBJ         Type = "NULL"
	RETURN_VALUE_OBJ Type = "RETURN_VALUE"
	ERROR_OBJ        Type = "ERROR"
	FUNCTION_OBJ     Type = "FUNCTION"
	STRING_OBJ       Type = "STRING"
	BUILTIN_OBJ      Type = "BUILTIN"
	RESOURCE_OBJ     Type = "RESOURCE"
	ARRAY_OBJ        Type = "ARRAY"
	HASH_OBJ         Type = "HASH"
)

// Immutable is the interface for all immutable objects which must implement
// the Clone() method used by binding names to values.
type Immutable interface {
	Clone() Object
}

type Object interface {
	Type() Type
	Inspect() string
}
