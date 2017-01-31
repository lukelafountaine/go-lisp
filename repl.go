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

		exp, err := Parse(program)
		if err != nil {
			fmt.Println(err)
			continue
		}

		result, err := Eval(exp, scope)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if result != nil {
			fmt.Println(result)
		} else {
			fmt.Println("ok.")
		}
	}
}

func main() {
	Repl()
}
