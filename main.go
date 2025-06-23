package main

import (
	"fmt"
	"os"

	"github.com/GabrielSathler/Compilador-MASClang/semantic_analyzer"
	"github.com/GabrielSathler/Compilador-MASClang/syntactic_analyzer"
)

func main() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	p := syntactic_analyzer.NewParser(file)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("error parsing: %v\n", r)
		}
	}()

	program := p.ParseProgram()

	analyzer := semantic_analyzer.NewSemanticAnalyzer()
	analyzer.Analyze(program)

	if len(analyzer.Errors()) > 0 {
		for _, err := range analyzer.Errors() {
			fmt.Println(err)
		}
	}
}
