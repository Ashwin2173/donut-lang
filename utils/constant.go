package utils

import (
	"regexp"
	"strings"
)

const (
	FILE_EXTENSION   = ".dt"
	TARGET_EXTENSION = ".class"
)

type TypeId int32
type TokenType int
type ExpOpTypeId int
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
	MOD            TokenType = iota
	BANG_EQUAL     TokenType = iota
	EQUAL_EQUAL    TokenType = iota
	LESSER_EQUALS  TokenType = iota
	GREATER_EQUALS TokenType = iota
	LESSER         TokenType = iota
	GREATER        TokenType = iota
	EQUAL          TokenType = iota
	NEWLINE        TokenType = iota
	K_EOF          TokenType = iota
)

var TOKEN_GRAMMAR = regexp.MustCompile(strings.Join([]string{
	`(?P<COMMENT>//[^\n]*)`,
	`(?P<KEYWORD_FUNCTION>\bfunction\b)`,
	`(?P<KEYWORD_RETURN>\breturn\b)`,
	`(?P<KEYWORD_PRINT>\bprint\b)`,
	`(?P<KEYWORD_TRUE>\btrue\b)`,
	`(?P<KEYWORD_FALSE>\bfalse\b)`,
	`(?P<KEYWORD_AND>\band\b)`,
	`(?P<KEYWORD_OR>\bor\b)`,
	`(?P<KEYWORD_NOT>\bnot\b)`,
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
	`(?P<MOD>%)`,
	`(?P<BANG_EQUAL>!=)`,
	`(?P<EQUAL_EQUAL>==)`,
	`(?P<LESSER_EQUALS><=)`,
	`(?P<GREATER_EQUALS>>=)`,
	`(?P<LESSER><)`,
	`(?P<GREATER>>)`,
	`(?P<EQUAL>=)`,
	`(?P<NEWLINE>\n)`,
}, "|"))

var GROUP_NAMES = TOKEN_GRAMMAR.SubexpNames()

// primary types
const (
	T_UNKNOWN TypeId = iota
	T_INT     TypeId = iota
	T_FLOAT   TypeId = iota
	T_STRING  TypeId = iota
	T_BOOLEAN TypeId = iota
)

// expression operation
const (
	UOP_INEG ExpOpTypeId = iota
	UOP_FNEG ExpOpTypeId = iota
	UOP_NOT  ExpOpTypeId = iota

	BOP_NONE ExpOpTypeId = iota

	BOP_IADD ExpOpTypeId = iota
	BOP_ISUB ExpOpTypeId = iota
	BOP_IMUL ExpOpTypeId = iota
	BOP_IDIV ExpOpTypeId = iota
	BOP_IREM ExpOpTypeId = iota

	BOP_FADD ExpOpTypeId = iota
	BOP_FSUB ExpOpTypeId = iota
	BOP_FMUL ExpOpTypeId = iota
	BOP_FDIV ExpOpTypeId = iota
	BOP_FREM ExpOpTypeId = iota

	BOP_I2F ExpOpTypeId = iota

	BOP_ICMP_EQ ExpOpTypeId = iota
	BOP_ICMP_NE ExpOpTypeId = iota
	BOP_ICMP_LT ExpOpTypeId = iota
	BOP_ICMP_LE ExpOpTypeId = iota
	BOP_ICMP_GT ExpOpTypeId = iota
	BOP_ICMP_GE ExpOpTypeId = iota

	BOP_FCMP_EQ ExpOpTypeId = iota
	BOP_FCMP_NE ExpOpTypeId = iota
	BOP_FCMP_LT ExpOpTypeId = iota
	BOP_FCMP_LE ExpOpTypeId = iota
	BOP_FCMP_GT ExpOpTypeId = iota
	BOP_FCMP_GE ExpOpTypeId = iota

	BOP_IAND ExpOpTypeId = iota
	BOP_IOR  ExpOpTypeId = iota

	BOP_S_CONCAT ExpOpTypeId = iota
	BOP_S_EQ     ExpOpTypeId = iota
	BOP_S_NE     ExpOpTypeId = iota
)

const (
	OP_ICONST_0        = 0x03
	OP_ICONST_1        = 0x04
	OP_LOAD_CONSTANT   = 0x12
	OP_LOAD_CONSTANT_W = 0x13
	OP_IAND            = 0x7E
	OP_IOR             = 0x80
	OP_IXOR            = 0x82
	OP_I2F             = 0x86
	OP_FCMPL           = 0x95
	OP_IFEQ            = 0x99
	OP_IFNE            = 0x9A
	OP_IFLT            = 0x9B
	OP_IFGE            = 0x9C
	OP_IFGT            = 0x9D
	OP_IFLE            = 0x9E
	OP_IF_ICMPEQ       = 0x9F
	OP_IF_ICMPNE       = 0xA0
	OP_IF_ICMPLT       = 0xA1
	OP_IF_ICMPGE       = 0xA2
	OP_IF_ICMPGT       = 0xA3
	OP_IF_ICMPLE       = 0xA4
	OP_GOTO            = 0xA7
	OP_RETURN          = 0xB1
	OP_INVOKE_VIRTUAL  = 0xb6
	OP_INVOKE_STATIC   = 0xB8
	OP_INVOKE_DYNAMIC  = 0xBA
	OP_IADD            = 0x60
	OP_ISUB            = 0x64
	OP_IMUL            = 0x68
	OP_IDIV            = 0x6C
	OP_IREM            = 0x70
	OP_FADD            = 0x62
	OP_FSUB            = 0x66
	OP_FMUL            = 0x6A
	OP_FDIV            = 0x6E
	OP_FREM            = 0x72
	OP_FNEG            = 0x76
	OP_INEG            = 0x74
	OP_BIPUSH          = 0x10
	OP_ISTORE_0        = 0x3B
	OP_ISTORE_1        = 0x3C
	OP_ISTORE_2        = 0x3D
	OP_ISTORE_3        = 0x3E
	OP_ISTORE          = 0x36
	OP_ILOAD_0         = 0x1A
	OP_ILOAD_1         = 0x1B
	OP_ILOAD_2         = 0x1C
	OP_ILOAD_3         = 0x1D
	OP_ILOAD           = 0x15
	OP_FSTORE_0        = 0x43
	OP_FSTORE_1        = 0x44
	OP_FSTORE_2        = 0x45
	OP_FSTORE_3        = 0x46
	OP_FSTORE          = 0x38
	OP_FLOAD_0         = 0x22
	OP_FLOAD_1         = 0x23
	OP_FLOAD_2         = 0x24
	OP_FLOAD_3         = 0x25
	OP_FLOAD           = 0x17
	OP_ASTORE_0        = 0x4B
	OP_ASTORE_1        = 0x4C
	OP_ASTORE_2        = 0x4D
	OP_ASTORE_3        = 0x4E
	OP_ASTORE          = 0x3A
	OP_ALOAD_0         = 0x2A
	OP_ALOAD_1         = 0x2B
	OP_ALOAD_2         = 0x2C
	OP_ALOAD_3         = 0x2D
	OP_ALOAD           = 0x19
	OP_WIDE            = 0xC4
)
