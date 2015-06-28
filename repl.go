package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lukechampine/minima/lang"
)

func main() {
	fmt.Println("Minima REPL v0.1")
	fmt.Print("λ> ")
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	sexp, err := lang.ReadString(lang.Desugar(line))
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	sexp, err = lang.Eval(sexp)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(sexp)
}
