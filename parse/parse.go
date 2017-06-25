package parse

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lukelafountaine/go-lisp/scan"
	"github.com/lukelafountaine/go-lisp/types"
)

func Parse(program string) (types.Expression, error) {
	tokens := tokenize(program)
	return readFromTokens(tokens)
}

func tokenize(program string) *[]scan.Token {

	scanner := scan.NewScanner(program)

	tokens := make([]scan.Token, 0)

	for token := scanner.NextToken(); token.Type != scan.EOF; token = scanner.NextToken() {
		tokens = append(tokens, token)
	}

	return &tokens
}

func readFromTokens(tokens *[]scan.Token) (types.Expression, error) {

	if len(*tokens) == 0 {
		return nil, errors.New("Unexpected EOF while reading")
	}

	// pop the first token off
	token := (*tokens)[0]
	(*tokens) = (*tokens)[1:]

	switch token.Type {

	case scan.OpenParen:

		lst := make([]types.Expression, 0)

		for len(*tokens) > 0 && (*tokens)[0].Type != scan.CloseParen {

			token = (*tokens)[0]
			i, err := readFromTokens(tokens)

			if err != nil {
				return lst, err
			} else if i != "" {
				lst = append(lst, i)
			}
		}

		if len(*tokens) == 0 {
			return nil, errors.New(fmt.Sprintf("Syntax Error: Line %d: Missing closing parenthesis", token.Line))
		}

		// pop off the closing paren
		*tokens = (*tokens)[1:]
		return lst, nil

	case scan.CloseParen:
		return nil, errors.New(fmt.Sprintf("Syntax Error: Line %d: Unexpected ')", token.Line))

	default:
		return atom(token), nil
	}
}

func atom(token scan.Token) interface{} {

	switch token.Type {

	case scan.StringLiteral:
		return types.String(token.Text)

	case scan.NumberLiteral:
		num, err := strconv.ParseFloat(token.Text, 64)
		if err == nil {
			return types.Number(num)
		}
		return nil
	default:
		if token.Text == "#t" {
			return true
		}

		if token.Text == "#f" {
			return false
		}

		return types.Symbol(token.Text)
	}
}
