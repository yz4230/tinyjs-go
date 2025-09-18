package subjs

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/samber/lo"
)

type Lexer struct {
	input  []rune
	offset int
	ch     rune
	Result any
	Err    error
}

func NewLexer(input []rune) *Lexer {
	return &Lexer{input: input, offset: 0, ch: input[0]}
}

func (l *Lexer) Lex(lval *yySymType) int {
	l.skipWhitespace()
	l.next()

	if l.ch == '+' || l.ch == '(' || l.ch == ')' || l.ch == '.' || l.ch == ',' {
		return int(l.ch)
	}
	if l.ch == '"' || l.ch == '\'' {
		var sb strings.Builder
		stop := l.ch
		for {
			l.next()
			if l.ch == stop {
				break
			}
			sb.WriteRune(l.ch)
		}

		lval.literal = string(stop) + sb.String() + string(stop)
		lval.val = sb.String()
		return STRING
	}
	if unicode.IsDigit(l.ch) {
		var sb strings.Builder
		sb.WriteRune(l.ch)
		for unicode.IsDigit(l.peek()) {
			l.next()
			sb.WriteRune(l.ch)
		}
		lval.literal = sb.String()
		lval.val = lo.Must(strconv.Atoi(lval.literal))
		return NUMBER
	}
	if unicode.IsLetter(l.ch) {
		var sb strings.Builder
		sb.WriteRune(l.ch)
		for unicode.IsLetter(l.peek()) {
			l.next()
			sb.WriteRune(l.ch)
		}
		lval.literal = sb.String()
		lval.val = lval.literal
		return IDENT
	}

	return 0
}

func (l *Lexer) next() {
	if l.offset >= len(l.input) {
		l.ch = 0
		return
	}
	l.ch = l.input[l.offset]
	l.offset++
}

func (l *Lexer) peek() rune {
	if l.offset >= len(l.input) {
		return 0
	}
	return l.input[l.offset]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.peek()) {
		l.next()
	}
}

func (l *Lexer) Error(s string) {
	l.Err = errors.New(s)
}
