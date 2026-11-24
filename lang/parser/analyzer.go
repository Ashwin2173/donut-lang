package parser

import (
	"fmt"
	"slices"

	"donut/lang"
)

type Context struct {
	variableStack   []map[string]*Variable
	functionPointer *FunctionStatement
}

func NewContext(function *FunctionStatement) *Context {
	return &Context{
		variableStack:   make([]map[string]*Variable, 0),
		functionPointer: function,
	}
}

type Analyzer struct {
	filePath     string
	program      Program
	functions    map[string][]string
	context      *Context
	typeRegister map[string]TypeId
}

func NewAnalyzer(filePath string, program Program) *Analyzer {
	return &Analyzer{
		filePath:  filePath,
		program:   program,
		functions: make(map[string][]string),
		typeRegister: map[string]TypeId{
			"void":             TypeVoid,
			"int":              TypeInt,
			"float":            TypeFloat,
			"boolean":          TypeBoolean,
			"java/lang/string": TypeString,
		},
	}
}

func (a *Analyzer) Analyze() error {
	for _, fnStmt := range a.program.Functions {
		a.context = NewContext(&fnStmt)
		err := a.handleStatement(&fnStmt.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Analyzer) handleStatement(stmt Statement) error {
	var err error = nil
	switch stmt := stmt.(type) {
	case *ExpressionStatement:
		_, err = a.handleExpression(stmt.Expression)
	case *BlockStatement:
		err = a.handleBlockstmt(*stmt)
	case *WhileStatement:
		err = a.handleWhileStmt(*stmt)
	case *IfStatement:
		err = a.handleIfStmt(*stmt)
	case *PrintStatement:
		// test
		var t TypeId
		t, err = a.handleExpression(stmt.Expression)
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
	return err
}

func (a *Analyzer) handleBlockstmt(stmt BlockStatement) error {
	a.context.variableStack = append(a.context.variableStack, make(map[string]*Variable))
	for _, stmt := range stmt.Body {
		if err := a.handleStatement(stmt); err != nil {
			return err
		}
	}
	a.context.variableStack = a.context.variableStack[:len(a.context.variableStack)-1]
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
		return a.getIdentifier(expr)
	default:
		panic(fmt.Sprintf("Unhandled Expression Type: %T", expr))
	}
}

func (a *Analyzer) handleAssignmentExpression(expr *AssignmentExpression) (TypeId, error) {
	rightType, err := a.handleExpression(expr.Right)
	if err != nil {
		return TypeUnknown, err
	}
	switch expr.Left.(type) {
	case *Identifier:
		return a.handleVariableAssignment(expr, rightType)
	default:
		panic(fmt.Sprintf("Unhandled Assignment Type: %T", expr.Left))
	}
}

func (a *Analyzer) handleVariableAssignment(expr *AssignmentExpression, rightType TypeId) (TypeId, error) {
	left := expr.Left.(*Identifier)
	switch expr.binaryType {
	case t_colonEqual: // new variable
		expr.IsNew = true
		left.Variable = &Variable{
			VarType: rightType,
		}
		a.context.variableStack[len(a.context.variableStack)-1][left.Name] = left.Variable
		return rightType, nil
	case t_equal: // assignment to existing
		leftType, err := a.getIdentifier(left)
		if err != nil {
			return TypeUnknown, err
		}
		if leftType != rightType {
			return TypeUnknown, a.fail("Assignment type mismatch", expr.line)
		}
		return leftType, nil
	default:
		panic("RUC at parser/analyzer/handleAssignmentExpression: Parser failed to catch this")
	}
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

func (a *Analyzer) getIdentifier(id *Identifier) (TypeId, error) {
	var variable *Variable
	for _, frame := range slices.Backward(a.context.variableStack) {
		var exists bool
		variable, exists = frame[id.Name]
		if exists {
			break
		}
	}
	if variable == nil {
		return TypeUnknown, a.fail(fmt.Sprintf("Undefined variable: %s", id.Name), id.line)
	}
	id.Variable = variable
	return variable.VarType, nil
}

func (a *Analyzer) fail(message string, line uint) error {
	return error(lang.CompileError{
		Message:  message,
		FileName: a.filePath,
		Line:     line,
		Col:      0,
	})
}
