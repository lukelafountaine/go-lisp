package lex

import (
	"fmt"
	"unicode/utf8"
	"unicode"
)

type Token struct {
	Type Type
	Text string
}

type Type int

const (
	EOF Type = iota-1
	Comment
	OpenParen
	CloseParen
	NewLine
	StringLiteral
	NumberLiteral
	Symbol
)

func (t Token) String() string {
	switch t.Type {

	case EOF:
		return "<EOF>"

	case NewLine:
		return "<NewLine>"

	case OpenParen:
		return "<Open Paren>"

	case CloseParen:
		return "<Close Paren>"

	case Comment, StringLiteral, Symbol, NumberLiteral:
		return fmt.Sprintf("<%s: %s>", t.Type, t.Text)

	default:
		return "Dont know this type"
	}
}

type stateFn func(*Scanner) stateFn

type Scanner struct {
	tokens chan Token
	state  stateFn
	input  string
	line   int
	start  int
	width int
	pos    int
}

type scanFn func(*Scanner) scanFn

func NewScanner(input string) *Scanner {
	return &Scanner{
		tokens : make(chan Token, 0),
		input : input,
		line : 0,
		start : 0,
		pos : 0,
	}
}

const eof = -1

func (l *Scanner) emitToken(t Type) {

	if t == NewLine {
		l.line += 1
	}
	l.tokens <- Token{t, l.input[l.start:l.pos]}
	l.start = l.pos
	l.width = 0
}

// next returns the next rune in the input.
func (l *Scanner) next() rune {

	if int(l.pos) < len(l.input) {
		r, w := utf8.DecodeRuneInString(l.input[l.pos:])
		l.width = w
		l.pos += l.width
		return r
	}

	return eof
}

// next returns the next rune in the input.
func (l *Scanner) backup() {
	l.pos -= l.width
}

func start(l *Scanner) scanFn {

	switch l.next() {

	case eof:
		close(l.tokens)
		return nil

	case '\n':
		l.emitToken(NewLine)
		return start

	case ' ', '\t':
		return lexSpace

	case '(':
		l.emitToken(OpenParen)
		return start

	case ')':
		l.emitToken(CloseParen)
		return start

	// skip comments
	case ';':
		return lexComment

	case '"':
		return lexString

	default:
		return lexSymbol
	}
}

func lexString(l *Scanner) scanFn {

	switch l.next() {
	case '"':
		l.emitToken(StringLiteral)
		return start

	case eof:
		l.emitToken(EOF)
		return nil
	}

	return lexString
}

func lexSymbol(l *Scanner) scanFn {

	r := l.next()
	switch {
	case r == eof, r == '\n', r == ' ':
		l.backup()
		l.emitToken(Symbol)
		return start

	case unicode.IsLetter(r) || unicode.IsNumber(r):
		return lexSymbol

	default:
		return start
	}
}

func lexComment(l *Scanner) scanFn {

	for {
		next := l.next()
		if next == '\n' || next == eof {
			break
		}
	}

	l.backup()
	l.emitToken(Comment)
	return start
}

func lexSpace(l *Scanner) scanFn {
	for unicode.IsSpace(l.next()) {}
	l.backup()
	return start
}

func Scan(l *Scanner) *[]Token {

	go (func() {
		for fn := start; fn != nil; {
			fn = fn(l)
		}
	})()

	tokens := make([]Token, 0)
	for token := range l.tokens {
		fmt.Println(token)
		tokens = append(tokens, token)
	}
	return &tokens
}
