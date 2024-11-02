package lexer

import (
	"testing"

	"github.com/kacperkrolak/scene-description-language/token"
)

func TestNewToken(t *testing.T) {
	sampleContent := `MODIFY CAMERA {
    position: (0, 1.5, -10),
    rotation: (0, 0, 0),
    focalDistance: 35.0,
    ambientIntensity 0.3,
}

NUMBER pi = 3.14159265359,

COLOR red = (1.0, 0.0, 0.0),

MATERIAL shiny = {
    ambientIntensity: 0.1 * 5 + 1 - 1 / 2,
    color: red,
}

SPHERE sphere1 = {
    material: shiny,
    radius: 1.5,
}

LIGHT light1 AT (0, 1.5, 0) {
    color: red,
    specularIntensity: 0.5,
}`
	tokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.MODIFY, "MODIFY"},
		{token.CAMERA, "CAMERA"},
		{token.LBRACE, "{"},
		{token.IDENT, "position"},
		{token.COLON, ":"},
		{token.LPAREN, "("},
		{token.FLOAT, "0"},
		{token.COMMA, ","},
		{token.FLOAT, "1.5"},
		{token.COMMA, ","},
		{token.MINUS, "-"},
		{token.FLOAT, "10"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.IDENT, "rotation"},
		{token.COLON, ":"},
		{token.LPAREN, "("},
		{token.FLOAT, "0"},
		{token.COMMA, ","},
		{token.FLOAT, "0"},
		{token.COMMA, ","},
		{token.FLOAT, "0"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.IDENT, "focalDistance"},
		{token.COLON, ":"},
		{token.FLOAT, "35.0"},
		{token.COMMA, ","},
		{token.IDENT, "ambientIntensity"},
		{token.FLOAT, "0.3"},
		{token.COMMA, ","},
		{token.RBRACE, "}"},
		{token.NUMBER, "NUMBER"},
		{token.IDENT, "pi"},
		{token.ASSIGN, "="},
		{token.FLOAT, "3.14159265359"},
		{token.COMMA, ","},
		{token.COLOR, "COLOR"},
		{token.IDENT, "red"},
		{token.ASSIGN, "="},
		{token.LPAREN, "("},
		{token.FLOAT, "1.0"},
		{token.COMMA, ","},
		{token.FLOAT, "0.0"},
		{token.COMMA, ","},
		{token.FLOAT, "0.0"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.MATERIAL, "MATERIAL"},
		{token.IDENT, "shiny"},
		{token.ASSIGN, "="},
		{token.LBRACE, "{"},
		{token.IDENT, "ambientIntensity"},
		{token.COLON, ":"},
		{token.FLOAT, "0.1"},
		{token.MULTIPLY, "*"},
		{token.FLOAT, "5"},
		{token.PLUS, "+"},
		{token.FLOAT, "1"},
		{token.MINUS, "-"},
		{token.FLOAT, "1"},
		{token.DIVIDE, "/"},
		{token.FLOAT, "2"},
		{token.COMMA, ","},
		{token.IDENT, "color"},
		{token.COLON, ":"},
		{token.IDENT, "red"},
		{token.COMMA, ","},
		{token.RBRACE, "}"},
		{token.SPHERE, "SPHERE"},
		{token.IDENT, "sphere1"},
		{token.ASSIGN, "="},
		{token.LBRACE, "{"},
		{token.IDENT, "material"},
		{token.COLON, ":"},
		{token.IDENT, "shiny"},
		{token.COMMA, ","},
		{token.IDENT, "radius"},
		{token.COLON, ":"},
		{token.FLOAT, "1.5"},
		{token.COMMA, ","},
		{token.RBRACE, "}"},
		{token.LIGHT, "LIGHT"},
		{token.IDENT, "light1"},
		{token.AT, "AT"},
		{token.LPAREN, "("},
		{token.FLOAT, "0"},
		{token.COMMA, ","},
		{token.FLOAT, "1.5"},
		{token.COMMA, ","},
		{token.FLOAT, "0"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "color"},
		{token.COLON, ":"},
		{token.IDENT, "red"},
		{token.COMMA, ","},
		{token.IDENT, "specularIntensity"},
		{token.COLON, ":"},
		{token.FLOAT, "0.5"},
		{token.COMMA, ","},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(sampleContent)

	for i, tt := range tokens {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
