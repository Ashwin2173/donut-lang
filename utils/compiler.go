package utils

import (
	"fmt"
	"path/filepath"
)

const MAGIC uint32 = 0xCAFEBABE
const MAJOR uint16 = 0
const MINOR uint16 = 62

type Compiler struct {
	fileName  string
	ClassName string
	parser    Parser
	bc        ByteCode
	cp        ConstantPool
}

func NewCompiler(fileName, source string) *Compiler {
	bc := NewByteCode()
	bc.writeU4(MAGIC)
	bc.writeU2(MAJOR)
	bc.writeU2(MINOR)
	return &Compiler{
		fileName: fileName,
		parser:   *NewParser(fileName, source),
		cp:       *NewConstantPool(),
		bc:       bc,
	}
}

func (c *Compiler) GetByteCode() []byte {
	return c.bc.byteCode
}

func (c *Compiler) Compile() error {
	program, err := c.parser.ParseProgram()
	if err != nil {
		return err
	}
	typeChecker := NewTypeChecker(c.parser.fileName, *program)
	err = typeChecker.resolve()
	if err != nil {
		return err
	}
	programByteCode := c.compileProgram(typeChecker.program)
	c.bc.writeU2(c.cp.length() + 1)
	c.bc.dump(c.cp.bc)
	fmt.Println(c.cp.constantPool)
	c.bc.dump(programByteCode)
	return err
}

func (c *Compiler) compileProgram(program Program) (pbc ByteCode) {
	fileName := filepath.Base(c.fileName)
	ext := filepath.Ext(fileName)
	c.ClassName = fileName[:len(fileName)-len(ext)]

	// class infos
	pbc.writeU2(0x1)                               // making source class public
	pbc.writeU2(c.cp.AddClass(c.ClassName))        // this class name
	pbc.writeU2(c.cp.AddClass("java/lang/Object")) // super class Name
	pbc.writeU2(0)                                 // interface count
	pbc.writeU2(0)                                 // fields count
	pbc.writeU2(uint16(len(program.functions)))

	// functions
	for _, function := range program.functions {
		pbc.dump(c.compileFunction(function))
	}

	// class attrs
	pbc.writeU2(0x1)
	pbc.writeU2(c.cp.AddUtf8("SourceFile"))
	pbc.writeU4(0x2)
	pbc.writeU2(c.cp.AddUtf8(fileName))
	return
}

func (c *Compiler) compileFunction(function FunctionStatement) (fbc ByteCode) {
	fbc.writeU2(0x9)                                    // access flag [public, static]; All functions are defaulted to this
	fbc.writeU2(c.cp.AddUtf8(function.name))            // function Name
	fbc.writeU2(c.cp.AddUtf8("([Ljava/lang/String;)V")) // function arg string; todo: unhardcode this
	fbc.writeU2(0x1)                                    // attr count; current only code exits

	// code attr
	fbc.writeU2(c.cp.AddUtf8("Code"))
	codeAttr := c.compileCodeAttr(function.body)
	fbc.writeU4(uint32(codeAttr.length()))
	fbc.dump(codeAttr)

	return
}

func (c *Compiler) compileCodeAttr(body BlockStatement) (bc ByteCode) {
	bc.writeU2(256) // max stack size; todo: unhardcode this.
	bc.writeU2(1)   //max locals; todo: unhardcode this.

	stmtByteCode := c.compileStatement(body.body)
	bc.writeU4(uint32(stmtByteCode.length()))
	bc.dump(stmtByteCode)

	bc.writeU2(0) // exception table
	bc.writeU2(0) // attribute count
	return
}

func (c *Compiler) compileStatement(stmts []Statement) (bc ByteCode) {
	for _, stmt := range stmts {
		switch stmt := stmt.(type) {
		case *ReturnStatement:
			bc.opReturn()
		case *PrintStatement:
			bc.dump(c.compileExpression(stmt.expression))
			bc.opInvokeStatic(c.cp.AddMethodRef("Console", "log", "(I)V")) // console.log(<string>) method; todo: unhardcode this
		case *ExpressionStatement:
			c.compileExpression(stmt.expression)
		default:
			panic(fmt.Sprintf("Unhandled statement type: %T", stmt))
		}
	}
	return
}

func (c *Compiler) compileExpression(expr Expression) (bc ByteCode) {
	switch expr := expr.(type) {
	case *StringLiteral:
		stringIndex := c.cp.AddString(expr.value)
		bc.opLoadConstant(stringIndex)
	case *IntegerLiteral:
		intIndex := c.cp.AddInteger(int32(expr.value))
		bc.opLoadConstant(intIndex)
	case *BooleanLiteral:
		var value int32 = 0
		if expr.value {
			value = 1
		}
		intIndex := c.cp.AddInteger(value)
		bc.opLoadConstant(intIndex)
	case *BinaryExpression:
		bc.dump(c.compileBinaryExpression(*expr))
	default:
		panic(fmt.Sprintf("Unhandled expression type: %T", expr))
	}
	return
}

func (c *Compiler) compileBinaryExpression(expr BinaryExpression) (bc ByteCode) {
	bc = c.compileExpression(expr.left)
	if expr.resolvedInfo.leftTo != nil {
		panic("binary leftTo is not supported yet")
	}
	bc.dump(c.compileExpression(expr.right))
	if expr.resolvedInfo.rightTo != nil {
		panic("binary rightTo is not supported yet")
	}
	exprPair := Pair[TypeId, TokenType]{expr.resolvedInfo.result, expr.binaryType}
	generator, exists := EXPR_OPERATION_TABLE[exprPair]
	if !exists {
		panic(fmt.Sprintln("Unhandled expr operation comination", exprPair))
	}
	generator(&bc)
	return bc
}
