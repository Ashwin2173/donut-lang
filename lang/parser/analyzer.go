package parser

import (
	"fmt"

	"donut/lang"
)

// primary types
const (
	TypeUnknown TypeId = iota
	TypeInt
	TypeFloat
	TypeString
	TypeBoolean
)

// expression operation
const (
	Uop_ineg ExpOpTypeId = iota
	Uop_fneg
	Uop_not
	Bop_none
	Bop_iadd
	Bop_isub
	Bop_imul
	Bop_idiv
	Bop_irem
	Bop_fadd
	Bop_fsub
	Bop_fmul
	Bop_fdiv
	Bop_frem
	Bop_i2f
	Bop_icmp_eq
	Bop_icmp_ne
	Bop_icmp_lt
	Bop_icmp_le
	Bop_icmp_gt
	Bop_icmp_ge
	Bop_fcmp_eq
	Bop_fcmp_ne
	Bop_fcmp_lt
	Bop_fcmp_le
	Bop_fcmp_gt
	Bop_fcmp_ge
	Bop_iand
	Bop_ior
	Bop_s_concat
	Bop_s_eq
	Bop_s_ne
)

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
	OperationType ExpOpTypeId
	LeftTo        ExpOpTypeId
	RightTo       ExpOpTypeId
}

var unaryRules = map[UnaryTypeKey]ExpOpTypeId{
	{t_minus, TypeInt}:   Uop_ineg,
	{t_minus, TypeFloat}: Uop_fneg,
	{k_not, TypeBoolean}: Uop_not,
}

