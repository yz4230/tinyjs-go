%{
package subjs

func NewParser() yyParser {
	return yyNewParser()
}

func init() {
    yyErrorVerbose = true
}

%}

%union {
	val any
	literal string
}

%token<val> STRING NUMBER IDENT
%type<val> program expr factor args call ident idents method_call

%left '+'

%%

program:
	expr {
        $$ = $1
        if yylex, ok := yylex.(*Lexer); ok {
            yylex.Result = $$
        }
	}

expr:
    factor { $$ = $1 }
|   expr '+' expr { $$ = &AddExpr{Left: $1, Right: $3} }

factor:
    STRING { $$ = StringValue($1.(string)) }
|   NUMBER { $$ = NumberValue($1.(int)) }
|   call { $$ = $1 }
|   idents { $$ = $1 }
|   method_call { $$ = $1 }
|   '(' expr ')' { $$ = $2 }

args:
    { $$ = []any{} }
|   expr { $$ = []any{$1} }
|   args ',' expr { $$ = append($1.([]any), $3) }

call:
    ident '(' args ')' { 
        $$ = &CallExpr{Name: $1.(Ident), Args: $3.([]any)} 
    }

ident:
    IDENT { $$ = Ident($1.(string)) }

idents:
    ident { $$ = []any{$1.(Ident)} }
|   idents '.' ident { $$ = append($1.([]any), $3.(Ident)) }

method_call:
    idents '.' ident '(' args ')' {
        $$ = &MethodCallExpr{
            Receiver: $1.([]any),
            Method:   $3.(Ident),
            Args:     $5.([]any),
        }
    }
|   factor '.' ident '(' args ')' {
        $$ = &MethodCallExpr{
            Receiver: []any{$1},
            Method:   $3.(Ident),
            Args:     $5.([]any),
        }
    }
