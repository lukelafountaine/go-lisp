package scan

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	Line int
	Type Type
	Text string
}

type Type int

const (
	EOF Type = iota - 1
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

type Scanner struct {
	reader io.ByteReader
	tokens chan Token
	state  stateFn
	input  string
	line   int
	start  int
	width  int
	pos    int
}

type stateFn func(*Scanner) stateFn

func NewScanner(input string) *Scanner {
	return &Scanner{
		tokens: make(chan Token, 2),
		state:  start,
		input:  input,
		line:   0,
		start:  0,
		pos:    0,
	}
}

const eof = -1

func (s *Scanner) emitToken(t Type) {

	if t == NewLine {
		s.line += 1
	}
	s.tokens <- Token{s.line, t, s.input[s.start:s.pos]}
	s.start = s.pos
	s.width = 0
}

// next returns the next rune in the input.
func (s *Scanner) next() rune {

	if int(s.pos) == len(s.input) {
		s.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(s.input[s.pos:])
	s.width = w
	s.pos += s.width
	return r
}

func (s *Scanner) peek() rune {
	next := s.next()
	s.backup()
	return next
}

func (s *Scanner) backup() {
	s.pos -= s.width
}

func (s *Scanner) ignore() {
	s.start = s.pos
}

func start(s *Scanner) stateFn {

	r := s.next()

	switch {

	case r == eof:
		close(s.tokens)
		return nil

	case r == '\n':
		s.emitToken(NewLine)
		return start

	case unicode.IsSpace(r):
		return lexSpace

	case r == '(':
		s.emitToken(OpenParen)
		return start

	case r == ')':
		s.emitToken(CloseParen)
		return start

	// skip comments
	case r == ';':
		return lexComment

	case r == '"':
		s.ignore()
		return lexString

	case r == '.':
		return lexDigitsOnly

	case r == '-':
		next := s.peek()

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

func lexString(s *Scanner) stateFn {

	switch s.next() {
	case '"':
		s.backup()
		s.emitToken(StringLiteral)
		s.next()
		return start

	case eof:
		s.emitToken(EOF)
		return nil
	}

	return lexString
}

func lexSymbol(s *Scanner) stateFn {

	r := s.next()
	switch {

	case !(unicode.IsLetter(r) || unicode.IsNumber(r) || strings.Contains(operators, string(r))):
		s.backup()
		s.emitToken(Symbol)
		return start

	default:
		return lexSymbol
	}
}

func lexNumber(s *Scanner) stateFn {
	r := s.next()

	switch {
	case unicode.IsNumber(r):
		return lexNumber

	case r == '.':
		return lexDigitsOnly

	default:
		s.backup()
		s.emitToken(NumberLiteral)
		return start
	}
}

func lexDigitsOnly(s *Scanner) stateFn {
	r := s.next()

	switch {
	case unicode.IsNumber(r):
		return lexDigitsOnly

	default:
		s.backup()
		s.emitToken(NumberLiteral)
		return start
	}
}

func lexComment(s *Scanner) stateFn {

	for {
		next := s.peek()
		if next == '\n' {
			s.emitToken(NewLine)
			break
		}
		if next == eof {
			break
		}
		s.next()
	}

	s.emitToken(Comment)
	return start
}

func lexSpace(s *Scanner) stateFn {

	for unicode.IsSpace(s.peek()) {
		s.next()
	}

	s.ignore()
	return start
}

func (s *Scanner) NextToken() Token {

	for s.state != nil {
		select {

		case token := <-s.tokens:
			switch token.Type {

			case NewLine, Comment:

			default:
				return token
			}

		default:
			s.state = s.state(s)
		}
	}

	return Token{0, EOF, ""}
}
