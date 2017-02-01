package main

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
)

func Parse(program string) (Expression, error) {
	// add 'begin' so it evaluates all expressions
	program = "(begin" + program + ")"
	tokens := tokenize(program)
	return readFromTokens(&tokens)
}

func tokenize(program string) []string {

	// add spaces to make the program splittable on whitespace
	program = strings.Replace(program, ")", " ) ", -1)
	program = strings.Replace(program, "(", " ( ", -1)

	// break on all whitespace
	return strings.Fields(program)
}

func readFromTokens(tokens *[]string) (Expression, error) {

	// pop the first token off
	token := (*tokens)[0]
	(*tokens) = (*tokens)[1:]

	switch token {

	case "(":

		L := make([]Expression, 0)

		for len(*tokens) > 0 && (*tokens)[0] != ")" {

			i, err := readFromTokens(tokens);

			if err != nil {
				return L, err
			} else if i != "" {
				L = append(L, i)
			}
		}

		if len(*tokens) == 0 {
			return nil, errors.New("Syntax Error: Missing closing parenthesis")
		}

		// pop off the closing paren
		*tokens = (*tokens)[1:]
		return L, nil

	case ")":
		return nil, errors.New("Syntax Error: Unexpected ')")

	default:
		return atom(token), nil

	}
}

func atom(value string) interface{} {

	num, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return Number(num)
	}

	return Symbol(value)
}

func getSymbol(symbol Symbol, env *Env) (Expression, error) {

	// get the symbol value if its there
	if val, ok := env.symbols[symbol]; ok {
		return val, nil
	}

	// otherwise check the next scope
	if env.outer != nil {
		return getSymbol(symbol, env.outer)
	}

	// otherwise the symbol is not found
	return nil, errors.New("'" + string(symbol) + "' is not defined")
}

func Eval(exp Expression, env *Env) (Expression, error) {

	switch val := exp.(type) {

	// the
	case Number:
		return val, nil

	case Symbol:
		return getSymbol(Symbol(val), env)

	case []Expression:

		// switch on the first word
		switch val[0].(Symbol) {

		case "begin":

			var value Expression
			var err error
			for _, i := range val[1:] {
				value, err = Eval(i, env)

				if err != nil {
					return nil, err
				}
			}
			return value, err

		case "if":
			if len(val) != 4 {
				return nil, errors.New("Syntax Error: Wrong number of arguments to 'if'")
			}

			condition, err := Eval(val[1], env)

			if err != nil {
				return nil, err
			}

			consequence := val[2]
			alternative := val[3]

			if condition.(bool) {
				return Eval(consequence, env)
			} else {
				return Eval(alternative, env)
			}

		case "lambda":
			return Function{val[1], val[2], env}, nil

		case "define":
			if len(val) != 3 {
				return nil, errors.New("Syntax Error: Wrong number of arguments to 'define'")
			}

			key := val[1].(Symbol)
			value, err := Eval(val[2], env)

			if err != nil {
				return nil, err
			}

			env.symbols[key] = value

		default:
			operands := val[1:]
			values := make([]Expression, len(operands))

			var err error
			// evaluate the operands
			for i, op := range operands {

				values[i], err = Eval(op, env)

				if err != nil {
					return nil, err
				}
			}

			// get the function from the name
			fn, err := Eval(val[0], env)
			if err != nil {
				return nil, err
			}

			// evaluate the function
			return apply(fn, values), nil
		}


	default:
		return nil, errors.New("Unknown Type")
	}

	return nil, nil
}

func apply(fn Expression, args []Expression) (value Expression) {

	value = nil

	switch f := fn.(type) {

	case func(...Expression) Expression:
		return f(args...)

	case Function:

		// make new environment with outer scope
		scope := &Env{make(map[Symbol]Expression), f.env}

		switch params := f.params.(type) {

		case []Expression:
			for i, key := range params {

				value, err := Eval(args[i], f.env)

				if err != nil {
					fmt.Println(err)
				}

				scope.symbols[key.(Symbol)] = value
			}

			result, err := Eval(f.body, scope)

			if err != nil {
				fmt.Println(err)
			}

			return result

		default:
			scope.symbols[params.(Symbol)] = args
			value = nil
		}

	default:
		value = nil
	}

	return value
}
