package utils

import (
	"regexp"
	"strings"
)

const (
	FILE_EXTENSION  = ".dt"
	TAGET_EXTENSION = ".class"
)

type TypeId int32
type TokenType int
type StatementType int

const (
	FUNCTION_STATEMENT = iota
	BLOCK_STATEMENT    = iota
)

const (
	X_EMPTY        TokenType = iota // Add this to handle the default in SubexpGroups
	COMMENT        TokenType = iota
	K_FUNCTION     TokenType = iota
	K_RETURN       TokenType = iota
	K_PRINT        TokenType = iota
	K_TRUE         TokenType = iota
	K_FALSE        TokenType = iota
	K_AND          TokenType = iota
	K_OR           TokenType = iota
	K_NOT          TokenType = iota
	SEMICOLON      TokenType = iota
	OPEN_BRACE     TokenType = iota
	CLOSE_BRACE    TokenType = iota
	OPEN_PARAM     TokenType = iota
	CLOSE_PARAM    TokenType = iota
	IDENTIFIER     TokenType = iota
	STRING         TokenType = iota
	FLOAT          TokenType = iota
	INTEGER        TokenType = iota
	PLUS           TokenType = iota
	MINUS          TokenType = iota
	STAR           TokenType = iota
	SLASH          TokenType = iota
	BANG_EQUAL     TokenType = iota
	EQUAL_EQUAL    TokenType = iota
	LESSER_EQUALS  TokenType = iota
	GREATER_EQUALS TokenType = iota
	LESSER         TokenType = iota
	GREATER        TokenType = iota
	NEWLINE        TokenType = iota
	K_EOF          TokenType = iota
)

var TOKEN_GRAMMAR = regexp.MustCompile(strings.Join([]string{
	`(?P<COMMENT>//[^\n]*)`,
	`(?P<KEYWORD_FUNCTION>\bfunction\b)`,
	`(?P<KETWORD_RETURN>\breturn\b)`,
	`(?P<KETWORD_PRINT>\bprint\b)`,
	`(?P<KETWORD_TRUE>\btrue\b)`,
	`(?P<KETWORD_FALSE>\bfalse\b)`,
	`(?P<KETWORD_AND>\band\b)`,
	`(?P<KETWORD_OR>\bor\b)`,
	`(?P<KETWORD_NOT>\bnot\b)`,
	`(?P<SEMI_COLON>\;)`,
	`(?P<OPEN_BRACE>\{)`,
	`(?P<CLOSE_BRACE>\})`,
	`(?P<OPEN_PARAM>\()`,
	`(?P<CLOSE_PARAM>\))`,
	`(?P<IDENTIFIER>[A-Za-z_][A-Za-z0-9_]*)`,
	`(?P<STRING_LITERAL>"(?:\\.|[^"\\])*")`,
	`(?P<FLOAT_LITERAL>[0-9]+\.[0-9]+)`,
	`(?P<INTEGER_LITERAL>[0-9]+)`,
	`(?P<PLUS>\+)`,
	`(?P<MINUS>-)`,
	`(?P<STAR>\*)`,
	`(?P<SLASH>/)`,
	`(?P<BANG_EQUAL>!=)`,
	`(?P<EQUAL_EQUAL>==)`,
	`(?P<LESSER_EQUALS><=)`,
	`(?P<GREATER_EQUALS>>=)`,
	`(?P<LESSER><)`,
	`(?P<GREATER>>)`,
	`(?P<NEWLINE>\n)`,
}, "|"))

var GROUP_NAMES = TOKEN_GRAMMAR.SubexpNames()

const (
	T_UNKNOWN TypeId = iota
	T_INT     TypeId = iota
	T_FLOAT   TypeId = iota
	T_STRING  TypeId = iota
	T_BOOLEAN TypeId = iota
)

const (
	OP_RETURN          = 0xB1
	OP_LOAD_CONSTANT   = 0x12
	OP_LOAD_CONSTANT_W = 0x13
	OP_INVOKE_STATIC   = 0xB8
	OP_IADD            = 0x60
	OP_FADD            = 0x62
	OP_ISUB            = 0x64
	OP_FSUB            = 0x66
	OP_I2F             = 0x86
)
