package main

import (
	"fmt"
	"strconv"
	"errors"
	"bufio"
	"os"
	"io/ioutil"

	"./lex"
)

// types
type Expression interface{}
type Symbol string
type String string
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
	tokens := tokenize(program)
	return readFromTokens(tokens)
}

func tokenize(program string) *[]lex.Token {
	scanner := lex.NewScanner(program)
	return lex.Scan(scanner)
}

func readFromTokens(tokens *[]lex.Token) (Expression, error) {

	if len(*tokens) == 0 {
		return nil, errors.New("Unexpected EOF while reading")
	}

	// pop the first token off
	token := (*tokens)[0]
	(*tokens) = (*tokens)[1:]

	switch token.Type {

	case lex.OpenParen:

		lst := make([]Expression, 0)

		for len(*tokens) > 0 && (*tokens)[0].Type != lex.CloseParen {

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

	case lex.CloseParen:
		return nil, errors.New(fmt.Sprintf("Syntax Error: Line %d: Unexpected ')", token.Line))

	default:
		return atom(token), nil
	}
}

func atom(token lex.Token) interface{} {

	switch token.Type {

	case lex.StringLiteral:
		return String(token.Text)

	case lex.NumberLiteral:
		num, err := strconv.ParseFloat(token.Text, 64)
		if err == nil {
			return Number(num)
		}
		return nil
	default:
		if token.Text == "#t" {
			return true
		}

		if token.Text == "#f" {
			return false
		}

		return Symbol(token.Text)
	}
}

func Eval(exp Expression, env *Scope) (result Expression, err error) {

	switch exp := exp.(type) {

	// variable reference
	case Symbol:
		var scope *Scope
		scope, err = getSymbol(Symbol(exp), env)

		if err == nil {
			result = scope.symbols[Symbol(exp)]
		}

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

				var scope *Scope
				scope, err = getSymbol(key, env)
				if err != nil {
					break
				}

				var value Expression
				value, err = Eval(exp[2], env)
				if err != nil {
					break
				}

				scope.symbols[key] = value

			case "begin":
				for _, i := range exp[1:] {
					result, err = Eval(i, env)

					if err != nil {
						break
					}
				}

			case "if":
				if len(exp) != 4 {
					err = errors.New("Syntax Error: Wrong number of arguments to 'if'")
					break
				}

				var condition Expression
				condition, err = Eval(exp[1], env)

				if err != nil {
					break
				}

				consequence := exp[2]
				alternative := exp[3]
				condResult := false

				switch condition := condition.(type) {
				case bool:
					if condition {
						condResult = true
					}

				case Number:
					if condition != 0 {
						condResult = true
					}

				case []Expression:
					if len(condition) > 0 {
						condResult = true
					}

				default:
					if condition != nil {
						condResult = true
					}
				}

				if condResult {
					result, err = Eval(consequence, env)
				} else {
					result, err = Eval(alternative, env)
				}

			case "lambda":
				if len(exp) != 3 {
					err = errors.New("Syntax Error: Wrong number of arguments to 'lambda'")
					break
				}
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

				var value Expression
				value, err = Eval(exp[2], env)

				if err != nil {
					break
				}

				env.symbols[key] = value

			// otherwise its a function
			default:
				result, err = applyFn(exp, env)
			}

		// otherwise its *probably* a function literal
		default:
			result, err = applyFn(exp, env)
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

func applyFn(fn []Expression, env *Scope) (result Expression, err error) {

	args := fn[1:]
	evaluated_args := make([]Expression, len(args))

	// get the function body
	var body Expression
	body, err = Eval(fn[0], env)

	// evaluate the arguments
	for i, op := range args {
		evaluated_args[i], err = Eval(op, env)
		if err != nil {
			return nil, err
		}
	}

	// call the function
	switch f := body.(type) {

	// built in functions
	case func(...Expression) (Expression, error):
		result, err = f(evaluated_args...)

	case func(...Expression) Expression:
		result = f(evaluated_args...)

	// user defined functions
	case Function:

		// make new environment with outer scope
		scope := &Scope{make(map[Symbol]Expression), f.env}

		switch params := f.params.(type) {

		case []Expression:
			if len(params) != len(evaluated_args) {
				err = errors.New(fmt.Sprintf("Wrong number of arguments to function. Expecting %s, got %s", len(params), len(args)))
			}

			for i, key := range params {
				scope.symbols[key.(Symbol)] = evaluated_args[i]
			}

		default:
			if len(evaluated_args) != 1 {
				err = errors.New(fmt.Sprintf("Wrong number of arguments to function. Expecting 1, got %s", len(evaluated_args)))
			}
			scope.symbols[params.(Symbol)] = evaluated_args[0]
		}

		result, err = Eval(f.body, scope)

	default:
		err = errors.New(fmt.Sprintf("%s is not callable", fn))
	}

	return result, err
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
