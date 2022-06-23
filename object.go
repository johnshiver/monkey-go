package monkey_interpreter

import "fmt"

type ObjectType string

const (
	INT_OBJ_TYPE          = "INTEGER"
	BOOL_OBJ_TYPE         = "BOOLEAN"
	NULL_OBJ_TYPE         = "NULL"
	RETURN_VALUE_OBJ_TYPE = "RETURN_VALUE"
	ERROR_OBJ_TYPE        = "ERROR"
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
