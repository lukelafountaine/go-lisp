package lex

import (
	"fmt"
	"unicode/utf8"
	"unicode"
	"io"
	"strings"
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

const operators = "+-*/=<>%"

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
	reader io.ByteReader
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

	if int(l.pos) == len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *Scanner) peek() rune {
	next := l.next()
	l.backup()
	return next
}

func (l *Scanner) backup() {
	l.pos -= l.width
}

func (l *Scanner) ignore() {
	l.start = l.pos
}

func start(l *Scanner) scanFn {

	r := l.next()

	switch {

	case r == eof:
		close(l.tokens)
		return nil

	case r == '\n':
		l.emitToken(NewLine)
		return start

	case unicode.IsSpace(r):
		return lexSpace

	case r == '(':
		l.emitToken(OpenParen)
		return start

	case r == ')':
		l.emitToken(CloseParen)
		return start

	// skip comments
	case r == ';':
		return lexComment

	case r == '"':
		l.ignore()
		return lexString

	case r == '.':
		return lexDigitsOnly

	case r == '-':
		next := l.peek()

		switch {
		case unicode.IsNumber(next), next == '.':
			return lexNumber

		default:
			return lexSymbol
		}

	case unicode.IsNumber(r):
		return lexNumber

	default:
		return lexSymbol
	}
}

func lexString(l *Scanner) scanFn {

	switch l.next() {
	case '"':
		l.backup()
		l.emitToken(StringLiteral)
		l.next()
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

	case r == eof:
		l.emitToken(EOF)
		return start

	case !(unicode.IsLetter(r) || unicode.IsNumber(r) || strings.Contains(operators, string(r))):
		l.backup()
		l.emitToken(Symbol)
		return start

	default:
		return lexSymbol
	}
}

func lexNumber(l *Scanner) scanFn {
	r := l.next()

	switch {
	case unicode.IsNumber(r):
		return lexNumber

	case r == '.':
		return lexDigitsOnly

	default:
		l.backup()
		l.emitToken(NumberLiteral)
		return start
	}
}

func lexDigitsOnly(l *Scanner) scanFn {
	r := l.next()

	switch {
	case unicode.IsNumber(r):
		return lexDigitsOnly

	default:
		l.backup()
		l.emitToken(NumberLiteral)
		return start
	}
}

func lexComment(l *Scanner) scanFn {
	for {
		next := l.peek()
		if next == '\n' || next == eof {
			break
		}
		l.next()
	}

	l.emitToken(Comment)
	return start
}

func lexSpace(l *Scanner) scanFn {

	for unicode.IsSpace(l.peek()) {
		l.next()
	}

	l.ignore()
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
		if token.Type != NewLine {
			tokens = append(tokens, token)
		}
	}

	return &tokens
}