var binaryRules = map[BinaryTypeKey]BinaryTypeValue{
	{TypeInt, t_plus, TypeInt}:               {TypeInt, Bop_iadd, Bop_none, Bop_none},
	{TypeInt, t_minus, TypeInt}:              {TypeInt, Bop_isub, Bop_none, Bop_none},
	{TypeInt, t_star, TypeInt}:               {TypeInt, Bop_imul, Bop_none, Bop_none},
	{TypeInt, t_slash, TypeInt}:              {TypeInt, Bop_idiv, Bop_none, Bop_none},
	{TypeInt, t_mod, TypeInt}:                {TypeInt, Bop_irem, Bop_none, Bop_none},
	{TypeInt, t_lesser, TypeInt}:             {TypeBoolean, Bop_icmp_lt, Bop_none, Bop_none},
	{TypeInt, t_greater, TypeInt}:            {TypeBoolean, Bop_icmp_gt, Bop_none, Bop_none},
	{TypeInt, t_lesserEquals, TypeInt}:       {TypeBoolean, Bop_icmp_le, Bop_none, Bop_none},
	{TypeInt, t_greaterEquals, TypeInt}:      {TypeBoolean, Bop_icmp_ge, Bop_none, Bop_none},
	{TypeInt, t_equalEqual, TypeInt}:         {TypeBoolean, Bop_icmp_eq, Bop_none, Bop_none},
	{TypeInt, t_bangEqual, TypeInt}:          {TypeBoolean, Bop_icmp_ne, Bop_none, Bop_none},
	{TypeInt, t_plus, TypeFloat}:             {TypeFloat, Bop_fadd, Bop_i2f, Bop_none},
	{TypeInt, t_minus, TypeFloat}:            {TypeFloat, Bop_fsub, Bop_i2f, Bop_none},
	{TypeInt, t_star, TypeFloat}:             {TypeFloat, Bop_fmul, Bop_i2f, Bop_none},
	{TypeInt, t_slash, TypeFloat}:            {TypeFloat, Bop_fdiv, Bop_i2f, Bop_none},
	{TypeInt, t_mod, TypeFloat}:              {TypeFloat, Bop_frem, Bop_i2f, Bop_none},
	{TypeInt, t_lesser, TypeFloat}:           {TypeBoolean, Bop_fcmp_lt, Bop_i2f, Bop_none},
	{TypeInt, t_greater, TypeFloat}:          {TypeBoolean, Bop_fcmp_gt, Bop_i2f, Bop_none},
	{TypeInt, t_lesserEquals, TypeFloat}:     {TypeBoolean, Bop_fcmp_le, Bop_i2f, Bop_none},
	{TypeInt, t_greaterEquals, TypeFloat}:    {TypeBoolean, Bop_fcmp_ge, Bop_i2f, Bop_none},
	{TypeInt, t_equalEqual, TypeFloat}:       {TypeBoolean, Bop_fcmp_eq, Bop_i2f, Bop_none},
	{TypeInt, t_bangEqual, TypeFloat}:        {TypeBoolean, Bop_fcmp_ne, Bop_i2f, Bop_none},
	{TypeFloat, t_plus, TypeFloat}:           {TypeFloat, Bop_fadd, Bop_none, Bop_none},
	{TypeFloat, t_minus, TypeFloat}:          {TypeFloat, Bop_fsub, Bop_none, Bop_none},
	{TypeFloat, t_star, TypeFloat}:           {TypeFloat, Bop_fmul, Bop_none, Bop_none},
	{TypeFloat, t_slash, TypeFloat}:          {TypeFloat, Bop_fdiv, Bop_none, Bop_none},
	{TypeFloat, t_mod, TypeFloat}:            {TypeFloat, Bop_frem, Bop_none, Bop_none},
	{TypeFloat, t_lesser, TypeFloat}:         {TypeBoolean, Bop_fcmp_lt, Bop_none, Bop_none},
	{TypeFloat, t_greater, TypeFloat}:        {TypeBoolean, Bop_fcmp_gt, Bop_none, Bop_none},
	{TypeFloat, t_lesserEquals, TypeFloat}:   {TypeBoolean, Bop_fcmp_le, Bop_none, Bop_none},
	{TypeFloat, t_greaterEquals, TypeFloat}:  {TypeBoolean, Bop_fcmp_ge, Bop_none, Bop_none},
	{TypeFloat, t_equalEqual, TypeFloat}:     {TypeBoolean, Bop_fcmp_eq, Bop_none, Bop_none},
	{TypeFloat, t_bangEqual, TypeFloat}:      {TypeBoolean, Bop_fcmp_ne, Bop_none, Bop_none},
	{TypeFloat, t_plus, TypeInt}:             {TypeFloat, Bop_fadd, Bop_none, Bop_i2f},
	{TypeFloat, t_minus, TypeInt}:            {TypeFloat, Bop_fsub, Bop_none, Bop_i2f},
	{TypeFloat, t_star, TypeInt}:             {TypeFloat, Bop_fmul, Bop_none, Bop_i2f},
	{TypeFloat, t_slash, TypeInt}:            {TypeFloat, Bop_fdiv, Bop_none, Bop_i2f},
	{TypeFloat, t_mod, TypeInt}:              {TypeFloat, Bop_frem, Bop_none, Bop_i2f},
	{TypeFloat, t_lesser, TypeInt}:           {TypeBoolean, Bop_fcmp_lt, Bop_none, Bop_i2f},
	{TypeFloat, t_greater, TypeInt}:          {TypeBoolean, Bop_fcmp_gt, Bop_none, Bop_i2f},
	{TypeFloat, t_lesserEquals, TypeInt}:     {TypeBoolean, Bop_fcmp_le, Bop_none, Bop_i2f},
	{TypeFloat, t_greaterEquals, TypeInt}:    {TypeBoolean, Bop_fcmp_ge, Bop_none, Bop_i2f},
	{TypeFloat, t_equalEqual, TypeInt}:       {TypeBoolean, Bop_fcmp_eq, Bop_none, Bop_i2f},
	{TypeFloat, t_bangEqual, TypeInt}:        {TypeBoolean, Bop_fcmp_ne, Bop_none, Bop_i2f},
	{TypeString, t_plus, TypeString}:         {TypeString, Bop_s_concat, Bop_none, Bop_none},
	{TypeString, t_equalEqual, TypeString}:   {TypeBoolean, Bop_s_eq, Bop_none, Bop_none},
	{TypeString, t_bangEqual, TypeString}:    {TypeBoolean, Bop_s_ne, Bop_none, Bop_none},
	{TypeBoolean, k_and, TypeBoolean}:        {TypeBoolean, Bop_iand, Bop_none, Bop_none},
	{TypeBoolean, k_or, TypeBoolean}:         {TypeBoolean, Bop_ior, Bop_none, Bop_none},
	{TypeBoolean, t_equalEqual, TypeBoolean}: {TypeBoolean, Bop_icmp_eq, Bop_none, Bop_none},
	{TypeBoolean, t_bangEqual, TypeBoolean}:  {TypeBoolean, Bop_icmp_ne, Bop_none, Bop_none},
}

type Analyzer struct {
	filePath      string
	program       Program
	functions     map[string][]string
	variableStack map[string]Variable // todo: convert this to a stack which will holder the scoped variables.
	typeRegister  map[string]TypeId
}

func NewAnalyzer(filePath string, program Program) *Analyzer {
	return &Analyzer{
		filePath:  filePath,
		program:   program,
		functions: make(map[string][]string),
		typeRegister: map[string]TypeId{
			"int":              TypeInt,
			"float":            TypeFloat,
			"boolean":          TypeBoolean,
			"java/lang/string": TypeString,
		},
	}
}

