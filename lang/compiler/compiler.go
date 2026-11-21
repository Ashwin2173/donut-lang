package compiler

import (
	"donut/lang/parser"
)

const (
	magic = 0xCAFEBABE
	major = 0
	minor = 50
)

type Compiler struct {
	filePath string
	program  *parser.Program
}

func New(filePath, source string) (*Compiler, error) {
	p := parser.New(filePath, source)
	program, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}
	return &Compiler{
		filePath: filePath,
		program:  program,
	}, nil
}

func (c Compiler) Compile() []byte {
	context := NewContext(c.filePath)
	class := NewClass(context, Class{
		AccessFlag: AccPublic,
		ThisClass:  context.ClassName,
		SuperClass: "java/lang/Object",
		Interfaces: []string{},
		Fields:     []string{},
		Methods:    c.getMethods(context),
		Attributes: "",
	})
	return c.packBytes(context, class)
}

func (c Compiler) packBytes(context *Context, class *Class) []byte {
	classBc := class.GetByteCode(context)
	bc := NewByteCode()
	bc.writeU4(magic)
	bc.writeU2(major)
	bc.writeU2(minor)
	bc.writeU2(context.cp.length() + 1)
	bc.dump(context.cp.bc)
	bc.dump(classBc)
	return bc.byteCode
}

func (c Compiler) getMethods(context *Context) []Method {
	methods := make([]Method, len(c.program.Functions))
	for i, function := range c.program.Functions {
		code := NewCode(context, &function.Body)
		methods[i] = *NewMethod(context, Method{
			AccessFlag: AccPublic | AccStatic,
			Name:       function.Name,
			Descriptor: "([Ljava/lang/String;)V", // todo: unhardcode this
			Code:       *code,
			Attributes: []string{},
		})
	}
	return methods
}
