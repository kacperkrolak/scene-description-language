package ast

import (
	"strings"
	"testing"

	"github.com/kacperkrolak/scene-description-language/token"
)

func TestString(t *testing.T) {
	file := &File{
		Statements: []Statement{
			&AssignStatement{
				Token: token.Token{Type: token.COLOR, Literal: "COLOR"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "red"},
					Value: "red",
				},
				Value: &ArrayExpression{
					Token: token.Token{Type: token.LBRACKET, Literal: "["},
					Elements: []Expression{
						&FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "1.0"}, Value: 1.0},
						&FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0.0"}, Value: 0.0},
						&FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0.0"}, Value: 0.0},
					},
				},
			},
			&AssignStatement{
				Token: token.Token{Type: token.COLOR, Literal: "NUMBER"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "pi"},
					Value: "pi",
				},
				Value: &FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "3.14159265359"}, Value: 3.14159265359},
			},
			&AssignStatement{
				Token: token.Token{Type: token.COLOR, Literal: "NUMBER"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "negative"},
					Value: "negative",
				},
				Value: &PrefixExpression{
					Token:    token.Token{Type: token.MINUS, Literal: "-"},
					Operator: "-",
					Right:    &FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "1.0"}, Value: 1.0},
				},
			},
		},
	}

	expectedContent := `COLOR red = [1.0, 0.0, 0.0]
NUMBER pi = 3.14159265359
NUMBER negative = (-1.0)
`

	got := file.String()
	if got != expectedContent {
		t.Errorf("file.String() wrong. got=%q", got)
	}
}

func TestPropertiesStatement(t *testing.T) {
	file := &File{
		Statements: []Statement{
			&AssignStatement{
				Token: token.Token{Type: token.COLOR, Literal: "MATERIAL"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "shiny"},
					Value: "shiny",
				},
				Value: &PropertiesExpression{
					Token: token.Token{Type: token.PROPERTIES, Literal: "{"},
					Properties: map[string]Expression{
						"ambientIntensity": &FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0.1"}, Value: 0.1},
						"color": &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "white"},
							Value: "white",
						},
						"diffuseIntensity":  &FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0.7"}, Value: 0.7},
						"specularIntensity": &FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "1.0"}, Value: 1.0},
					},
				},
			},
		},
	}

	expectedLines := []string{
		"MATERIAL shiny = {",
		"ambientIntensity: 0.1,",
		"color: white,",
		"diffuseIntensity: 0.7,",
		"specularIntensity: 1.0,",
		"}",
	}

	// Split the output into lines and trim spaces.
	output := file.String()
	outputLines := strings.Split(output, "\n")
	for i := range outputLines {
		outputLines[i] = strings.TrimSpace(outputLines[i])
	}

	// Create a map for quick lookup of lines.
	outputMap := make(map[string]bool)
	for _, line := range outputLines {
		outputMap[line] = true
	}

	// Check that each expected line is present.
	for _, expectedLine := range expectedLines {
		if !outputMap[expectedLine] {
			t.Errorf("Expected line '%s' not found in output: %q", expectedLine, output)
		}
	}
}

func TestModifyStatement(t *testing.T) {
	file := &File{
		Statements: []Statement{
			&ModifyStatement{
				Token: token.Token{Type: token.MODIFY, Literal: "MODIFY"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "CAMERA"},
					Value: "CAMERA",
				},
				Value: &PropertiesExpression{
					Token: token.Token{Type: token.PROPERTIES, Literal: "{"},
					Properties: map[string]Expression{
						"position": &ArrayExpression{
							Token: token.Token{Type: token.LBRACKET, Literal: "["},
							Elements: []Expression{
								&FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0.0"}, Value: 0.0},
								&FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0.0"}, Value: 0.0},
								&FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0.0"}, Value: 0.0},
							},
						},
					},
				},
			},
		},
	}

	expectedContent := `MODIFY CAMERA {
position: [0.0, 0.0, 0.0],
}
`

	got := file.String()
	if got != expectedContent {
		t.Errorf("file.String() wrong. got=%q", got)
	}
}
