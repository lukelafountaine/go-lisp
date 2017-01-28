package main

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
	"os"
	"bufio"
)

type AST struct {
	nodeType string
	val      interface{}
	children []*AST
}

func (node *AST) String() string {

	switch node.nodeType {

	case "int":
		return fmt.Sprintf(" (%v, int) ", node.val.(int))
	case "float":
		return fmt.Sprintf(" (%v, float) ",node.val.(float64))
	case "string":
		return fmt.Sprintf(" (%v, string) ", node.val.(string))
	case "expression":
		result := ""
		for _, arg := range node.children {
			result += arg.String()
		}
		return result
	default:
		return "nothing"
	}
}

func Parse(program string) (*AST, error) {
	tokens, _, err := readFromTokens(tokenize(program))

	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

func Eval(node *AST, env *Env) interface{} {

	switch node.nodeType {

	case "symbol":
		return (*env)[node.val.(string)]

	case "int", "float":
		return node.val

	case "expression":
		fn := Eval(node.children[0], env).(func (int, int) int)
		initial := Eval(node.children[1], env).(int)

		for i := 2; i < len(node.children); i++ {
			next := Eval(node.children[i], env).(int)
			initial = fn(initial, next)
		}
		return initial

	default:
		return nil
	}
}

func tokenize(program string) []string {

	// add spaces to make the program splittable on whitespace
	program = strings.Replace(program, ")", " ) ", len(program))
	program = strings.Replace(program, "(", " ( ", len(program))

	return strings.Fields(program)
}

func readFromTokens(tokens []string) (*AST, []string, error) {

	if len(tokens) == 0 {
		return &AST{}, nil, errors.New("Unexpected EOF while parsing")
	}

	// pop the first token off
	token, tokens := tokens[0], tokens[1:]

	if token == "(" {
		root := new(AST)
		root.nodeType = "expression"

		for len(tokens) > 0 && tokens[0] != ")" {
			args, newTokens, err := readFromTokens(tokens)
			tokens = newTokens
			if err != nil {
				return root, tokens, err
			}

			root.children = append(root.children, args)
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

func atom(value string) *AST {

	node := AST{}

	intVal, err := strconv.Atoi(value)
	if err == nil {
		node.val = intVal
		node.nodeType = "int"
		return &node
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err == nil {
		node.val = floatVal
		node.nodeType = "float"
		return &node
	}

	node.val = value
	node.nodeType = "symbol"
	return &node
}

func Repl() {
	for {
		reader := bufio.NewReader(os.Stdin)
		program, _ := reader.ReadString('\n')

		exp, err := Parse(program)
		env := NewEnv()
		result := Eval(exp, env)
		fmt.Println(result)
		//fmt.Println(exp)

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}

func main() {
	Repl()
}