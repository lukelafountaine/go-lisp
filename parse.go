package main

import (
	"fmt"
	"strings"
	"strconv"
	"reflect"
)

type Expression interface{}
type Symbol string
type Number float64

func Parse(program string) Expression {
	tokens := tokenize(program)
	tree := readFromTokens(&tokens)
	return tree
}


func tokenize(program string) []string {

	// add spaces to make the program splittable on whitespace
	program = strings.Replace(program, ")", " ) ", -1)
	program = strings.Replace(program, "(", " ( ", -1)

	return strings.Fields(program)
}

func Eval(exp Expression, env *Env) Expression {

	switch val := exp.(type) {

	case Number:
		return val

	case Symbol:
		return env.symbols[Symbol(val)]

	case []Expression:

		switch start := val[0].(Symbol); start {

		case "define":
			if len(val) < 3 {
				fmt.Println("Syntax Error: Wrong number of arguments to 'define'")
			}

			key := val[1].(Symbol)
			value := val[2]

			env.symbols[key] = value

		default:
			operands := val[1:]
			values := make([]Expression, len(operands))

			// evaluate the operands
			for i, op := range operands {
				values[i] = Eval(op, env)
			}

			fn := (*env).symbols[val[0].(Symbol)].(func (...Expression) Expression)
			return fn(values...)
		}


	default:
		fmt.Println("Unknown Type: ", reflect.TypeOf(val))
	}

	return nil
}

func readFromTokens(tokens *[]string) Expression {

	// pop the first token off
	token := (*tokens)[0]
	(*tokens) = (*tokens)[1:]

	switch token {

	case "(":

		L := make([]Expression, 0)

		for (*tokens)[0] != ")" {

			if i := readFromTokens(tokens); i != "" {
				L = append(L, i)
			}
		}

		// pop off the closing paren
		*tokens = (*tokens)[1:]
		return L

	case ")":
		fmt.Println("Syntax Error")

	default:
		return atom(token)

	}

	return nil
}

func atom(value string) interface{} {

	num, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return Number(num)
	}

	return Symbol(value)
}
