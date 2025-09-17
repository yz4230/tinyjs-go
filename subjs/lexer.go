package subjs

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/samber/lo"
)

type Lexer struct {
	input  []byte
	offset int
	ch     rune
	Result any
	Err    error
}

func NewLexer(input []byte) *Lexer {
	return &Lexer{input: input, offset: 0}
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
		for {
			if unicode.IsDigit(l.peek()) {
				l.next()
				sb.WriteRune(l.ch)
			} else {
				break
			}
		}
		lval.literal = sb.String()
		lval.val = lo.Must(strconv.Atoi(lval.literal))
		return NUMBER
	}
	if unicode.IsLetter(l.ch) {
		var sb strings.Builder
		sb.WriteRune(l.ch)
		for {
			if unicode.IsLetter(l.peek()) {
				l.next()
				sb.WriteRune(l.ch)
			} else {
				break
			}
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
	ch := l.input[l.offset]
	if ch < utf8.RuneSelf {
		l.offset++
		l.ch = rune(ch)
		return
	}
	r, size := utf8.DecodeRune(l.input[l.offset:])
	l.offset += size
	l.ch = r
}

func (l *Lexer) peek() rune {
	if l.offset >= len(l.input) {
		return 0
	}
	ch := l.input[l.offset]
	if ch < utf8.RuneSelf {
		return rune(ch)
	}
	return lo.Must(utf8.DecodeRune(l.input[l.offset:]))
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.next()
	}
}

func (l *Lexer) Error(s string) {
	l.Err = errors.New(s)
}
