package main

type symbols map[Symbol]Expression
type Env struct {
	symbols
	outer   *Env
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

func NewEnv() *Env {

	env := Env{
		symbols{
			"+" : add,
			"-" : sub,
			"*" : mult,
			"/" : div,
		},
		nil,
	}

	return &env
}
