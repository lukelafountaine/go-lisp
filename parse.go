package main

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
	"os"
)

type ASTNode struct {
	nodeType  string
	intVal    int
	floatVal  float64
	symbolVal string
	args      []*ASTNode
}

func (n *ASTNode) traverse() {
	if n.nodeType == "int" {
		fmt.Print(n.intVal)
	} else if n.nodeType == "float" {
		fmt.Print(n.floatVal)
	} else if n.nodeType == "symbol" {
		fmt.Print(n.symbolVal)
	}

	for _, arg := range n.args {
		arg.traverse()
	}
}

func Parse(program string) (*ASTNode, error) {
	tokens, _, err := readFromTokens(tokenize(program))

	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

func tokenize(program string) []string {

	// add spaces to make the program splittable on whitespace
	program = strings.Replace(program, ")", " ) ", len(program))
	program = strings.Replace(program, "(", " ( ", len(program))

	return strings.Fields(program)
}

func readFromTokens(tokens []string) (*ASTNode, []string, error) {

	if len(tokens) == 0 {
		return &ASTNode{}, nil, errors.New("Unexpected EOF while parsing")
	}

	token, tokens := tokens[0], tokens[1:]
	if token == "(" {
		root := new(ASTNode)

		for len(tokens) > 0 && tokens[0] != ")" {
			args, newTokens, err := readFromTokens(tokens)
			tokens = newTokens
			if err != nil {
				return root, tokens, err
			}

			root.args = append(root.args, args)
		}

		if len(tokens) == 0 {
			return nil, nil, errors.New("Unexpected EOF while parsing")
		}

		tokens = tokens[1:]
		return root, tokens, nil

	} else if token == ")" {
		return nil, tokens, errors.New("Unexpected ')' while parsing")

	} else {
		return atom(token), tokens, nil
	}
}

func atom(value string) *ASTNode {

	intVal, err := strconv.Atoi(value)
	if err == nil {
		return &ASTNode{nodeType:"int", intVal:intVal}
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return &ASTNode{nodeType:"float", floatVal:floatVal}
	}

	fmt.Println("getting the symbol", value)
	return &ASTNode{nodeType:"symbol", symbolVal:value}
}

func main() {

	expression, err := Parse("(+ 1 2")

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	expression.traverse()

}