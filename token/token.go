package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

var keywords = map[string]TokenType{
	"MODIFY":   MODIFY,
	"CAMERA":   CAMERA,
	"PLACE":    PLACE,
	"AT":       AT,
	"NUMBER":   NUMBER,
	"COLOR":    COLOR,
	"MATERIAL": MATERIAL,
	"SPHERE":   SPHERE,
	"LIGHT":    LIGHT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifiers + literals
	IDENT      = "IDENT"      // x, y, sphere1, light_blue ...
	PROPERTIES = "PROPERTIES" // {x: 1, y: 2, z: 3}
	FLOAT      = "FLOAT"

	// Operators
	ASSIGN   = "="
	MINUS    = "-"
	PLUS     = "+"
	MULTIPLY = "*"
	DIVIDE   = "/"

	// Delimiters
	COMMA    = ","
	COLON    = ":"
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	MODIFY = "MODIFY"
	CAMERA = "CAMERA"
	PLACE  = "PLACE"
	AT     = "AT"

	// Object types
	NUMBER   = "NUMBER"
	COLOR    = "COLOR"
	MATERIAL = "MATERIAL"
	SPHERE   = "SPHERE"
	LIGHT    = "LIGHT"

	// Special token for statements that don't need a token
	NONE = "NONE"
)
