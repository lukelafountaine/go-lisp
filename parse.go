package main

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
)

// types
type Expression interface{}
type Symbol string
type Number float64
type Function struct {
	params, body Expression
	env          *Scope
}
type Scope struct {
	symbols map[Symbol]Expression
	outer   *Scope
}

func Parse(program string) (Expression, error) {
	// add 'begin' so it evaluates all expressions
	//program = "(begin" + program + ")"
	tokens := tokenize(program)
	return readFromTokens(&tokens)
}

func tokenize(program string) []string {

	// add spaces to make the program splittable on whitespace
	program = strings.Replace(program, "\n", " ", -1)
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

func Eval(exp Expression, env *Scope) (Expression, error) {

	switch exp := exp.(type) {

	case Number, Function:
		return exp, nil

	case Symbol:
		return getSymbol(Symbol(exp), env)

	case []Expression:

		// make sure we have something look at
		if len(exp) == 0 {
			return nil, nil
		}

		switch t := exp[0].(type) {

		// switch on the first word
		case Symbol:

			switch t {

			case "quote":
				return exp[1:], nil

			case "set!":
				if len(exp) != 3 {
					return nil, errors.New("Syntax Error: Wrong number of arguments to 'set!'")
				}

				key := exp[1].(Symbol)
				if _, ok := (*env).symbols[key]; !ok {
					return nil, errors.New(fmt.Sprintf("Symbol '%s' not defined", key))
				}

				value, err := Eval(exp[2], env)
				if err != nil {
					return nil, err
				}

				env.symbols[key] = value
				return nil, nil


			case "begin":

				var value Expression
				var err error
				for _, i := range exp[1:] {
					value, err = Eval(i, env)

					if err != nil {
						return nil, err
					}
				}
				return value, err

			case "if":
				if len(exp) != 4 {
					return nil, errors.New("Syntax Error: Wrong number of arguments to 'if'")
				}

				condition, err := Eval(exp[1], env)

				if err != nil {
					return nil, err
				}

				consequence := exp[2]
				alternative := exp[3]

				if condition.(bool) {
					return Eval(consequence, env)
				} else {
					return Eval(alternative, env)
				}

			case "lambda":
				return Function{exp[1], exp[2], env}, nil

			case "define":
				if len(exp) != 3 {
					return nil, errors.New("Syntax Error: Wrong number of arguments to 'define'")
				}

				key := exp[1].(Symbol)
				value, err := Eval(exp[2], env)

				if err != nil {
					return nil, err
				}

				env.symbols[key] = value
				return nil, nil

			// the default case is that it is a user defined function call
			default:
				operands := exp[1:]
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
				fn, err := Eval(exp[0], env)
				if err != nil {
					return nil, err
				}

				// evaluate the function
				return apply(fn, values)


			}

		// the default is that the first elemet in the list is a function literal
		default:
			fn, err := Eval(t, env)

			if err != nil {
				return nil, err
			}

			return apply(fn, exp[1:])
		}

	default:
		return nil, errors.New(fmt.Sprintf("Unknown Type: %T for var %s", exp, exp))
	}

	return nil, nil
}

func getSymbol(symbol Symbol, env *Scope) (Expression, error) {

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

func apply(fn Expression, args []Expression) (value Expression, err error) {

	value = nil
	err = nil

	switch f := fn.(type) {

	// built in functions
	case func(...Expression) (Expression, error):
		return f(args...)

	// user defined functions
	case Function:

		// make new environment with outer scope
		scope := &Scope{make(map[Symbol]Expression), f.env}

		switch params := f.params.(type) {

		case []Expression:
			for i, key := range params {

				value, err := Eval(args[i], f.env)

				if err != nil {
					fmt.Println(err)
				}

				scope.symbols[key.(Symbol)] = value
			}

			return Eval(f.body, scope)

		default:
			scope.symbols[params.(Symbol)] = args
			value = nil
		}

	default:
		value = nil
	}

	return value, err
}
