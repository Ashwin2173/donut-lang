package utils

import (
	"strings"
	"regexp"
)

type TokenType 		int
type StatementType 	int

const (
	FUNCTION_STATEMENT 		= iota
	BLOCK_STATEMENT 		= iota
)

const (
	X_EMPTY			TokenType = iota // Add this to handle the default in SubexpGroups
	COMMENT			TokenType = iota
	K_FUNCTION		TokenType = iota
	K_RETURN		TokenType = iota
	K_TRUE			TokenType = iota
	K_FALSE			TokenType = iota
	SEMICOLON		TokenType = iota
	OPEN_BRACE		TokenType = iota
	CLOSE_BRACE		TokenType = iota
	OPEN_PARAM		TokenType = iota
	CLOSE_PARAM		TokenType = iota
	IDENTIFIER		TokenType = iota
	STRING			TokenType = iota
	FLOAT			TokenType = iota
	INTEGER			TokenType = iota
	NEWLINE			TokenType = iota
	K_EOF			TokenType = iota
)

var TOKEN_GRAMMAR = regexp.MustCompile(strings.Join([]string{
	`(?P<COMMENT>//[^\n]*)`,
	`(?P<KEYWORD_FUNCTION>\bfunction\b)`,
	`(?P<KETWORD_RETURN>\breturn\b)`,
	`(?P<KETWORD_TRUE>\btrue\b)`,
	`(?P<KETWORD_FALSE>\bfalse\b)`,
	`(?P<SEMI_COLON>\;)`,
	`(?P<OPEN_BRACE>\{)`,
	`(?P<CLOSE_BRACE>\})`,
	`(?P<OPEN_PARAM>\()`,
	`(?P<CLOSE_PARAM>\))`,
	`(?P<IDENTIFIER>[A-Za-z_][A-Za-z0-9_]*)`,
	`(?P<STRING_LITERAL>"(?:\\.|[^"\\])*")`,
	`(?P<FLOAT_LITERAL>[0-9]+\.[0-9]+)`,
	`(?P<INTEGER_LITERAL>[0-9]+)`,
	`(?P<NEWLINE>\n)`,
}, "|"))

var GROUP_NAMES = TOKEN_GRAMMAR.SubexpNames()
