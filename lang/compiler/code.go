package compiler

import (
	"fmt"
	"maps"
	"slices"

	"donut/lang/parser"
)

var opcodeForOpr = map[parser.ExpOpTypeId]byte{
	parser.Bop_iadd: op_iadd,
	parser.Bop_isub: op_isub,
	parser.Bop_imul: op_imul,
	parser.Bop_idiv: op_idiv,
	parser.Bop_irem: op_irem,
	parser.Bop_fadd: op_fadd,
	parser.Bop_fsub: op_fsub,
	parser.Bop_fmul: op_fmul,
	parser.Bop_fdiv: op_fdiv,
	parser.Bop_frem: op_frem,
	parser.Bop_i2f:  op_i2f,
	parser.Bop_iand: op_iand,
	parser.Bop_ior:  op_ior,
	parser.Uop_ineg: op_ineg,
	parser.Uop_fneg: op_fneg,
}

type cmpOp struct {
	opcode  byte
	isFloat bool
}

var comparisonOpcode = map[parser.ExpOpTypeId]cmpOp{
	parser.Bop_fcmp_lt: {op_iflt, true},
	parser.Bop_fcmp_gt: {op_ifgt, true},
	parser.Bop_fcmp_le: {op_ifle, true},
	parser.Bop_fcmp_ge: {op_ifge, true},
	parser.Bop_fcmp_eq: {op_ifeq, true},
	parser.Bop_fcmp_ne: {op_ifne, true},
	parser.Bop_icmp_lt: {op_if_icmplt, false},
	parser.Bop_icmp_gt: {op_if_icmpgt, false},
	parser.Bop_icmp_le: {op_if_icmple, false},
	parser.Bop_icmp_ge: {op_if_icmpge, false},
	parser.Bop_icmp_eq: {op_if_icmpeq, false},
	parser.Bop_icmp_ne: {op_if_icmpne, false},
}

type Code struct {
	codeBlock *parser.BlockStatement
	cp        *ConstantPool

	stackSize uint16
	maxLocals uint16

	LocalsStack []map[string]uint16
	LocalsPool  []uint16
}

func NewCode(context *Context, body *parser.BlockStatement) *Code {
	return &Code{
		codeBlock: body,
		cp:        context.cp,
		stackSize: 255, // todo: unhardcode this
		maxLocals: 0,
	}
}

func (c Code) GetByteCode(context *Context) (bc ByteCode) {
	stmtBc := c.compileStatement(c.codeBlock)
	bc.writeU2(c.stackSize)
	bc.writeU2(c.maxLocals)
	bc.writeU4(uint32(stmtBc.length()))
	bc.dump(stmtBc)
	bc.writeU2(0) // exception table
	bc.writeU2(0) // attribute count
	return
}

func (c *Code) compileStatement(stmt parser.Statement) (bc ByteCode) {
	switch stmt := stmt.(type) {
	case *parser.ReturnStatement:
		bc.dump(c.compileReturnStmt(stmt))
	case *parser.ExpressionStatement:
		bc.dump(c.compileExpression(stmt.Expression))
	case *parser.WhileStatement:
		bc.dump(c.compileWhileStmt(stmt))
	case *parser.IfStatement:
		bc.dump(c.compileIfStmt(stmt))
	case *parser.BlockStatement:
		c.pushLocalsStack()
		for _, s := range stmt.Body {
			bc.dump(c.compileStatement(s))
		}
		c.popLocalsStack()
	case *parser.PrintStatement: // todo: drop this once function calls are implemented
		bc.dump(c.compileExpression(stmt.Expression))
		bc.opInvokeStatic(c.cp.AddMethodRef("Console", "log", stmt.ResolvedArgs))
	default:
		panic(fmt.Sprintf("Unhandled statement type: %T", stmt))
	}
	return
}

