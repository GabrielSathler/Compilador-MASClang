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

	parser := syntactic_analyzer.NewParser(file)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error parsing: %v\n", r)
		}
	}()

	program := parser.ParseProgram()

	analyzer := semantic_analyzer.NewSemanticAnalyzer()
	analyzer.Analyze(program)

	if errs := analyzer.Errors; len(errs) > 0 {
		fmt.Println("Semantic errors:")

		for _, err := range errs {
			fmt.Println(" -", err)
		}

		return
	}
}
