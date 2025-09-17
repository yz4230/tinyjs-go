package subjs

type Ident string
type StringValue string
type NumberValue int

type AddExpr struct {
	Left, Right any
}

func (e *AddExpr) Eval() any {
	if l, ok := e.Left.(int); ok {
		if r, ok := e.Right.(int); ok {
			return l + r
		}
	}
	if l, ok := e.Left.(string); ok {
		if r, ok := e.Right.(string); ok {
			return l + r
		}
	}

	panic("unsupported types for addition")
}

type CallExpr struct {
	Name Ident
	Args []any
}

type MethodCallExpr struct {
	Receiver []any
	Method   Ident
	Args     []any
}
