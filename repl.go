package main

import (
	"fmt"
	"os"

	"github.com/lukechampine/minima/lang"
)

func main() {
	fmt.Println("Minima REPL v0.1")
	fmt.Print("λ> ")
	sexp, err := lang.Read(os.Stdin)
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	fmt.Println(sexp)
}
