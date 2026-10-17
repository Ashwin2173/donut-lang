package main

import (
	"fmt"
	"os"

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
	compiler := utils.NewCompiler(args[1], string(file))
	err = compiler.Compile()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile(args[1] + utils.TAGET_EXTENSION, compiler.GetByteCode(), 0644)
	if err !=  nil {
		fmt.Println(err)
	}
}

func printUsage(shouldExit bool) {
	fmt.Println("USAGE:")
	fmt.Println("    donut *.dt [SUBCOMMANDS]")
	if shouldExit {
		os.Exit(1)
	}
}
