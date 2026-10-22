package utils

import (
	"fmt"
)

type TypeChecker struct {
	fileName     string
	program      Program
	functions    map[string][]string
	typeRegister map[string]TypeId
}

var binaryRules = map[Pair[TypeId, TypeId]]BinaryTypeResult{
	{T_INT, T_INT}:     {T_INT, nil, nil},
	{T_INT, T_FLOAT}:   {T_FLOAT, (*ByteCode).opI2F, nil},
	{T_FLOAT, T_INT}:   {T_FLOAT, nil, (*ByteCode).opI2F},
	{T_FLOAT, T_FLOAT}: {T_FLOAT, nil, nil},
}

var EXPR_OPERATION_TABLE = map[Pair[TypeId, TokenType]]func(*ByteCode){
	{T_INT, PLUS}:    (*ByteCode).opIAdd,
	{T_INT, MINUS}:   (*ByteCode).opISub,
	{T_FLOAT, PLUS}:  (*ByteCode).opFAdd,
	{T_FLOAT, MINUS}: (*ByteCode).opFSub,
}

func NewTypeChecker(fileName string, program Program) *TypeChecker {
	defaultType := map[string]TypeId{
		"int":     T_INT,
		"float":   T_FLOAT,
		"string":  T_STRING,
		"boolean": T_BOOLEAN,
	}
	return &TypeChecker{
		fileName:     fileName,
		program:      program,
		functions:    make(map[string][]string),
		typeRegister: defaultType,
	}
}

func (t *TypeChecker) resolve() error {
	for _, fnStmt := range t.program.functions {
		err := t.handleFnBodyStatement(fnStmt.body)
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
			_, err := t.handleExpression(stmt.expression)
			if err != nil {
				return err
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
	}
	return T_UNKNOWN, nil
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
	pair := Pair[TypeId, TypeId]{left, right}
	expr.typeInfo, exists = binaryRules[pair]
	if !exists {
		return T_UNKNOWN, t.fail("Invalid Binary Expression", expr.line)
	}
	opKey := Pair[TypeId, TokenType]{expr.typeInfo.result, expr.binaryType}
	expr.resolvedOpType, exists = EXPR_OPERATION_TABLE[opKey]
	if !exists {
		// this should be unreachable when all the rules are in place.
		panic(fmt.Sprintf("Missing binary operation implementation for %v", opKey))
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
