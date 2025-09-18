package subjs

import (
	"reflect"
	"testing"
)

func TestParse_Factor(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"0", NumberValue(0)},
		{"42", NumberValue(42)},
		{"123456", NumberValue(123456)},
		{`"hello"`, StringValue("hello")},
		{`'world'`, StringValue("world")},
		{`"12345"`, StringValue("12345")},
		{`(0)`, NumberValue(0)},
		{`(42)`, NumberValue(42)},
		{`("hello")`, StringValue("hello")},
		{`('world')`, StringValue("world")},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]rune(tt.input))
			parser := NewParser()
			parser.Parse(lexer)

			if lexer.Err != nil {
				t.Error(lexer.Err)
			}

			if !reflect.DeepEqual(lexer.Result, tt.want) {
				t.Errorf("want %v, got %v", tt.want, lexer.Result)
			}
		})
	}
}

func TestParse_Expr(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"0+0", &AddExpr{Left: NumberValue(0), Right: NumberValue(0)}},
		{"1+2", &AddExpr{Left: NumberValue(1), Right: NumberValue(2)}},
		{"42+100", &AddExpr{Left: NumberValue(42), Right: NumberValue(100)}},
		{"1+2+3", &AddExpr{Left: &AddExpr{Left: NumberValue(1), Right: NumberValue(2)}, Right: NumberValue(3)}},
		{"10+20+30+40", &AddExpr{
			Left: &AddExpr{
				Left: &AddExpr{
					Left:  NumberValue(10),
					Right: NumberValue(20)},
				Right: NumberValue(30)},
			Right: NumberValue(40)}},
		{"(1+2)+3", &AddExpr{Left: &AddExpr{Left: NumberValue(1), Right: NumberValue(2)}, Right: NumberValue(3)}},
		{"1+(2+3)", &AddExpr{Left: NumberValue(1), Right: &AddExpr{Left: NumberValue(2), Right: NumberValue(3)}}},
		{"(1+2)+(3+4)", &AddExpr{
			Left: &AddExpr{
				Left:  NumberValue(1),
				Right: NumberValue(2)},
			Right: &AddExpr{
				Left:  NumberValue(3),
				Right: NumberValue(4)}}},
		{`"hello"+"world"`, &AddExpr{Left: StringValue("hello"), Right: StringValue("world")}},
		{`"foo"+"bar"+"baz"`, &AddExpr{Left: &AddExpr{Left: StringValue("foo"), Right: StringValue("bar")}, Right: StringValue("baz")}},
		{`("foo"+"bar")+"baz"`, &AddExpr{Left: &AddExpr{Left: StringValue("foo"), Right: StringValue("bar")}, Right: StringValue("baz")}},
		{`"foo"+("bar"+"baz")`, &AddExpr{Left: StringValue("foo"), Right: &AddExpr{Left: StringValue("bar"), Right: StringValue("baz")}}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]rune(tt.input))
			parser := NewParser()
			parser.Parse(lexer)

			if lexer.Err != nil {
				t.Error(lexer.Err)
			}

			if !reflect.DeepEqual(lexer.Result, tt.want) {
				t.Errorf("want %v, got %v", tt.want, lexer.Result)
			}
		})
	}
}

func TestParse_Call(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"foo()", &CallExpr{Name: "foo", Args: []any{}}},
		{"bar(1)", &CallExpr{Name: "bar", Args: []any{NumberValue(1)}}},
		{"baz(1,2,3)", &CallExpr{Name: "baz", Args: []any{NumberValue(1), NumberValue(2), NumberValue(3)}}},
		{`qux("hello","world")`, &CallExpr{Name: "qux", Args: []any{StringValue("hello"), StringValue("world")}}},
		{"sum(1+2,3+4)", &CallExpr{
			Name: "sum",
			Args: []any{
				&AddExpr{
					Left:  NumberValue(1),
					Right: NumberValue(2)},
				&AddExpr{
					Left:  NumberValue(3),
					Right: NumberValue(4)}}}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]rune(tt.input))
			parser := NewParser()
			parser.Parse(lexer)

			if lexer.Err != nil {
				t.Error(lexer.Err)
			}

			if !reflect.DeepEqual(lexer.Result, tt.want) {
				t.Errorf("want %v, got %v", tt.want, lexer.Result)
			}
		})
	}
}

