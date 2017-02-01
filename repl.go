package main

import (
	"bufio"
	"os"
	"fmt"
	"io/ioutil"
)

func Repl() {

	scope := NewEnv()
	scanner := bufio.NewScanner(os.Stdin)

	for fmt.Print("> "); scanner.Scan(); fmt.Print("> "){


		exp, err := Parse(scanner.Text())
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

	if len(os.Args) == 1 {
		Repl()
	} else {
		scope := NewEnv()
		data, err := ioutil.ReadFile(os.Args[1])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		exp, err := Parse(string(data))
		if err != nil {
			fmt.Println(err)
		}

		result, err := Eval(exp, scope)
		if err != nil {
			fmt.Println(err)
		}

		if result != nil {
			fmt.Println(result)
		} else {
			fmt.Println("ok.")
		}
	}
}
