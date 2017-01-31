package main

type symbols map[Symbol]Expression
type Env struct {
	symbols
	outer *Env
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

func NewEnv() *Env {

	env := Env{
		symbols{
			"+" : add,
			"-" : sub,
			"*" : mult,
			"/" : div,
			"<" : lt,
			"<=" : lte,
			">" : gt,
			">=" :gte,
			"==": equals,
			"max": max,
			"min": min,
		},
		nil,
	}

	return &env
}
