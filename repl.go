package main

import (
	"bufio"
	"os"
	"fmt"
	"io/ioutil"
)

func Run(program string, scope *Env) {

	// read
	exp, err := Parse(program)

	if err != nil {
		fmt.Println(err)
		return
	}

	// evaluate
	result, err := Eval(exp, scope)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print
	if result != nil {
		fmt.Println(result)
	} else {
		fmt.Println("ok.")
	}
}

func Repl(scope *Env) {

	scanner := bufio.NewScanner(os.Stdin)
	var program string

	for fmt.Print("> "); scanner.Scan(); fmt.Print("> ") {

		program = scanner.Text()
		if program == "exit" {
			fmt.Println("bye!")
			os.Exit(0)
		}

		Run(program, scope)
	}
}

func main() {

	scope := NewEnv()
	for _, file := range os.Args[1:] {
		program, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		Run(string(program), scope)
	}

	Repl(scope)
}