func (c *Code) compileReturnStmt(stmt *parser.ReturnStatement) (bc ByteCode) {
	bc.dump(c.compileExpression(stmt.Expression))
	switch stmt.ReturnType {
	case parser.TypeBoolean, parser.TypeInt:
		bc.writeU1(op_ireturn)
	case parser.TypeFloat:
		bc.writeU1(op_freturn)
	case parser.TypeVoid:
		bc.writeU1(op_return)
	default:
		bc.writeU1(op_areturn)
	}
	return
}

func (c *Code) compileIfStmt(stmt *parser.IfStatement) (bc ByteCode) {
	bc.dump(c.compileExpression(stmt.Test))
	consequentPointer := bc.writeJumpPlaceholder(op_ifeq)
	bc.dump(c.compileStatement(stmt.Consequent))
	if stmt.Alternate != nil {
		alternatePointer := bc.writeJumpPlaceholder(op_goto)
		bc.patchJump(consequentPointer)
		bc.dump(c.compileStatement(stmt.Alternate))
		bc.patchJump(alternatePointer)
	} else {
		bc.patchJump(consequentPointer)
	}
	return
}

func (c *Code) compileWhileStmt(stmt *parser.WhileStatement) (bc ByteCode) {
	start := len(bc.byteCode)
	bc.dump(c.compileExpression(stmt.Test))
	exitJump := bc.writeJumpPlaceholder(op_ifeq)
	bc.dump(c.compileStatement(stmt.Body))
	loopJump := bc.writeJumpPlaceholder(op_goto)
	bc.patchJumpTo(loopJump, start)
	bc.patchJump(exitJump)
	return
}

func (c *Code) compileExpression(expr parser.Expression) (bc ByteCode) {
	switch expr := expr.(type) {
	case *parser.UnaryExpression:
		bc.dump(c.compileUnaryExpression(*expr))
	case *parser.BinaryExpression:
		bc.dump(c.compileBinaryExpression(*expr))
	case *parser.AssignmentExpression:
		bc.dump(c.compileAssignmentExpression(*expr))
	case *parser.StringLiteral:
		stringIndex := c.cp.AddString(expr.Value)
		bc.opLoadConstant(stringIndex)
	case *parser.FloatLiteral:
		floatIndex := c.cp.AddFloat(float32(expr.Value))
		bc.opLoadConstant(floatIndex)
	case *parser.BooleanLiteral:
		if expr.Value {
			bc.writeU1(op_iconst_1)
		} else {
			bc.writeU1(op_iconst_0)
		}
	case *parser.IntegerLiteral:
		number := int32(expr.Value)
		if number >= -128 && number < 128 {
			bc.opBipush(byte(number))
		} else {
			intIndex := c.cp.AddInteger(int32(expr.Value))
			bc.opLoadConstant(intIndex)
		}
	case *parser.Identifier:
		index := c.getLocal(expr.Name)
		switch expr.Variable.VarType {
		case parser.TypeInt, parser.TypeBoolean:
			bc.opLocal(op_iload_0, op_iload, index)
		case parser.TypeFloat:
			bc.opLocal(op_fload_0, op_fload, index)
		default:
			bc.opLocal(op_aload_0, op_aload, index)
		}
	default:
		panic(fmt.Sprintf("Unhandled expression type: %T", expr))
	}
	return
}

func (c *Code) compileAssignmentExpression(expr parser.AssignmentExpression) (bc ByteCode) {
	bc.dump(c.compileExpression(expr.Right))
	switch left := expr.Left.(type) {
	case *parser.Identifier:
		var index uint16
		if expr.IsNew {
			index = c.newLocal(left.Name)
		} else {
			index = c.getLocal(left.Name)
		}
		switch left.Variable.VarType {
		case parser.TypeInt, parser.TypeBoolean:
			bc.opLocal(op_istore_0, op_istore, index)
		case parser.TypeFloat:
			bc.opLocal(op_fstore_0, op_fstore, index)
		default:
			bc.opLocal(op_astore_0, op_astore, index)
		}
	default:
		panic(fmt.Sprintf("Unhandled Assignment type: %T", expr))
	}
	return
}

