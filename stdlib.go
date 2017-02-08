package main

import (
	"fmt"
	"errors"
)

// type check for functions that expect arguments of type Number
func numberType(f func(...Expression) Expression) func(...Expression) (Expression, error) {

	return func(args...Expression) (Expression, error) {
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

	return func(args...Expression) (Expression, error) {
		for _, arg := range args {
			if _, ok := arg.(bool); !ok {
				return nil, errors.New(fmt.Sprintf("Type Error: Recieved %T, Expected Boolean", arg))
			}
		}
		return f(args...), nil
	}
}

// built in functions
func abs(args... Expression) Expression {
	num := args[0].(Number)

	if num < 0 {
		num *= -1
	}
	return num
}

func add(args... Expression) Expression {

	initial := args[0].(Number)
	for _, num := range args[1:] {
		initial += num.(Number)
	}

	return initial
}

func sub(args... Expression) Expression {

	initial := args[0].(Number)
	for _, num := range args[1:] {
		initial -= num.(Number)
	}

	return initial
}

func mult(args... Expression) Expression {

	initial := args[0].(Number)
	for _, num := range args[1:] {
		initial *= num.(Number)
	}

	return initial
}

func div(args... Expression) Expression {

	initial := args[0].(Number)
	for _, num := range args[1:] {
		initial /= num.(Number)
	}

	return initial
}

func mod(args... Expression) Expression {
	return Number(int(args[0].(Number) / 1) % int(args[1].(Number) / 1))
}

func lt(args... Expression) Expression {
	return args[0].(Number) < args[1].(Number)
}

func lte(args... Expression) Expression {
	return args[0].(Number) <= args[1].(Number)
}

func gt(args... Expression) Expression {
	return args[0].(Number) > args[1].(Number)
}

func gte(args... Expression) Expression {
	return args[0].(Number) >= args[1].(Number)
}

func equals(args... Expression) Expression {
	return args[0].(Number) == args[1].(Number)
}

func and(args... Expression) Expression {
	return args[0].(bool) && args[1].(bool)
}

func or(args... Expression) Expression {
	return args[0].(bool) || args[1].(bool)
}

func not(args... Expression) Expression {
	return !args[0].(bool)
}

func max(args...Expression) Expression {
	biggest := args[0].(Number)

	for _, num := range args {
		if num.(Number) > biggest {
			biggest = num.(Number)
		}
	}
	return biggest
}

func min(args...Expression) Expression {
	smallest := args[0].(Number)

	for _, num := range args {
		if num.(Number) < smallest {
			smallest = num.(Number)
		}
	}
	return smallest
}

func car(args...Expression) Expression {
	return args[0].([]Expression)[0]
}

func cdr(args...Expression) Expression {
	return args[0].([]Expression)[1:]
}

func cons(args...Expression) Expression {
	first := args[0]

	switch rest := args[1].(type) {
	case []Expression:
		return append([]Expression{first}, rest...)
	default:
		return []Expression{first, rest}
	}
	return args[0].([]Expression)[1:]
}

func list(args...Expression) Expression {
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
			"abs": numberType(abs),
			"+" : numberType(add),
			"-" : numberType(sub),
			"*" : numberType(mult),
			"/" : numberType(div),
			"%" : numberType(mod),
			"max": numberType(max),
			"min": numberType(min),
			"<" : numberType(lt),
			"<=" : numberType(lte),
			">" : numberType(gt),
			">=" : numberType(gte),
			"=": numberType(equals),
			"&&" : boolType(and),
			"||" : boolType(or),
			"!" : boolType(not),
			"car": car,
			"cdr" : cdr,
			"cons" : cons,
			"list" : list,
		},
		nil,
	}
	return &env
}
