package main

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
	"bufio"
	"os"
	"io/ioutil"
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
	tokens := tokenize("(begin " + program + ")")
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

	if len(*tokens) == 0 {
		return nil, errors.New("Unexpected EOF while reading")
	}

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

	// variable reference
	case Symbol:
		scope, err := getSymbol(Symbol(exp), env)

		if err != nil {
			return nil, err
		}

		return scope.symbols[Symbol(exp)], nil

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
				return exp[1], nil

			case "set!":
				if len(exp) != 3 {
					return nil, errors.New("Syntax Error: Wrong number of arguments to 'set!'")
				}

				key, ok := exp[1].(Symbol)
				if !ok {
					return nil, errors.New("Syntax Error: Cannot assign to a literal")
				}

				scope, err := getSymbol(key, env)
				if err != nil {
					return nil, err
				}

				value, err := Eval(exp[2], env)
				if err != nil {
					return nil, err
				}

				scope.symbols[key] = value
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
				result := false

				switch condition := condition.(type) {
				case bool:
					if condition {
						result = true
					}

				case Number:
					if condition != 0 {
						result  = true
					}

				case []Expression:
					if len(condition) > 0 {
						result = true
					}

				default:
					if condition != nil {
						result = true
					}
				}

				if result {
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

				key, ok := exp[1].(Symbol)
				if !ok {
					return nil, errors.New("Syntax Error: Cannot assign to a literal")
				}

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

				// get the function from the name
				fn, err := Eval(exp[0], env)
				if err != nil {
					return nil, err
				}

				// evaluate the operands
				for i, op := range operands {

					values[i], err = Eval(op, env)

					if err != nil {
						return nil, err
					}
				}

				// apply the function
				return apply(fn, values)
			}

		// the default is that the first element in the list is a function literal
		default:
			fn, err := Eval(t, env)

			if err != nil {
				return nil, err
			}

			return apply(fn, exp[1:])
		}

	// constant literal
	default:
		return exp, nil
	}
}

func getSymbol(symbol Symbol, env *Scope) (*Scope, error) {

	// get the symbol value if its there
	if _, ok := env.symbols[symbol]; ok {
		return env, nil
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

	case func(...Expression) Expression:
		return f(args...), nil

	// user defined functions
	case Function:

		// make new environment with outer scope
		scope := &Scope{make(map[Symbol]Expression), f.env}

		switch params := f.params.(type) {

		case []Expression:
			for i, key := range params {
				scope.symbols[key.(Symbol)] = args[i]
			}

		default:
			scope.symbols[params.(Symbol)] = args
		}

		value, err = Eval(f.body, scope)

	default:
		fmt.Println(fn, "is not callable")
		value = nil
	}

	return value, err
}

func Run(program string, scope *Scope) {

	// read
	exp, err := Parse(program)

	if err != nil {
		fmt.Println(err)
		return
	}

	// evaluate
	result, err := Eval(exp, scope)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print
	if result != nil {
		fmt.Println(result)
	} else {
		fmt.Println("ok.")
	}
}

func Repl(scope *Scope) {

	scanner := bufio.NewScanner(os.Stdin)
	var program string

	for fmt.Print("> "); scanner.Scan(); fmt.Print("> ") {

		program = scanner.Text()
		if program == "exit" {
			fmt.Println("bye!")
			os.Exit(0)
		}

		Run(program, scope)
	}
}

func main() {

	scope := NewEnv()
	for _, file := range os.Args[1:] {
		program, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		Run(string(program), scope)
	}

	Repl(scope)
}
