package parser

import (
	"fmt"

	"donut/lang"
)

// primary types
const (
	typeUnknown TypeId = iota
	typeInt
	typeFloat
	typeString
	typeBoolean
)

// expression operation
const (
	uop_ineg ExpOpTypeId = iota
	uop_fneg
	uop_not
	bop_none
	bop_iadd
	bop_isub
	bop_imul
	bop_idiv
	bop_irem
	bop_fadd
	bop_fsub
	bop_fmul
	bop_fdiv
	bop_frem
	bop_i2f
	bop_icmp_eq
	bop_icmp_ne
	bop_icmp_lt
	bop_icmp_le
	bop_icmp_gt
	bop_icmp_ge
	bop_fcmp_eq
	bop_fcmp_ne
	bop_fcmp_lt
	bop_fcmp_le
	bop_fcmp_gt
	bop_fcmp_ge
	bop_iand
	bop_ior
	bop_s_concat
	bop_s_eq
	bop_s_ne
)

type TypeId int32
type ExpOpTypeId int

type UnaryTypeKey struct {
	operation TokenType
	right     TypeId
}

type BinaryTypeKey struct {
	left      TypeId
	operation TokenType
	right     TypeId
}

type BinaryTypeValue struct {
	result        TypeId
	operationType ExpOpTypeId
	leftTo        ExpOpTypeId
	rightTo       ExpOpTypeId
}

var unaryRules = map[UnaryTypeKey]ExpOpTypeId{
	{t_minus, typeInt}:   uop_ineg,
	{t_minus, typeFloat}: uop_fneg,
	{k_not, typeBoolean}: uop_not,
}

