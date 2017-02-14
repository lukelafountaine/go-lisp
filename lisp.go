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
	replacer := strings.NewReplacer(")", " ) ", "(", " ( ", "\n", " ")
	return strings.Fields(replacer.Replace(program))
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

func Eval(exp Expression, env *Scope) (result Expression, err error) {

	switch exp := exp.(type) {

	// variable reference
	case Symbol:
		scope, err := getSymbol(Symbol(exp), env)

		if err != nil {
			break
		}

		result = scope.symbols[Symbol(exp)]

	case []Expression:

		// make sure we have something look at
		if len(exp) == 0 {
			break
		}

		switch t := exp[0].(type) {

		// switch on the first word
		case Symbol:

			switch t {

			case "quote":
				result = exp[1]

			case "set!":
				if len(exp) != 3 {
					err = errors.New("Syntax Error: Wrong number of arguments to 'set!'")
					break
				}

				key, ok := exp[1].(Symbol)
				if !ok {
					err = errors.New("Syntax Error: Cannot assign to a literal")
					break
				}

				scope, err := getSymbol(key, env)
				if err != nil {
					break
				}

				value, err := Eval(exp[2], env)
				if err != nil {
					break
				}

				scope.symbols[key] = value

			case "begin":
				for _, i := range exp[1:] {
					result, err = Eval(i, env)

					if err != nil {
						return nil, err
					}
				}

			case "if":
				if len(exp) != 4 {
					err = errors.New("Syntax Error: Wrong number of arguments to 'if'")
					break
				}

				condition, err := Eval(exp[1], env)

				if err != nil {
					break
				}

				consequence := exp[2]
				alternative := exp[3]
				boolean := false

				switch condition := condition.(type) {
				case bool:
					if condition {
						boolean = true
					}

				case Number:
					if condition != 0 {
						boolean  = true
					}

				case []Expression:
					if len(condition) > 0 {
						boolean = true
					}

				default:
					if condition != nil {
						boolean = true
					}
				}

				if boolean {
					result, err = Eval(consequence, env)
				} else {
					result, err = Eval(alternative, env)
				}

			case "lambda":
				result = Function{exp[1], exp[2], env}

			case "define":
				if len(exp) != 3 {
					err = errors.New("Syntax Error: Wrong number of arguments to 'define'")
					break
				}

				key, ok := exp[1].(Symbol)
				if !ok {
					err = errors.New("Syntax Error: Cannot assign to a literal")
					break
				}

				value, err := Eval(exp[2], env)

				if err != nil {
					break
				}

				env.symbols[key] = value

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
				result, err = apply(fn, operands)
			}

		// otherwise its probably a function literal
		default:
			fn, err := Eval(t, env)

			if err != nil {
				return nil, err
			}

			return apply(fn, exp[1:])
		}

	// constant literal
	default:
		result = exp
	}

	return result, err
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
			if len(params) != len(args) {
				value = nil
				err = errors.New(fmt.Sprintf("Wrong number of arguments to function. Expecting %s, got %s", len(params), len(args)))
			}

			for i, key := range params {
				scope.symbols[key.(Symbol)] = args[i]
			}

		default:
			if len(args) != 1 {
				value = nil
				err = errors.New(fmt.Sprintf("Wrong number of arguments to function. Expecting 1, got %s", len(args)))
			}
			scope.symbols[params.(Symbol)] = args[0]
		}

		value, err = Eval(f.body, scope)

	default:
		err = errors.New(fmt.Sprintf("%s is not callable", fn))
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
	}
}

func Repl(scope *Scope) {
	scanner := bufio.NewScanner(os.Stdin)
	for fmt.Print("> "); scanner.Scan(); fmt.Print("> ") {
		Run(scanner.Text(), scope)
	}
}

func main() {

	scope := NewEnv()

	// load the standard library
	program, _ := ioutil.ReadFile("stdlib.lisp")
	Run(string(program), scope)

	// evaluate any files provided
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
