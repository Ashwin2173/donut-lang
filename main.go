package main

import (
	"fmt"
	"os"

	"donut/lang"
	"donut/lang/compiler"
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
	compiler, err := compiler.New(args[1], string(file))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// test write this will be updated.
	err = os.WriteFile(args[1][:len(args[1])-3]+lang.TargetExtension, compiler.Compile(), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printUsage(shouldExit bool) {
	fmt.Println("USAGE:")
	fmt.Println("    donut *.dt [SUBCOMMANDS]")
	if shouldExit {
		os.Exit(1)
	}
}