func TestParse_Idents(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"obj", []any{Ident("obj")}},
		{"obj.method", []any{Ident("obj"), Ident("method")}},
		{"a.b.c.d", []any{Ident("a"), Ident("b"), Ident("c"), Ident("d")}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]rune(tt.input))
			parser := NewParser()
			parser.Parse(lexer)

			if lexer.Err != nil {
				t.Error(lexer.Err)
			}

			if !reflect.DeepEqual(lexer.Result, tt.want) {
				t.Errorf("want %v, got %v", tt.want, lexer.Result)
			}
		})
	}
}

func TestParse_MethodCall(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"obj.method()", &MethodCallExpr{
			Receiver: []any{Ident("obj")},
			Method:   Ident("method"),
			Args:     []any{}}},
		{"obj.method(1)", &MethodCallExpr{
			Receiver: []any{Ident("obj")},
			Method:   Ident("method"),
			Args:     []any{NumberValue(1)}}},
		{"obj.method(1,2,3)", &MethodCallExpr{
			Receiver: []any{Ident("obj")},
			Method:   Ident("method"),
			Args:     []any{NumberValue(1), NumberValue(2), NumberValue(3)}}},
		{`obj.method("hello","world")`, &MethodCallExpr{
			Receiver: []any{Ident("obj")},
			Method:   Ident("method"),
			Args:     []any{StringValue("hello"), StringValue("world")}}},
		{"obj.method(1+2,3+4)", &MethodCallExpr{
			Receiver: []any{Ident("obj")},
			Method:   Ident("method"),
			Args: []any{&AddExpr{
				Left:  NumberValue(1),
				Right: NumberValue(2)}, &AddExpr{
				Left:  NumberValue(3),
				Right: NumberValue(4)}}}},
		{"a.b.c.d.method()", &MethodCallExpr{
			Receiver: []any{Ident("a"), Ident("b"), Ident("c"), Ident("d")},
			Method:   Ident("method"),
			Args:     []any{}}},
		{"a.b.c.d.method(42)", &MethodCallExpr{
			Receiver: []any{Ident("a"), Ident("b"), Ident("c"), Ident("d")},
			Method:   Ident("method"),
			Args:     []any{NumberValue(42)}}},
		{`"abc".substr(1,2)`, &MethodCallExpr{
			Receiver: []any{StringValue("abc")},
			Method:   Ident("substr"),
			Args:     []any{NumberValue(1), NumberValue(2)}}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]rune(tt.input))
			parser := NewParser()
			parser.Parse(lexer)

			if lexer.Err != nil {
				t.Error(lexer.Err)
			}

			if !reflect.DeepEqual(lexer.Result, tt.want) {
				t.Errorf("want %v, got %v", tt.want, lexer.Result)
			}
		})
	}
}

func TestParse_Complex(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{`("hello"+"world").substr(1,2)`, &MethodCallExpr{
			Receiver: []any{&AddExpr{
				Left:  StringValue("hello"),
				Right: StringValue("world")}},
			Method: Ident("substr"),
			Args:   []any{NumberValue(1), NumberValue(2)}}},
		{`("b" + "c" + "d").substr(1, 2)`, &MethodCallExpr{
			Receiver: []any{&AddExpr{
				Left: &AddExpr{
					Left:  StringValue("b"),
					Right: StringValue("c")},
				Right: StringValue("d")}},
			Method: Ident("substr"),
			Args:   []any{NumberValue(1), NumberValue(2)}}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]rune(tt.input))
			parser := NewParser()
			parser.Parse(lexer)

			if lexer.Err != nil {
				t.Error(lexer.Err)
			}

			if !reflect.DeepEqual(lexer.Result, tt.want) {
				t.Errorf("want %v, got %v", tt.want, lexer.Result)
			}
		})
	}
}

func TestEval_Spaces(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"  1  + 2   ", &AddExpr{Left: NumberValue(1), Right: NumberValue(2)}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]rune(tt.input))
			parser := NewParser()
			parser.Parse(lexer)

			if lexer.Err != nil {
				t.Error(lexer.Err)
			}

			if !reflect.DeepEqual(lexer.Result, tt.want) {
				t.Errorf("want %v, got %v", tt.want, lexer.Result)
			}
		})
	}
}