func (c *Code) compileUnaryExpression(expr parser.UnaryExpression) (bc ByteCode) {
	bc = c.compileExpression(expr.Value)
	bc.dump(c.compileExpressionOperation(expr.TypeInfo))
	return
}

func (c *Code) compileBinaryExpression(expr parser.BinaryExpression) (bc ByteCode) {
	bc = c.compileExpression(expr.Left)
	if expr.TypeInfo.LeftTo != parser.Bop_none {
		bc.dump(c.compileExpressionOperation(expr.TypeInfo.LeftTo))
	}
	bc.dump(c.compileExpression(expr.Right))
	if expr.TypeInfo.RightTo != parser.Bop_none {
		bc.dump(c.compileExpressionOperation(expr.TypeInfo.RightTo))
	}
	bc.dump(c.compileExpressionOperation(expr.TypeInfo.OperationType))
	return bc
}

func (c *Code) compileExpressionOperation(opType parser.ExpOpTypeId) (bc ByteCode) {
	if opcode, exists := opcodeForOpr[opType]; exists {
		bc.writeU1(opcode)
		return
	}
	if cmp, exists := comparisonOpcode[opType]; exists {
		c.compileComparison(&bc, cmp.opcode, cmp.isFloat)
		return
	}
	switch opType {
	case parser.Uop_not:
		bc.writeU1(op_iconst_1)
		bc.writeU1(op_ixor)
	case parser.Bop_s_concat:
		bc.opInvokeVirtual(c.cp.AddMethodRef("java/lang/String", "concat", "(Ljava/lang/String;)Ljava/lang/String;"))
	case parser.Bop_s_eq:
		bc.opInvokeVirtual(c.cp.AddMethodRef("java/lang/String", "equals", "(Ljava/lang/Object;)Z"))
	case parser.Bop_s_ne:
		bc.opInvokeVirtual(c.cp.AddMethodRef("java/lang/String", "equals", "(Ljava/lang/Object;)Z"))
		bc.writeU1(op_iconst_1)
		bc.writeU1(op_ixor)
	default:
		panic(fmt.Sprintf("Unhandled opType: %d", opType))
	}
	return
}

func (c *Code) compileComparison(bc *ByteCode, cmpOpcode byte, isFloat bool) {
	if isFloat {
		bc.writeU1(op_fcmpl)
	}
	trueJump := bc.writeJumpPlaceholder(cmpOpcode)
	bc.writeU1(op_iconst_0)
	endJump := bc.writeJumpPlaceholder(op_goto)
	bc.patchJump(trueJump)
	bc.writeU1(op_iconst_1)
	bc.patchJump(endJump)
}

func (c *Code) newLocal(name string) (index uint16) {
	LocalsPoolLength := len(c.LocalsPool)
	frameLength := len(c.LocalsStack)
	if LocalsPoolLength != 0 {
		index = c.LocalsPool[LocalsPoolLength-1]
		c.LocalsPool = c.LocalsPool[:LocalsPoolLength-1]
		c.LocalsStack[frameLength-1][name] = index
		return
	}
	c.LocalsStack[frameLength-1][name] = c.maxLocals
	index = c.maxLocals
	c.maxLocals++
	return
}

func (c *Code) getLocal(name string) uint16 {
	for _, frame := range slices.Backward(c.LocalsStack) {
		index, exists := frame[name]
		if exists {
			return index
		}
	}
	panic("RUC at compiler/code/getLocal: Analyzer failed to catch this")
}

func (c *Code) pushLocalsStack() {
	c.LocalsStack = append(c.LocalsStack, make(map[string]uint16))
}

func (c *Code) popLocalsStack() {
	frame := c.LocalsStack[len(c.LocalsStack)-1]
	c.LocalsPool = append(c.LocalsPool, slices.Collect(maps.Values(frame))...)
	c.LocalsStack = c.LocalsStack[:len(c.LocalsStack)-1]
}
