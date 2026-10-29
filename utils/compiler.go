package utils

import (
	"fmt"
	"path/filepath"
)

const MAGIC uint32 = 0xCAFEBABE
const MAJOR uint16 = 0
const MINOR uint16 = 50

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
	// fmt.Println(c.cp.constantPool)
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
	codeAttr := c.compileCodeAttr(function)
	fbc.writeU4(uint32(codeAttr.length()))
	fbc.dump(codeAttr)

	return
}

func (c *Compiler) compileCodeAttr(function FunctionStatement) (bc ByteCode) {
	body := function.body
	bc.writeU2(256)                        // max stack size; todo: unhardcode this.
	bc.writeU2(uint16(function.maxLocals)) //max locals; todo: unhardcode this.

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
			bc.opInvokeStatic(c.cp.AddMethodRef("Console", "log", stmt.resolvedArgs))
		case *ExpressionStatement:
			bc.dump(c.compileExpression(stmt.expression))
		default:
			panic(fmt.Sprintf("Unhandled statement type: %T", stmt))
		}
	}
	return
}

func (c *Compiler) compileExpression(expr Expression) (bc ByteCode) {
	switch expr := expr.(type) {
	case *BinaryExpression:
		bc.dump(c.compileBinaryExpression(*expr))
	case *AssignmentExpression:
		bc.dump(c.compileAssignmentExpression(*expr))
	case *StringLiteral:
		stringIndex := c.cp.AddString(expr.value)
		bc.opLoadConstant(stringIndex)
	case *FloatLiteral:
		floatIndex := c.cp.AddFloat(float32(expr.value))
		bc.opLoadConstant(floatIndex)
	case *BooleanLiteral:
		var value int32 = 0
		if expr.value {
			value = 1
		}
		intIndex := c.cp.AddInteger(value)
		bc.opLoadConstant(intIndex)
	case *IntegerLiteral:
		number := int32(expr.value)
		if number >= -128 && number < 128 {
			bc.opBipush(byte(int8(number)))
		} else {
			intIndex := c.cp.AddInteger(int32(expr.value))
			bc.opLoadConstant(intIndex)
		}
	case *Identifier:
		switch expr.variable.varType {
		case T_INT, T_BOOLEAN:
			bc.opIload(uint16(expr.variable.index))
		case T_FLOAT:
			bc.opFload(uint16(expr.variable.index))
		default:
			bc.opAload(uint16(expr.variable.index))
		}
	default:
		panic(fmt.Sprintf("Unhandled expression type: %T", expr))
	}
	return
}

func (c *Compiler) compileAssignmentExpression(expr AssignmentExpression) (bc ByteCode) {
	bc.dump(c.compileExpression(expr.right))
	switch expr := expr.left.(type) {
	case *Identifier:
		switch expr.variable.varType {
		case T_INT, T_BOOLEAN:
			bc.opIstore(uint16(expr.variable.index))
		case T_FLOAT:
			bc.opFstore(uint16(expr.variable.index))
		default:
			bc.opAstore(uint16(expr.variable.index))
		}
	default:
		panic(fmt.Sprintf("Unhandled Assignment type: %T", expr))
	}
	return
}

func (c *Compiler) compileBinaryExpression(expr BinaryExpression) (bc ByteCode) {
	bc = c.compileExpression(expr.left)
	if expr.typeInfo.leftTo != BOP_NONE {
		bc.dump(c.compileExpressionNode(expr.typeInfo.leftTo))
	}
	bc.dump(c.compileExpression(expr.right))
	if expr.typeInfo.rightTo != BOP_NONE {
		bc.dump(c.compileExpressionNode(expr.typeInfo.rightTo))
	}
	bc.dump(c.compileExpressionNode(expr.typeInfo.operationType))
	return bc
}

func (c *Compiler) compileExpressionNode(opType BinOpTypeId) (bc ByteCode) {
	switch opType {
	case BOP_IADD:
		bc.opIAdd()
	case BOP_ISUB:
		bc.opISub()
	case BOP_IMUL:
		bc.opIMul()
	case BOP_IDIV:
		bc.opIDiv()
	case BOP_IREM:
		bc.opIRem()
	case BOP_FADD:
		bc.opFAdd()
	case BOP_FSUB:
		bc.opFSub()
	case BOP_FMUL:
		bc.opFMul()
	case BOP_FDIV:
		bc.opFDiv()
	case BOP_FREM:
		bc.opFRem()
	case BOP_I2F:
		bc.opI2F()
	case BOP_IAND:
		bc.opIAnd()
	case BOP_IOR:
		bc.opIOr()
	case BOP_ICMP_LT:
		c.compileComparison(&bc, OP_IF_ICMPLT, false)
	case BOP_ICMP_GT:
		c.compileComparison(&bc, OP_IF_ICMPGT, false)
	case BOP_ICMP_LE:
		c.compileComparison(&bc, OP_IF_ICMPLE, false)
	case BOP_ICMP_GE:
		c.compileComparison(&bc, OP_IF_ICMPGE, false)
	case BOP_ICMP_EQ:
		c.compileComparison(&bc, OP_IF_ICMPEQ, false)
	case BOP_ICMP_NE:
		c.compileComparison(&bc, OP_IF_ICMPNE, false)
	case BOP_FCMP_LT:
		c.compileComparison(&bc, OP_IFLT, true)
	case BOP_FCMP_GT:
		c.compileComparison(&bc, OP_IFGT, true)
	case BOP_FCMP_LE:
		c.compileComparison(&bc, OP_IFLE, true)
	case BOP_FCMP_GE:
		c.compileComparison(&bc, OP_IFGE, true)
	case BOP_FCMP_EQ:
		c.compileComparison(&bc, OP_IFEQ, true)
	case BOP_FCMP_NE:
		c.compileComparison(&bc, OP_IFNE, true)
	case BOP_S_CONCAT:
		bc.opInvokeVirtual(c.cp.AddMethodRef("java/lang/String", "concat", "(Ljava/lang/String;)Ljava/lang/String;"))
	case BOP_S_EQ:
		bc.opInvokeVirtual(c.cp.AddMethodRef("java/lang/String", "equals", "(Ljava/lang/Object;)Z"))
	case BOP_S_NE:
		bc.opInvokeVirtual(c.cp.AddMethodRef("java/lang/String", "equals", "(Ljava/lang/Object;)Z"))
		bc.opIconst1()
		bc.opIXor()
	default:
		panic(fmt.Sprintf("Unhandled opType: %d", opType))
	}
	return
}

func (c *Compiler) compileComparison(bc *ByteCode, cmpOpcode byte, isFloat bool) {
	if isFloat {
		bc.opFcmpl()
	}
	trueJump := bc.writeJumpPlaceHolder(cmpOpcode)
	bc.opIconst0()
	endJump := bc.writeJumpPlaceHolder(OP_GOTO)
	bc.patchJump(trueJump)
	bc.opIconst1()
	bc.patchJump(endJump)
}
