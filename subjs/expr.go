package subjs

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
	Name string
	Args []any
}

type MethodCallExpr struct {
	Receiver []string
	Method   string
	Args     []any
}
