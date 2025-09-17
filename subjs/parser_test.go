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
		{"0", 0},
		{"42", 42},
		{"123456", 123456},
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`"12345"`, "12345"},
		{`(0)`, 0},
		{`(42)`, 42},
		{`("hello")`, "hello"},
		{`('world')`, "world"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input))
			parser := NewParser()
			parser.Parse(lexer)

			if lexer.Err != nil {
				t.Error(lexer.Err)
			}

			if lexer.Result != tt.want {
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
		{"0+0", &AddExpr{Left: 0, Right: 0}},
		{"1+2", &AddExpr{Left: 1, Right: 2}},
		{"42+100", &AddExpr{Left: 42, Right: 100}},
		{"1+2+3", &AddExpr{Left: &AddExpr{Left: 1, Right: 2}, Right: 3}},
		{"10+20+30+40", &AddExpr{Left: &AddExpr{Left: &AddExpr{Left: 10, Right: 20}, Right: 30}, Right: 40}},
		{"(1+2)+3", &AddExpr{Left: &AddExpr{Left: 1, Right: 2}, Right: 3}},
		{"1+(2+3)", &AddExpr{Left: 1, Right: &AddExpr{Left: 2, Right: 3}}},
		{"(1+2)+(3+4)", &AddExpr{Left: &AddExpr{Left: 1, Right: 2}, Right: &AddExpr{Left: 3, Right: 4}}},
		{`"hello"+"world"`, &AddExpr{Left: "hello", Right: "world"}},
		{`"foo"+"bar"+"baz"`, &AddExpr{Left: &AddExpr{Left: "foo", Right: "bar"}, Right: "baz"}},
		{`("foo"+"bar")+"baz"`, &AddExpr{Left: &AddExpr{Left: "foo", Right: "bar"}, Right: "baz"}},
		{`"foo"+("bar"+"baz")`, &AddExpr{Left: "foo", Right: &AddExpr{Left: "bar", Right: "baz"}}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input))
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
		{"bar(1)", &CallExpr{Name: "bar", Args: []any{1}}},
		{"baz(1,2,3)", &CallExpr{Name: "baz", Args: []any{1, 2, 3}}},
		{`qux("hello","world")`, &CallExpr{Name: "qux", Args: []any{"hello", "world"}}},
		{"sum(1+2,3+4)", &CallExpr{Name: "sum", Args: []any{&AddExpr{Left: 1, Right: 2}, &AddExpr{Left: 3, Right: 4}}}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input))
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
		{"obj", []string{"obj"}},
		{"obj.method", []string{"obj", "method"}},
		{"a.b.c.d", []string{"a", "b", "c", "d"}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input))
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
		{"obj.method()", &MethodCallExpr{Receiver: []string{"obj"}, Method: "method", Args: []any{}}},
		{"obj.method(1)", &MethodCallExpr{Receiver: []string{"obj"}, Method: "method", Args: []any{1}}},
		{"obj.method(1,2,3)", &MethodCallExpr{Receiver: []string{"obj"}, Method: "method", Args: []any{1, 2, 3}}},
		{`obj.method("hello","world")`, &MethodCallExpr{Receiver: []string{"obj"}, Method: "method", Args: []any{"hello", "world"}}},
		{"obj.method(1+2,3+4)", &MethodCallExpr{Receiver: []string{"obj"}, Method: "method", Args: []any{&AddExpr{Left: 1, Right: 2}, &AddExpr{Left: 3, Right: 4}}}},
		{"a.b.c.d.method()", &MethodCallExpr{Receiver: []string{"a", "b", "c", "d"}, Method: "method", Args: []any{}}},
		{"a.b.c.d.method(42)", &MethodCallExpr{Receiver: []string{"a", "b", "c", "d"}, Method: "method", Args: []any{42}}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tt.input))
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
