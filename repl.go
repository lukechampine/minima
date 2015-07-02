package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lukechampine/minima/lang"
)

func main() {
	fmt.Println("Minima REPL v0.2. ^D to quit.")
	r := bufio.NewReader(os.Stdin)

REPL:
	for {
		fmt.Print("λ> ")

		var expr string
		// read until parens match
		var balance int
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("\nLeaving Minima REPL.")
					return
				}
				fmt.Println("read error:", err)
				continue REPL
			}
			balance += strings.Count(line, "(") - strings.Count(line, ")")
			expr += line
			if balance == 0 {
				break
			}
		}
		sexp, err := lang.ReadString(expr)
		if err != nil {
			fmt.Println("parse error:", err)
			continue
		}

		sexp, err = lang.Eval(sexp)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		fmt.Println(sexp)
	}
}
