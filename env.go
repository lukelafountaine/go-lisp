package main

import (
	"math"
)

type Symbol struct {
	nodeType string
	val      interface{}
}

type Env map[string]interface{}

func NewEnv() *Env {
	env := Env {
		"x" : 5,
		"abs" : math.Abs,
		"+" : func (a, b int) int {return a + b},
		"-" : func (a, b int) int {return a - b},
		"/" : func (a, b int) int {return a / b},
		"*" : func (a, b int) int {return a * b},
	}

	return &env
}
