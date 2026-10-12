package utils

import "fmt"

type CompileError struct {
	Message  string
	FileName string
	Line     uint
	Col      uint
}

func (e CompileError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.FileName, e.Line, 0, e.Message)
}