var binaryRules = map[BinaryTypeKey]BinaryTypeValue{
	{typeInt, t_plus, typeInt}:               {typeInt, bop_iadd, bop_none, bop_none},
	{typeInt, t_minus, typeInt}:              {typeInt, bop_isub, bop_none, bop_none},
	{typeInt, t_star, typeInt}:               {typeInt, bop_imul, bop_none, bop_none},
	{typeInt, t_slash, typeInt}:              {typeInt, bop_idiv, bop_none, bop_none},
	{typeInt, t_mod, typeInt}:                {typeInt, bop_irem, bop_none, bop_none},
	{typeInt, t_lesser, typeInt}:             {typeBoolean, bop_icmp_lt, bop_none, bop_none},
	{typeInt, t_greater, typeInt}:            {typeBoolean, bop_icmp_gt, bop_none, bop_none},
	{typeInt, t_lesserEquals, typeInt}:       {typeBoolean, bop_icmp_le, bop_none, bop_none},
	{typeInt, t_greaterEquals, typeInt}:      {typeBoolean, bop_icmp_ge, bop_none, bop_none},
	{typeInt, t_equalEqual, typeInt}:         {typeBoolean, bop_icmp_eq, bop_none, bop_none},
	{typeInt, t_bangEqual, typeInt}:          {typeBoolean, bop_icmp_ne, bop_none, bop_none},
	{typeInt, t_plus, typeFloat}:             {typeFloat, bop_fadd, bop_i2f, bop_none},
	{typeInt, t_minus, typeFloat}:            {typeFloat, bop_fsub, bop_i2f, bop_none},
	{typeInt, t_star, typeFloat}:             {typeFloat, bop_fmul, bop_i2f, bop_none},
	{typeInt, t_slash, typeFloat}:            {typeFloat, bop_fdiv, bop_i2f, bop_none},
	{typeInt, t_mod, typeFloat}:              {typeFloat, bop_frem, bop_i2f, bop_none},
	{typeInt, t_lesser, typeFloat}:           {typeBoolean, bop_fcmp_lt, bop_i2f, bop_none},
	{typeInt, t_greater, typeFloat}:          {typeBoolean, bop_fcmp_gt, bop_i2f, bop_none},
	{typeInt, t_lesserEquals, typeFloat}:     {typeBoolean, bop_fcmp_le, bop_i2f, bop_none},
	{typeInt, t_greaterEquals, typeFloat}:    {typeBoolean, bop_fcmp_ge, bop_i2f, bop_none},
	{typeInt, t_equalEqual, typeFloat}:       {typeBoolean, bop_fcmp_eq, bop_i2f, bop_none},
	{typeInt, t_bangEqual, typeFloat}:        {typeBoolean, bop_fcmp_ne, bop_i2f, bop_none},
	{typeFloat, t_plus, typeFloat}:           {typeFloat, bop_fadd, bop_none, bop_none},
	{typeFloat, t_minus, typeFloat}:          {typeFloat, bop_fsub, bop_none, bop_none},
	{typeFloat, t_star, typeFloat}:           {typeFloat, bop_fmul, bop_none, bop_none},
	{typeFloat, t_slash, typeFloat}:          {typeFloat, bop_fdiv, bop_none, bop_none},
	{typeFloat, t_mod, typeFloat}:            {typeFloat, bop_frem, bop_none, bop_none},
	{typeFloat, t_lesser, typeFloat}:         {typeBoolean, bop_fcmp_lt, bop_none, bop_none},
	{typeFloat, t_greater, typeFloat}:        {typeBoolean, bop_fcmp_gt, bop_none, bop_none},
	{typeFloat, t_lesserEquals, typeFloat}:   {typeBoolean, bop_fcmp_le, bop_none, bop_none},
	{typeFloat, t_greaterEquals, typeFloat}:  {typeBoolean, bop_fcmp_ge, bop_none, bop_none},
	{typeFloat, t_equalEqual, typeFloat}:     {typeBoolean, bop_fcmp_eq, bop_none, bop_none},
	{typeFloat, t_bangEqual, typeFloat}:      {typeBoolean, bop_fcmp_ne, bop_none, bop_none},
	{typeFloat, t_plus, typeInt}:             {typeFloat, bop_fadd, bop_none, bop_i2f},
	{typeFloat, t_minus, typeInt}:            {typeFloat, bop_fsub, bop_none, bop_i2f},
	{typeFloat, t_star, typeInt}:             {typeFloat, bop_fmul, bop_none, bop_i2f},
	{typeFloat, t_slash, typeInt}:            {typeFloat, bop_fdiv, bop_none, bop_i2f},
	{typeFloat, t_mod, typeInt}:              {typeFloat, bop_frem, bop_none, bop_i2f},
	{typeFloat, t_lesser, typeInt}:           {typeBoolean, bop_fcmp_lt, bop_none, bop_i2f},
	{typeFloat, t_greater, typeInt}:          {typeBoolean, bop_fcmp_gt, bop_none, bop_i2f},
	{typeFloat, t_lesserEquals, typeInt}:     {typeBoolean, bop_fcmp_le, bop_none, bop_i2f},
	{typeFloat, t_greaterEquals, typeInt}:    {typeBoolean, bop_fcmp_ge, bop_none, bop_i2f},
	{typeFloat, t_equalEqual, typeInt}:       {typeBoolean, bop_fcmp_eq, bop_none, bop_i2f},
	{typeFloat, t_bangEqual, typeInt}:        {typeBoolean, bop_fcmp_ne, bop_none, bop_i2f},
	{typeString, t_plus, typeString}:         {typeString, bop_s_concat, bop_none, bop_none},
	{typeString, t_equalEqual, typeString}:   {typeBoolean, bop_s_eq, bop_none, bop_none},
	{typeString, t_bangEqual, typeString}:    {typeBoolean, bop_s_ne, bop_none, bop_none},
	{typeBoolean, k_and, typeBoolean}:        {typeBoolean, bop_iand, bop_none, bop_none},
	{typeBoolean, k_or, typeBoolean}:         {typeBoolean, bop_ior, bop_none, bop_none},
	{typeBoolean, t_equalEqual, typeBoolean}: {typeBoolean, bop_icmp_eq, bop_none, bop_none},
	{typeBoolean, t_bangEqual, typeBoolean}:  {typeBoolean, bop_icmp_ne, bop_none, bop_none},
}

type Analyzer struct {
	fileName      string
	program       Program
	functions     map[string][]string
	variableStack map[string]Variable // todo: convert this to a stack which will holder the scoped variables.
	typeRegister  map[string]TypeId
}

func NewAnalyzer(fileName string, program Program) *Analyzer {
	return &Analyzer{
		fileName:  fileName,
		program:   program,
		functions: make(map[string][]string),
		typeRegister: map[string]TypeId{
			"int":              typeInt,
			"float":            typeFloat,
			"boolean":          typeBoolean,
			"java/lang/string": typeString,
		},
	}
}

