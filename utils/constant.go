package utils

import "strings"

var TOKEN_GRAMMAR = strings.Join([]string{
	`(?P<COMMENT>//[^\n]*)`,
	`(?P<K_FUNCTION>\bfunction\b)`,
	`(?P<K_RETURN>\breturn\b)`,
	`(?P<OPEN_BRACE>\{)`,
	`(?P<CLOSE_BRACE>\})`,
	`(?P<IDENTIFIER>[A-Za-z_][A-Za-z0-9_]*)`,
	`(?P<FLOAT>[0-9]+\.[0-9]+)`,
	`(?P<INTEGER>[0-9]+)`,
	`(?P<NEWLINE>\n)`,
}, "|")
