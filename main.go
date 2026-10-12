package main

import (
	"os"
	"fmt"

	"donut/utils"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		printUsage(true)
	}
	file, err := os.ReadFile(args[1])
	if err != nil {
		fmt.Println("ERROR:", err)
		printUsage(true)
	}
	parser := utils.NewParser(args[1], string(file))
	// debugging code
	program, err := parser.ParseProgram()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(program)
}

func printUsage(shouldExit bool) {
	fmt.Println("USAGE:")
	fmt.Println("    donut *.dt [SUBCOMMANDS]")
	if shouldExit {
		os.Exit(1)
	}
}
