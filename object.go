package monkey_interpreter

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

type ObjectType string

const (
	INT_OBJ_TYPE          = "INTEGER"
	STRING_OBJ_TYPE       = "STRING"
	BOOL_OBJ_TYPE         = "BOOLEAN"
	NULL_OBJ_TYPE         = "NULL"
	RETURN_VALUE_OBJ_TYPE = "RETURN_VALUE"
	ERROR_OBJ_TYPE        = "ERROR"
	FUNCTION_OBJ_TYPE     = "FUNCTION"
	BULTIN_OBJ_TYPE       = "BUILTIN"
	ARRAY_OBJ_TYPE        = "ARRAY"
	HASH_OBJ_TYPE         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INT_OBJ_TYPE }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ_TYPE }
func (s *String) Inspect() string  { return s.Value }

type BooleanObject struct {
	Value bool
}

func (b *BooleanObject) Type() ObjectType { return BOOL_OBJ_TYPE }
func (b *BooleanObject) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ_TYPE }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ_TYPE }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ_TYPE }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*Identifier
	Body       *BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ_TYPE }
func (f *Function) Inspect() string {
	var (
		out    bytes.Buffer
		params []string
	)
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BULTIN_OBJ_TYPE }
func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ_TYPE }

func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *BooleanObject) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ_TYPE }

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

type Hashable interface {
	HashKey() HashKey
}
