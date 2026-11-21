package compiler

import (
	"fmt"

	"donut/lang/parser"
)

type Code struct {
	codeBlock *parser.BlockStatement
	cp        *ConstantPool

	stackSize uint16
	maxLocals uint16
}

func NewCode(context *Context, body *parser.BlockStatement) *Code {
	return &Code{
		codeBlock: body,
		cp:        context.cp,
		stackSize: 255, // todo: unhardcode this
		maxLocals: 200, // todo: unhardcode this
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
		bc.writeU1(op_return)
	case *parser.ExpressionStatement:
		bc.dump(c.compileExpression(stmt.Expression))
	case *parser.WhileStatement:
		bc.dump(c.compileWhileStmt(*stmt))
	case *parser.IfStatement:
		bc.dump(c.compileIfStmt(*stmt))
	case *parser.BlockStatement:
		for _, s := range stmt.Body {
			bc.dump(c.compileStatement(s))
		}
	case *parser.PrintStatement: // todo: drop this once function calls are implemented
		bc.dump(c.compileExpression(stmt.Expression))
		bc.opInvokeStatic(c.cp.AddMethodRef("Console", "log", stmt.ResolvedArgs))
	default:
		panic(fmt.Sprintf("Unhandled statement type: %T", stmt))
	}
	return
}

func (c *Code) compileIfStmt(stmt parser.IfStatement) (bc ByteCode) {
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

func (c *Code) compileWhileStmt(stmt parser.WhileStatement) (bc ByteCode) {
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
		var value int32 = 0
		if expr.Value {
			value = 1
		}
		intIndex := c.cp.AddInteger(value)
		bc.opLoadConstant(intIndex)
	case *parser.IntegerLiteral:
		number := int32(expr.Value)
		if number >= -128 && number < 128 {
			bc.opBipush(byte(int8(number)))
		} else {
			intIndex := c.cp.AddInteger(int32(expr.Value))
			bc.opLoadConstant(intIndex)
		}
	case *parser.Identifier:
		switch expr.Variable.VarType {
		case parser.TypeInt, parser.TypeBoolean:
			bc.opIload(uint16(expr.Variable.Index))
		case parser.TypeFloat:
			bc.opFload(uint16(expr.Variable.Index))
		default:
			bc.opAload(uint16(expr.Variable.Index))
		}
	default:
		panic(fmt.Sprintf("Unhandled expression type: %T", expr))
	}
	return
}

func (c *Code) compileAssignmentExpression(expr parser.AssignmentExpression) (bc ByteCode) {
	bc.dump(c.compileExpression(expr.Right))
	switch expr := expr.Left.(type) {
	case *parser.Identifier:
		switch expr.Variable.VarType {
		case parser.TypeInt, parser.TypeBoolean:
			bc.opIstore(uint16(expr.Variable.Index))
		case parser.TypeFloat:
			bc.opFstore(uint16(expr.Variable.Index))
		default:
			bc.opAstore(uint16(expr.Variable.Index))
		}
	default:
		panic(fmt.Sprintf("Unhandled Assignment type: %T", expr))
	}
	return
}

func (c *Code) compileUnaryExpression(expr parser.UnaryExpression) (bc ByteCode) {
	bc = c.compileExpression(expr.Value)
	bc.dump(c.compileExpressionNode(expr.TypeInfo))
	return
}

func (c *Code) compileBinaryExpression(expr parser.BinaryExpression) (bc ByteCode) {
	bc = c.compileExpression(expr.Left)
	if expr.TypeInfo.LeftTo != parser.Bop_none {
		bc.dump(c.compileExpressionNode(expr.TypeInfo.LeftTo))
	}
	bc.dump(c.compileExpression(expr.Right))
	if expr.TypeInfo.RightTo != parser.Bop_none {
		bc.dump(c.compileExpressionNode(expr.TypeInfo.RightTo))
	}
	bc.dump(c.compileExpressionNode(expr.TypeInfo.OperationType))
	return bc
}

func (c *Code) compileExpressionNode(opType parser.ExpOpTypeId) (bc ByteCode) {
	switch opType {
	case parser.Bop_iadd:
		bc.writeU1(op_iadd)
	case parser.Bop_isub:
		bc.writeU1(op_isub)
	case parser.Bop_imul:
		bc.writeU1(op_imul)
	case parser.Bop_idiv:
		bc.writeU1(op_idiv)
	case parser.Bop_irem:
		bc.writeU1(op_irem)
	case parser.Bop_fadd:
		bc.writeU1(op_fadd)
	case parser.Bop_fsub:
		bc.writeU1(op_fsub)
	case parser.Bop_fmul:
		bc.writeU1(op_fmul)
	case parser.Bop_fdiv:
		bc.writeU1(op_fdiv)
	case parser.Bop_frem:
		bc.writeU1(op_frem)
	case parser.Bop_i2f:
		bc.writeU1(op_i2f)
	case parser.Bop_iand:
		bc.writeU1(op_iand)
	case parser.Bop_ior:
		bc.writeU1(op_ior)
	case parser.Uop_ineg:
		bc.writeU1(op_ineg)
	case parser.Uop_fneg:
		bc.writeU1(op_fneg)
	case parser.Uop_not:
		bc.writeU1(op_iconst_1)
		bc.writeU1(op_ixor)
	case parser.Bop_icmp_lt:
		c.compileComparison(&bc, op_if_icmplt, false)
	case parser.Bop_icmp_gt:
		c.compileComparison(&bc, op_if_icmpgt, false)
	case parser.Bop_icmp_le:
		c.compileComparison(&bc, op_if_icmple, false)
	case parser.Bop_icmp_ge:
		c.compileComparison(&bc, op_if_icmpge, false)
	case parser.Bop_icmp_eq:
		c.compileComparison(&bc, op_if_icmpeq, false)
	case parser.Bop_icmp_ne:
		c.compileComparison(&bc, op_if_icmpne, false)
	case parser.Bop_fcmp_lt:
		c.compileComparison(&bc, op_iflt, true)
	case parser.Bop_fcmp_gt:
		c.compileComparison(&bc, op_ifgt, true)
	case parser.Bop_fcmp_le:
		c.compileComparison(&bc, op_ifle, true)
	case parser.Bop_fcmp_ge:
		c.compileComparison(&bc, op_ifge, true)
	case parser.Bop_fcmp_eq:
		c.compileComparison(&bc, op_ifeq, true)
	case parser.Bop_fcmp_ne:
		c.compileComparison(&bc, op_ifne, true)
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
		bc.writeU2(op_fcmpl)
	}
	trueJump := bc.writeJumpPlaceholder(cmpOpcode)
	bc.writeU1(op_iconst_0)
	endJump := bc.writeJumpPlaceholder(op_goto)
	bc.patchJump(trueJump)
	bc.writeU1(op_iconst_1)
	bc.patchJump(endJump)
}