func (a *Analyzer) Analyze() error {
	for i, fnStmt := range a.program.functions {
		a.variableStack = make(map[string]Variable)
		err := a.handleStatement(&fnStmt.body)
		a.program.functions[i].maxLocals = len(a.variableStack) + 1 // todo: unhardcode this
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyzer) handleStatement(stmt Statement) error {
	switch stmt := stmt.(type) {
	case *ExpressionStatement:
		_, err := a.handleExpression(stmt.expression)
		if err != nil {
			return err
		}
	case *BlockStatement:
		err := a.handleBlockstmt(*stmt)
		if err != nil {
			return err
		}
	case *WhileStatement:
		err := a.handleWhileStmt(*stmt)
		if err != nil {
			return err
		}
	case *IfStatement:
		err := a.handleIfStmt(*stmt)
		if err != nil {
			return err
		}
	case *PrintStatement:
		// test
		t, err := a.handleExpression(stmt.expression)
		if err != nil {
			return err
		}
		switch t {
		case typeInt:
			stmt.resolvedArgs = "(I)V"
		case typeFloat:
			stmt.resolvedArgs = "(F)V"
		case typeBoolean:
			stmt.resolvedArgs = "(Z)V"
		case typeString:
			stmt.resolvedArgs = "(Ljava/lang/String;)V"
		default:
			stmt.resolvedArgs = "(Ljava/lang/Object;)V"
		}
	default:
		fmt.Printf("DEV-INFO: Skipped %T Type for typechecking\n", stmt)
	}
	return nil
}

func (a *Analyzer) handleBlockstmt(stmt BlockStatement) error {
	for _, stmt := range stmt.body {
		err := a.handleStatement(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyzer) handleIfStmt(stmt IfStatement) error {
	exprType, err := a.handleExpression(stmt.test)
	if err != nil {
		return err
	}
	if exprType != typeBoolean {
		return a.fail("Required a boolean test expression for the if statement", stmt.line)
	}
	err = a.handleStatement(stmt.consequent)
	if err != nil {
		return nil
	}
	if stmt.alternate != nil {
		return a.handleStatement(stmt.alternate)
	}
	return nil
}

func (a *Analyzer) handleWhileStmt(stmt WhileStatement) error {
	exprType, err := a.handleExpression(stmt.test)
	if err != nil {
		return err
	}
	if exprType != typeBoolean {
		return a.fail("Required a boolean test expression for the while loop", stmt.line)
	}
	return a.handleStatement(stmt.body)
}

func (a *Analyzer) handleExpression(expr Expression) (TypeId, error) {
	switch expr := expr.(type) {
	case *IntegerLiteral:
		return typeInt, nil
	case *FloatLiteral:
		return typeFloat, nil
	case *StringLiteral:
		return typeString, nil
	case *BooleanLiteral:
		return typeBoolean, nil
	case *UnaryExpression:
		return a.handleUnaryExpression(expr)
	case *BinaryExpression:
		return a.handleBinaryExpression(expr)
	case *AssignmentExpression:
		return a.handleAssignmentExpression(expr)
	case *Identifier:
		variable, exists := a.variableStack[expr.name]
		if !exists {
			return typeUnknown, a.fail(fmt.Sprintf("Undefined variable: %s", expr.name), expr.line)
		}
		expr.variable = &variable
		return variable.varType, nil
	default:
		panic(fmt.Sprintf("Unhandled Expression Type: %T", expr))
	}
}

func (a *Analyzer) handleAssignmentExpression(expr *AssignmentExpression) (TypeId, error) {
	right, err := a.handleExpression(expr.right)
	if err != nil {
		return typeUnknown, err
	}
	switch left := expr.left.(type) {
	case *Identifier:
		variable, exists := a.variableStack[left.name]
		if !exists {
			variable = Variable{
				varType: right,
				index:   len(a.variableStack),
			}
		}
		left.variable = &variable
		a.variableStack[left.name] = variable
	default:
		panic(fmt.Sprintf("Unhandled Assignment Type: %T", expr.left))
	}
	return right, nil
}

func (a *Analyzer) handleUnaryExpression(expr *UnaryExpression) (TypeId, error) {
	right, err := a.handleExpression(expr.value)
	if err != nil {
		return typeUnknown, err
	}
	key := UnaryTypeKey{expr.unaryType, right}
	result, exists := unaryRules[key]
	if !exists {
		return typeUnknown, a.fail("Invalid Unary Expression", expr.line)
	}
	expr.typeInfo = result
	return right, nil
}

func (a *Analyzer) handleBinaryExpression(expr *BinaryExpression) (TypeId, error) {
	left, err := a.handleExpression(expr.left)
	if err != nil {
		return typeUnknown, err
	}
	right, err := a.handleExpression(expr.right)
	if err != nil {
		return typeUnknown, err
	}
	var exists bool
	key := BinaryTypeKey{left, expr.binaryType, right}
	expr.typeInfo, exists = binaryRules[key]
	if !exists {
		return typeUnknown, a.fail("Invalid Binary Expression", expr.line)
	}
	return expr.typeInfo.result, nil
}

func (a *Analyzer) fail(message string, line uint) error {
	return error(lang.CompileError{
		Message:  message,
		FileName: a.fileName,
		Line:     line,
		Col:      0,
	})
}
