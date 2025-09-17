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
%type<val> program expr factor args call idents method_call

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
|   call { $$ = $1 }
|   idents { $$ = $1 }
|   method_call { $$ = $1 }
|   expr '+' expr { $$ = &AddExpr{Left: $1, Right: $3} }

factor:
    STRING { $$ = $1 }
|   NUMBER { $$ = $1 }
|   '(' expr ')' { $$ = $2 }

args:
    { $$ = []any{} }
|   expr { $$ = []any{$1} }
|   args ',' expr { $$ = append($1.([]any), $3) }

call:
    IDENT '(' args ')' { $$ = &CallExpr{Name: $1.(string), Args: $3.([]any)} }

idents:
    IDENT { $$ = []string{$1.(string)} }
|   idents '.' IDENT { $$ = append($1.([]string), $3.(string)) }

method_call:
    idents '.' IDENT '(' args ')' {
        $$ = &MethodCallExpr{Receiver: $1.([]string), Method: $3.(string), Args: $5.([]any)}
    }
