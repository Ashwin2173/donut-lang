package compiler

import (
	"path/filepath"

	"donut/lang/parser"
)

type Context struct {
	FileName string

	className string
	program   *parser.Program
	cp        *ConstantPool
}

func NewContext(fileName string) *Context {
	fileName = filepath.Base(fileName)
	extention := filepath.Ext(fileName)
	return &Context{
		FileName:  fileName,
		className: fileName[:len(fileName)-len(extention)],
		cp:        NewConstantPool(),
	}
}

func (c *Context) parseSource(source string) (err error) {
	c.program, err = parser.New(c.FileName, source).ParseProgram()
	if err != nil {
		return err
	}
	return nil
}
