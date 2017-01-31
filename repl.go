package main

import (
	"bufio"
	"os"
	"fmt"
)

func Repl() {

	scope := NewEnv()
	for {
		reader := bufio.NewReader(os.Stdin)
		program, _ := reader.ReadString('\n')

		exp := Parse(program)
		result := Eval(exp, scope)

		fmt.Println(result)

	}
}

func main() {
	Repl()
}