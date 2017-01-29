package main

import (
	"bufio"
	"os"
	"fmt"
)

func Repl() {
	for {
		reader := bufio.NewReader(os.Stdin)
		program, _ := reader.ReadString('\n')

		exp := Parse(program)
		scope := NewEnv()
		result := Eval(exp, scope)

		fmt.Println(result)

	}
}

func main() {
	Repl()
}