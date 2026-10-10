package utils

import (
	"regexp"
)

type TokenType string

const (
    K_FUNCTION     TokenType = "K_FUNCTION"
    K_RETURN       TokenType = "K_RETURN"
    OPEN_BRACE     TokenType = "OPEN_BRACE"
    CLOSE_BRACE    TokenType = "CLOSE_BRACE"
    IDENTIFIER     TokenType = "IDENTIFIER"
    FLOAT          TokenType = "FLOAT"
    INTEGER        TokenType = "INTEGER"
    NEWLINE        TokenType = "NEWLINE"
    COMMENT        TokenType = "COMMENT"
    K_EOF          TokenType = "K_EOF"
)

func Tokenize(source string) []Word {
	tokens := make([]Word, 0)
	line := 1

	compiledGrammar := regexp.MustCompile(TOKEN_GRAMMAR)
	groupNames := compiledGrammar.SubexpNames()

	for _, match := range compiledGrammar.FindAllStringSubmatchIndex(source, -1) {
		start, end := match[0], match[1]
		raw := source[start:end]
		var tokenType TokenType
		for i := 1; i < len(match)/2; i++ {
			groupStart := match[i*2]
			if groupStart >= 0 {
				tokenType = TokenType(groupNames[i])
				break
			}
		}
		switch tokenType {
		case COMMENT:
			continue
		case NEWLINE:
			line += 1
			continue
		default:
			tokens = append(tokens, Word{
				Raw: raw,
				Type: tokenType, 
				Line: line,
				Start: start,
				End: end,
			})
		}
	}
	tokens = append(tokens, Word{"EOF", K_EOF, 0, 0, 0})
	return tokens
}
