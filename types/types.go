package types

import (
	"errors"
	"fmt"
	"reflect"
)

// types
type Expression interface{}
type Symbol string
type String string
type Number float64
type Function struct {
	Params, Body Expression
	Env          *Scope
}
type Scope struct {
	Symbols map[Symbol]Expression
	Outer   *Scope
}

// type check for functions that expect arguments of type Number
func numberType(f func(...Expression) Expression) func(...Expression) (Expression, error) {

	return func(args ...Expression) (Expression, error) {
		for _, arg := range args {
			if _, ok := arg.(Number); !ok {
				return nil, errors.New(fmt.Sprintf("Type Error: Recieved %T, Expected Number", arg))
			}
		}
		return f(args...), nil
	}
}

// type check for functions that expect arguments of type Boolean
func boolType(f func(...Expression) Expression) func(...Expression) (Expression, error) {

	return func(args ...Expression) (Expression, error) {
		for _, arg := range args {
			if _, ok := arg.(bool); !ok {
				return nil, errors.New(fmt.Sprintf("Type Error: Recieved %T, Expected Boolean", arg))
			}
		}
		return f(args...), nil
	}
}

func add(args ...Expression) Expression {

	initial := args[0].(Number)
	for _, num := range args[1:] {
		initial += num.(Number)
	}

	return initial
}

func sub(args ...Expression) Expression {

	initial := args[0].(Number)
	for _, num := range args[1:] {
		initial -= num.(Number)
	}

	return initial
}

func mult(args ...Expression) Expression {

	initial := args[0].(Number)
	for _, num := range args[1:] {
		initial *= num.(Number)
	}

	return initial
}

func div(args ...Expression) Expression {

	initial := args[0].(Number)
	for _, num := range args[1:] {
		initial /= num.(Number)
	}

	return initial
}

func mod(args ...Expression) Expression {
	return Number(int(args[0].(Number)/1) % int(args[1].(Number)/1))
}

func lt(args ...Expression) Expression {
	return args[0].(Number) < args[1].(Number)
}

func gt(args ...Expression) Expression {
	return args[0].(Number) > args[1].(Number)
}

func equals(args ...Expression) Expression {
	if len(args) < 2 {
		return true
	}
	return reflect.DeepEqual(args[0], args[1])
}

func and(args ...Expression) Expression {
	return args[0].(bool) && args[1].(bool)
}

func or(args ...Expression) Expression {
	return args[0].(bool) || args[1].(bool)
}

func not(args ...Expression) Expression {
	return !args[0].(bool)
}

func car(args ...Expression) Expression {
	return args[0].([]Expression)[0]
}

func cdr(args ...Expression) Expression {
	return args[0].([]Expression)[1:]
}

func cons(args ...Expression) Expression {
	first := args[0]

	switch rest := args[1].(type) {

	// if the second arg is a list, then append it to the end of a new list
	case []Expression:
		return append([]Expression{first}, rest...)

	// otherwise assume create a new list from the two
	default:
		return []Expression{first, rest}
	}
}

func list(args ...Expression) Expression {
	result := make([]Expression, 0)
	for _, arg := range args {
		result = append(result, arg)
	}
	return result
}

// global scope
func NewEnv() *Scope {
	env := Scope{
		map[Symbol]Expression{
			"+":    numberType(add),
			"-":    numberType(sub),
			"*":    numberType(mult),
			"/":    numberType(div),
			"%":    numberType(mod),
			"<":    numberType(lt),
			">":    numberType(gt),
			"=":    equals,
			"and":  boolType(and),
			"or":   boolType(or),
			"not":  boolType(not),
			"car":  car,
			"cdr":  cdr,
			"cons": cons,
			"list": list,
		},
		nil,
	}
	return &env
}
