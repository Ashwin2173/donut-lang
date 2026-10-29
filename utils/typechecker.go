package utils

import (
	"fmt"
)

type TypeChecker struct {
	fileName      string
	program       Program
	functions     map[string][]string
	variableStack map[string]Variable // todo: convert this to a stack which will holder the scoped variables.
	typeRegister  map[string]TypeId
}

var binaryRules = map[BinaryTypeKey]BinaryTypeValue{
	{T_INT, PLUS, T_INT}:  {T_INT, BOP_IADD, BOP_NONE, BOP_NONE},
	{T_INT, MINUS, T_INT}: {T_INT, BOP_ISUB, BOP_NONE, BOP_NONE},
	{T_INT, STAR, T_INT}:  {T_INT, BOP_IMUL, BOP_NONE, BOP_NONE},
	{T_INT, SLASH, T_INT}: {T_INT, BOP_IDIV, BOP_NONE, BOP_NONE},
	{T_INT, MOD, T_INT}:   {T_INT, BOP_IREM, BOP_NONE, BOP_NONE},

	{T_INT, LESSER, T_INT}:         {T_BOOLEAN, BOP_ICMP_LT, BOP_NONE, BOP_NONE},
	{T_INT, GREATER, T_INT}:        {T_BOOLEAN, BOP_ICMP_GT, BOP_NONE, BOP_NONE},
	{T_INT, LESSER_EQUALS, T_INT}:  {T_BOOLEAN, BOP_ICMP_LE, BOP_NONE, BOP_NONE},
	{T_INT, GREATER_EQUALS, T_INT}: {T_BOOLEAN, BOP_ICMP_GE, BOP_NONE, BOP_NONE},
	{T_INT, EQUAL_EQUAL, T_INT}:    {T_BOOLEAN, BOP_ICMP_EQ, BOP_NONE, BOP_NONE},
	{T_INT, BANG_EQUAL, T_INT}:     {T_BOOLEAN, BOP_ICMP_NE, BOP_NONE, BOP_NONE},

	{T_INT, PLUS, T_FLOAT}:  {T_FLOAT, BOP_FADD, BOP_I2F, BOP_NONE},
	{T_INT, MINUS, T_FLOAT}: {T_FLOAT, BOP_FSUB, BOP_I2F, BOP_NONE},
	{T_INT, STAR, T_FLOAT}:  {T_FLOAT, BOP_FMUL, BOP_I2F, BOP_NONE},
	{T_INT, SLASH, T_FLOAT}: {T_FLOAT, BOP_FDIV, BOP_I2F, BOP_NONE},
	{T_INT, MOD, T_FLOAT}:   {T_FLOAT, BOP_FREM, BOP_I2F, BOP_NONE},

	{T_INT, LESSER, T_FLOAT}:         {T_BOOLEAN, BOP_FCMP_LT, BOP_I2F, BOP_NONE},
	{T_INT, GREATER, T_FLOAT}:        {T_BOOLEAN, BOP_FCMP_GT, BOP_I2F, BOP_NONE},
	{T_INT, LESSER_EQUALS, T_FLOAT}:  {T_BOOLEAN, BOP_FCMP_LE, BOP_I2F, BOP_NONE},
	{T_INT, GREATER_EQUALS, T_FLOAT}: {T_BOOLEAN, BOP_FCMP_GE, BOP_I2F, BOP_NONE},
	{T_INT, EQUAL_EQUAL, T_FLOAT}:    {T_BOOLEAN, BOP_FCMP_EQ, BOP_I2F, BOP_NONE},
	{T_INT, BANG_EQUAL, T_FLOAT}:     {T_BOOLEAN, BOP_FCMP_NE, BOP_I2F, BOP_NONE},

	{T_FLOAT, PLUS, T_FLOAT}:  {T_FLOAT, BOP_FADD, BOP_NONE, BOP_NONE},
	{T_FLOAT, MINUS, T_FLOAT}: {T_FLOAT, BOP_FSUB, BOP_NONE, BOP_NONE},
	{T_FLOAT, STAR, T_FLOAT}:  {T_FLOAT, BOP_FMUL, BOP_NONE, BOP_NONE},
	{T_FLOAT, SLASH, T_FLOAT}: {T_FLOAT, BOP_FDIV, BOP_NONE, BOP_NONE},
	{T_FLOAT, MOD, T_FLOAT}:   {T_FLOAT, BOP_FREM, BOP_NONE, BOP_NONE},

	{T_FLOAT, LESSER, T_FLOAT}:         {T_BOOLEAN, BOP_FCMP_LT, BOP_NONE, BOP_NONE},
	{T_FLOAT, GREATER, T_FLOAT}:        {T_BOOLEAN, BOP_FCMP_GT, BOP_NONE, BOP_NONE},
	{T_FLOAT, LESSER_EQUALS, T_FLOAT}:  {T_BOOLEAN, BOP_FCMP_LE, BOP_NONE, BOP_NONE},
	{T_FLOAT, GREATER_EQUALS, T_FLOAT}: {T_BOOLEAN, BOP_FCMP_GE, BOP_NONE, BOP_NONE},
	{T_FLOAT, EQUAL_EQUAL, T_FLOAT}:    {T_BOOLEAN, BOP_FCMP_EQ, BOP_NONE, BOP_NONE},
	{T_FLOAT, BANG_EQUAL, T_FLOAT}:     {T_BOOLEAN, BOP_FCMP_NE, BOP_NONE, BOP_NONE},

	{T_FLOAT, PLUS, T_INT}:  {T_FLOAT, BOP_FADD, BOP_NONE, BOP_I2F},
	{T_FLOAT, MINUS, T_INT}: {T_FLOAT, BOP_FSUB, BOP_NONE, BOP_I2F},
	{T_FLOAT, STAR, T_INT}:  {T_FLOAT, BOP_FMUL, BOP_NONE, BOP_I2F},
	{T_FLOAT, SLASH, T_INT}: {T_FLOAT, BOP_FDIV, BOP_NONE, BOP_I2F},
	{T_FLOAT, MOD, T_INT}:   {T_FLOAT, BOP_FREM, BOP_NONE, BOP_I2F},

	{T_FLOAT, LESSER, T_INT}:         {T_BOOLEAN, BOP_FCMP_LT, BOP_NONE, BOP_I2F},
	{T_FLOAT, GREATER, T_INT}:        {T_BOOLEAN, BOP_FCMP_GT, BOP_NONE, BOP_I2F},
	{T_FLOAT, LESSER_EQUALS, T_INT}:  {T_BOOLEAN, BOP_FCMP_LE, BOP_NONE, BOP_I2F},
	{T_FLOAT, GREATER_EQUALS, T_INT}: {T_BOOLEAN, BOP_FCMP_GE, BOP_NONE, BOP_I2F},
	{T_FLOAT, EQUAL_EQUAL, T_INT}:    {T_BOOLEAN, BOP_FCMP_EQ, BOP_NONE, BOP_I2F},
	{T_FLOAT, BANG_EQUAL, T_INT}:     {T_BOOLEAN, BOP_FCMP_NE, BOP_NONE, BOP_I2F},

	{T_STRING, PLUS, T_STRING}:        {T_STRING, BOP_S_CONCAT, BOP_NONE, BOP_NONE},
	{T_STRING, EQUAL_EQUAL, T_STRING}: {T_BOOLEAN, BOP_S_EQ, BOP_NONE, BOP_NONE},
	{T_STRING, BANG_EQUAL, T_STRING}:  {T_BOOLEAN, BOP_S_NE, BOP_NONE, BOP_NONE},

	{T_BOOLEAN, K_AND, T_BOOLEAN}:       {T_BOOLEAN, BOP_IAND, BOP_NONE, BOP_NONE},
	{T_BOOLEAN, K_OR, T_BOOLEAN}:        {T_BOOLEAN, BOP_IOR, BOP_NONE, BOP_NONE},
	{T_BOOLEAN, EQUAL_EQUAL, T_BOOLEAN}: {T_BOOLEAN, BOP_ICMP_EQ, BOP_NONE, BOP_NONE},
	{T_BOOLEAN, BANG_EQUAL, T_BOOLEAN}:  {T_BOOLEAN, BOP_ICMP_NE, BOP_NONE, BOP_NONE},
}

