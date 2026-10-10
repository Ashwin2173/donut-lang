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
	for _, item := range utils.Tokenize(string(file)) {
		fmt.Println(item)
	}
}

func printUsage(shouldExit bool) {
	fmt.Println("USAGE:")
	fmt.Println("    donut *.dt [SUBCOMMANDS]")
	if shouldExit {
		os.Exit(1)
	}
}
