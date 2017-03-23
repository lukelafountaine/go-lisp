package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lukelafountaine/go-lisp/types"
	"github.com/lukelafountaine/go-lisp/parse"
	"github.com/lukelafountaine/go-lisp/eval"
)

func Run(program string, scope *types.Scope) {

	// read
	exp, err := parse.Parse(program)

	if err != nil {
		fmt.Println(err)
		return
	}

	// evaluate
	result, err := eval.Eval(exp, scope)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print
	if result != nil {
		fmt.Println(result)
	}
}

func Repl(scope *types.Scope) {
	scanner := bufio.NewScanner(os.Stdin)
	for fmt.Print("> "); scanner.Scan(); fmt.Print("> ") {
		Run(scanner.Text(), scope)
	}
}

func main() {

	scope := types.NewEnv()

	// load the standard library
	program, _ := ioutil.ReadFile("stdlib.lisp")
	Run(string(program), scope)

	// evaluate any files provided
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
