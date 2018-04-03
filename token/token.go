package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL Type = "ILLEGAL"
	EOF          = "EOF"

	SLASH     = "/"
	COMMA     = ","
	DOLLAR    = "$"
	BACKSLASH = "\\"
	INT       = "INT"
	LIT       = "LITERAL"
	CMD       = "COMMAND"
	DIV       = "DIVIDER"
	IDENT     = "IDENTIFIER"
	LBRACE    = "{"
	RBRACE    = "}"
	EXPLMARK  = "!"

	SEMICOLON = ";"
	NEWLINE   = "\\n"
	LPAREN    = "("
	RPAREN    = ")"
	QUESTION  = "?"
	PERIOD    = "."
	COLON     = ":"
)

// If the program was to have any key words, like
// for or if, we can use this map to figure out what
// type of token it is
var keywords = map[string]Type{}

// LookupIdent takes a string and returns its corresponding type
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return COLON
}
