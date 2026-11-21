package parser

import (
	"regexp"
	"strings"
)

var tokenGrammar = regexp.MustCompile(strings.Join([]string{
	`(?P<whitespace>[ \t\r]+)`,
	`(?P<comment>//[^\n]*)`,
	`(?P<k_function>\bfunction\b)`,
	`(?P<k_return>\breturn\b)`,
	`(?P<k_while>\bwhile\b)`,
	`(?P<k_print>\bprint\b)`,
	`(?P<k_if>\bif\b)`,
	`(?P<k_else>\belse\b)`,
	`(?P<k_true>\btrue\b)`,
	`(?P<k_false>\bfalse\b)`,
	`(?P<k_and>\band\b)`,
	`(?P<k_or>\bor\b)`,
	`(?P<k_not>\bnot\b)`,
	`(?P<t_semicolon>\;)`,
	`(?P<t_openBrace>\{)`,
	`(?P<t_closeBrace>\})`,
	`(?P<t_openParan>\()`,
	`(?P<t_closeParan>\))`,
	`(?P<t_identifier>[A-Za-z_][A-Za-z0-9_]*)`,
	`(?P<t_stringLit>"(?:\\.|[^"\\])*")`,
	`(?P<t_floatLit>[0-9]+\.[0-9]+)`,
	`(?P<t_integerLit>[0-9]+)`,
	`(?P<t_plus>\+)`,
	`(?P<t_minus>-)`,
	`(?P<t_star>\*)`,
	`(?P<t_slash>/)`,
	`(?P<t_mod>%)`,
	`(?P<t_bangEqual>!=)`,
	`(?P<t_equalEqual>==)`,
	`(?P<t_lesserEquals><=)`,
	`(?P<t_greaterEquals>>=)`,
	`(?P<t_lesser><)`,
	`(?P<t_greater>>)`,
	`(?P<t_equal>=)`,
	`(?P<newline>\n)`,
}, "|"))

var groupNames = tokenGrammar.SubexpNames()

type TokenType int

// this should be in order of tokenGrammer and invalid will be the default value if a token is invalid
const (
	invalid TokenType = iota
	whitespace
	comment
	k_function
	k_return
	k_while
	k_print
	k_if
	k_else
	k_true
	k_false
	k_and
	k_or
	k_not
	t_semicolon
	t_openBrace
	t_closeBrace
	t_openParen
	t_closeParen
	t_identifier
	t_stringLit
	t_floatLit
	t_integerLit
	t_plus
	t_minus
	t_star
	t_slash
	t_mod
	t_bangEqual
	t_equalEqual
	t_lesserEquals
	t_greaterEquals
	t_lesser
	t_greater
	t_equal
	newline
	k_eof
)

type Word struct {
	Raw    string
	Type   TokenType
	Line   uint
	Column uint
	End    uint
}

func Tokenize(source string) []Word {
	if len(groupNames) != int(k_eof) {
		panic("CORE_ERROR: Grammer contains unhandled groups")
	}
	tokens := make([]Word, 0)
	var line uint = 1
	var lineStart uint = 0
	for _, match := range tokenGrammar.FindAllStringSubmatchIndex(source, -1) {
		start, end := uint(match[0]), uint(match[1])
		raw := source[start:end]
		var tokenType TokenType
		for i := 1; i < len(match)/2; i++ {
			groupStart := match[i*2]
			if groupStart >= 0 {
				tokenType = TokenType(i)
				break
			}
		}
		switch tokenType {
		case comment, whitespace:
			continue
		case newline:
			line += 1
			lineStart = end
			continue
		default:
			tokens = append(tokens, Word{
				Raw:    raw,
				Type:   tokenType,
				Line:   line,
				Column: start - lineStart + 1,
				End:    end,
			})
		}
	}
	// tokens = append(tokens, Word{"EOF", K_EOF, line, 0, 0})
	return tokens
}
