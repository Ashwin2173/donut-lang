package utils

func Tokenize(source string) []Word {
	tokens := make([]Word, 0)
	var line uint = 1
	if len(GROUP_NAMES) != int(K_EOF) {
		panic("CORE_ERROR: Grammer contains unhandled groups")
	}
	for _, match := range TOKEN_GRAMMAR.FindAllStringSubmatchIndex(source, -1) {
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
