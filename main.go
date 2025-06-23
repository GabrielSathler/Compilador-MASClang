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
			fmt.Printf("Error parsing: %v\n", r)
		}
	}()

	program := p.ParseProgram()
	fmt.Printf("Declarações encontradas no programa: %d\n", len(program.Decls))

	analyzer := semantic_analyzer.NewSemanticAnalyzer()
	analyzer.Analyze(program)

	/*errs := analyzer.Errors()
	if len(errs) > 0 {
		fmt.Println("Erros semânticos encontrados:")
		for _, e := range errs {
			fmt.Println(" -", e)
		}
	} else {
		fmt.Println("Análise semântica concluída sem erros!")
	}*/
}