func NewTypeChecker(fileName string, program Program) *TypeChecker {
	defaultType := map[string]TypeId{
		"int":              T_INT,
		"float":            T_FLOAT,
		"boolean":          T_BOOLEAN,
		"java/lang/String": T_STRING,
	}
	return &TypeChecker{
		fileName:     fileName,
		program:      program,
		functions:    make(map[string][]string),
		typeRegister: defaultType,
	}
}

func (t *TypeChecker) resolve() error {
	for i, fnStmt := range t.program.functions {
		t.variableStack = make(map[string]Variable)
		err := t.handleFnBodyStatement(fnStmt.body)
		t.program.functions[i].maxLocals = len(t.variableStack)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TypeChecker) handleFnBodyStatement(body BlockStatement) error {
	for _, stmt := range body.body {
		switch stmt := stmt.(type) {
		case *ExpressionStatement:
			_, err := t.handleExpression(stmt.expression)
			if err != nil {
				return err
			}
		case *PrintStatement:
			// test
			t, err := t.handleExpression(stmt.expression)
			if err != nil {
				return err
			}
			switch t {
			case T_INT:
				stmt.resolvedArgs = "(I)V"
			case T_FLOAT:
				stmt.resolvedArgs = "(F)V"
			case T_BOOLEAN:
				stmt.resolvedArgs = "(Z)V"
			case T_STRING:
				stmt.resolvedArgs = "(Ljava/lang/String;)V"
			default:
				stmt.resolvedArgs = "(Ljava/lang/Object;)V"
			}
		default:
			fmt.Printf("DEV-INFO: Skipped %T Type for typechecking\n", stmt)
		}
	}
	return nil
}

func (t *TypeChecker) handleExpression(expr Expression) (TypeId, error) {
	switch expr := expr.(type) {
	case *IntegerLiteral:
		return T_INT, nil
	case *FloatLiteral:
		return T_FLOAT, nil
	case *StringLiteral:
		return T_STRING, nil
	case *BooleanLiteral:
		return T_BOOLEAN, nil
	case *BinaryExpression:
		return t.handleBinaryExpression(expr)
	case *AssignmentExpression:
		return t.handleAssignmentExpression(expr)
	case *Identifier:
		variable, exists := t.variableStack[expr.name]
		if !exists {
			return T_UNKNOWN, t.fail(
				fmt.Sprintf("Undefined variable: %s", expr.name),
				expr.line,
			)
		}
		expr.variable = &variable
		return variable.varType, nil
	default:
		panic(fmt.Sprintf("Unhandled Expression Type: %T", expr))
	}
}

func (t *TypeChecker) handleAssignmentExpression(expr *AssignmentExpression) (TypeId, error) {
	right, err := t.handleExpression(expr.right)
	if err != nil {
		return T_UNKNOWN, err
	}
	switch left := expr.left.(type) {
	case *Identifier:
		variable := Variable{
			varType: right,
			index:   len(t.variableStack),
		}
		left.variable = &variable
		t.variableStack[left.name] = variable
	default:
		panic(fmt.Sprintf("Unhandled Assignment Type: %T", expr.left))
	}
	return right, nil
}

func (t *TypeChecker) handleBinaryExpression(expr *BinaryExpression) (TypeId, error) {
	left, err := t.handleExpression(expr.left)
	if err != nil {
		return T_UNKNOWN, err
	}
	right, err := t.handleExpression(expr.right)
	if err != nil {
		return T_UNKNOWN, err
	}
	var exists bool
	key := BinaryTypeKey{left, expr.binaryType, right}
	expr.typeInfo, exists = binaryRules[key]
	if !exists {
		return T_UNKNOWN, t.fail("Invalid Binary Expression", expr.line)
	}
	return expr.typeInfo.result, nil
}

func (t *TypeChecker) fail(message string, line uint) error {
	return error(CompileError{
		Message:  message,
		FileName: t.fileName,
		Line:     line,
		Col:      0,
	})
}
