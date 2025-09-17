package subjs

import (
	"testing"
	"unicode/utf8"
)

func TestNext_EmptyInput(t *testing.T) {
	l := NewLexer([]byte{})
	l.next()
	if l.ch != 0 {
		t.Fatalf("expected ch==0 for empty input, got %q", l.ch)
	}
	if l.offset != 0 {
		t.Fatalf("expected offset==0 for empty input, got %d", l.offset)
	}
}

func TestNext_ASCIISequence(t *testing.T) {
	l := NewLexer([]byte("ab"))
	l.next()
	if l.ch != 'a' {
		t.Fatalf("expected first ch='a', got %q", l.ch)
	}
	if l.offset != 1 {
		t.Fatalf("expected offset==1 after first next, got %d", l.offset)
	}

	l.next()
	if l.ch != 'b' {
		t.Fatalf("expected second ch='b', got %q", l.ch)
	}
	if l.offset != 2 {
		t.Fatalf("expected offset==2 after second next, got %d", l.offset)
	}

	l.next()
	if l.ch != 0 {
		t.Fatalf("expected ch==0 after consuming all input, got %q", l.ch)
	}
	if l.offset != 2 {
		t.Fatalf("expected offset to remain at len(input)==2, got %d", l.offset)
	}
}

func TestNext_MultiByteRunes(t *testing.T) {
	// '世' is a 3-byte rune, followed by 'a'
	s := "世a"
	l := NewLexer([]byte(s))

	l.next()
	if l.ch != '世' {
		t.Fatalf("expected first ch='世', got %q", l.ch)
	}
	if l.offset != utf8.RuneLen('世') {
		t.Fatalf("expected offset==%d after reading '世', got %d", utf8.RuneLen('世'), l.offset)
	}

	l.next()
	if l.ch != 'a' {
		t.Fatalf("expected second ch='a', got %q", l.ch)
	}
	if l.offset != len(s) {
		t.Fatalf("expected offset==len(%q)==%d after reading 'a', got %d", s, len(s), l.offset)
	}

	l.next()
	if l.ch != 0 {
		t.Fatalf("expected ch==0 after exhausting input, got %q", l.ch)
	}
}

func TestNext_EmojiFourByte(t *testing.T) {
	// '😊' is a 4-byte rune
	s := "😊x"
	l := NewLexer([]byte(s))

	l.next()
	if l.ch != '😊' {
		t.Fatalf("expected first ch='😊', got %q", l.ch)
	}
	if l.offset != utf8.RuneLen('😊') {
		t.Fatalf("expected offset==%d after reading '😊', got %d", utf8.RuneLen('😊'), l.offset)
	}

	l.next()
	if l.ch != 'x' {
		t.Fatalf("expected second ch='x', got %q", l.ch)
	}
	if l.offset != len(s) {
		t.Fatalf("expected offset==len(%q)==%d after reading 'x', got %d", s, len(s), l.offset)
	}
}
