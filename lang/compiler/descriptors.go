package compiler

import "fmt"

func parserfnDescriptors(arguments, returnType string) string {
	// todo: make the arguments and slice of Arguments
	return arguments + parseType(returnType)
}

func parseType(rawType string) string {
	switch rawType {
	case "void":
		return "V"
	case "float":
		return "F"
	case "int":
		return "I"
	case "bool":
		return "Z"
	case "string":
		return "Ljava/lang/String;"
	default:
		return fmt.Sprintf("L%s;", rawType)
	}
}
