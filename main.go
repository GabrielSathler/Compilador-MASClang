package main

import (
	"fmt"
	"os"

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

	p.ParseProgram()
	fmt.Println("parsing finalizado")
}
