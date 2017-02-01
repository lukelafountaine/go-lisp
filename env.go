package main

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
	env := Env{
		map[Symbol]Expression{
			"abs": abs,
			"+" : add,
			"-" : sub,
			"*" : mult,
			"/" : div,
			"%" : mod,
			"<" : lt,
			"<=" : lte,
			">" : gt,
			">=" :gte,
			"==": equals,
			"&&" : and,
			"||" : or,
			"!" : not,
			"max": max,
			"min": min,
		},
		nil,
	}
	return &env
}