func (a *Analyzer) Analyze() error {
	for _, fnStmt := range a.program.Functions {
		a.variableStack = make(map[string]Variable)
		err := a.handleStatement(&fnStmt.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyzer) handleStatement(stmt Statement) error {
	switch stmt := stmt.(type) {
	case *ExpressionStatement:
		_, err := a.handleExpression(stmt.Expression)
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
		t, err := a.handleExpression(stmt.Expression)
		if err != nil {
			return err
		}
		switch t {
		case TypeInt:
			stmt.ResolvedArgs = "(I)V"
		case TypeFloat:
			stmt.ResolvedArgs = "(F)V"
		case TypeBoolean:
			stmt.ResolvedArgs = "(Z)V"
		case TypeString:
			stmt.ResolvedArgs = "(Ljava/lang/String;)V"
		default:
			stmt.ResolvedArgs = "(Ljava/lang/Object;)V"
		}
	default:
		fmt.Printf("DEV-INFO: Skipped %T Type for typechecking\n", stmt)
	}
	return nil
}

func (a *Analyzer) handleBlockstmt(stmt BlockStatement) error {
	for _, stmt := range stmt.Body {
		err := a.handleStatement(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyzer) handleIfStmt(stmt IfStatement) error {
	exprType, err := a.handleExpression(stmt.Test)
	if err != nil {
		return err
	}
	if exprType != TypeBoolean {
		return a.fail("Required a boolean test expression for the if statement", stmt.line)
	}
	err = a.handleStatement(stmt.Consequent)
	if err != nil {
		return err
	}
	if stmt.Alternate != nil {
		return a.handleStatement(stmt.Alternate)
	}
	return nil
}

func (a *Analyzer) handleWhileStmt(stmt WhileStatement) error {
	exprType, err := a.handleExpression(stmt.Test)
	if err != nil {
		return err
	}
	if exprType != TypeBoolean {
		return a.fail("Required a boolean test expression for the while loop", stmt.line)
	}
	return a.handleStatement(stmt.Body)
}

func (a *Analyzer) handleExpression(expr Expression) (TypeId, error) {
	switch expr := expr.(type) {
	case *IntegerLiteral:
		return TypeInt, nil
	case *FloatLiteral:
		return TypeFloat, nil
	case *StringLiteral:
		return TypeString, nil
	case *BooleanLiteral:
		return TypeBoolean, nil
	case *UnaryExpression:
		return a.handleUnaryExpression(expr)
	case *BinaryExpression:
		return a.handleBinaryExpression(expr)
	case *AssignmentExpression:
		return a.handleAssignmentExpression(expr)
	case *Identifier:
		variable, exists := a.variableStack[expr.name]
		if !exists {
			return TypeUnknown, a.fail(fmt.Sprintf("Undefined variable: %s", expr.name), expr.line)
		}
		expr.Variable = &variable
		return variable.VarType, nil
	default:
		panic(fmt.Sprintf("Unhandled Expression Type: %T", expr))
	}
}

func (a *Analyzer) handleAssignmentExpression(expr *AssignmentExpression) (TypeId, error) {
	right, err := a.handleExpression(expr.Right)
	if err != nil {
		return TypeUnknown, err
	}
	switch left := expr.Left.(type) {
	case *Identifier:
		variable, exists := a.variableStack[left.name]
		if !exists {
			variable = Variable{
				VarType: right,
				Index:   len(a.variableStack),
			}
		}
		left.Variable = &variable
		a.variableStack[left.name] = variable
	default:
		panic(fmt.Sprintf("Unhandled Assignment Type: %T", expr.Left))
	}
	return right, nil
}

func (a *Analyzer) handleUnaryExpression(expr *UnaryExpression) (TypeId, error) {
	right, err := a.handleExpression(expr.Value)
	if err != nil {
		return TypeUnknown, err
	}
	key := UnaryTypeKey{expr.unaryType, right}
	result, exists := unaryRules[key]
	if !exists {
		return TypeUnknown, a.fail("Invalid Unary Expression", expr.line)
	}
	expr.TypeInfo = result
	return right, nil
}

func (a *Analyzer) handleBinaryExpression(expr *BinaryExpression) (TypeId, error) {
	left, err := a.handleExpression(expr.Left)
	if err != nil {
		return TypeUnknown, err
	}
	right, err := a.handleExpression(expr.Right)
	if err != nil {
		return TypeUnknown, err
	}
	var exists bool
	key := BinaryTypeKey{left, expr.binaryType, right}
	expr.TypeInfo, exists = binaryRules[key]
	if !exists {
		return TypeUnknown, a.fail("Invalid Binary Expression", expr.line)
	}
	return expr.TypeInfo.result, nil
}

func (a *Analyzer) fail(message string, line uint) error {
	return error(lang.CompileError{
		Message:  message,
		FileName: a.filePath,
		Line:     line,
		Col:      0,
	})
}
