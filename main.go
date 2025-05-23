package main

import (
	"fmt"
	"os"

	"github.com/GabrielSathler/Compilador-MASClang/lexical_analyzer"
	"github.com/GabrielSathler/Compilador-MASClang/tokens"
)

func main() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}

	lexer := lexical_analyzer.NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == tokens.EOF {
			break
		}

		fmt.Printf("%d:%d\t%s\t%s\n", pos.Line, pos.Column, tok, lit)
	}
}
