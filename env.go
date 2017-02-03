package main

import (
	"fmt"
	"errors"
)

// types
type Expression interface{}
type Symbol string
type Number float64
type Function struct {
	params, body Expression
	env          *Env
}
type Env struct {
	symbols map[Symbol]Expression
	outer   *Env
}

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

// global scope
func NewEnv() *Env {
	env := Env {
		map[Symbol]Expression {
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
		},
		nil,
	}
	return &env
}
