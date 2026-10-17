package utils

import (
	"regexp"
	"strings"
)

const (
	FILE_EXTENSION = ".dt"
	TAGET_EXTENSION = ".class"
)

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

// Compiler const
const (
	OP_RETURN = 0xB1
)
