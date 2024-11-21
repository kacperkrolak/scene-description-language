package parser

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/kacperkrolak/scene-description-language/ast"
	"github.com/kacperkrolak/scene-description-language/lexer"
	"github.com/kacperkrolak/scene-description-language/token"
)

func TestAssignStatements(t *testing.T) {
	input := `
   NUMBER pi = 3.14159265359
   NUMBER red = [1.0, 0.0, 0.0]
   NUMBER negative = pi
   `
	l := lexer.New(input)
	p := New(l)
	program := p.ParseFile()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseFile() returned nil")
	}

	tests := []struct {
		expectedIdentifier string
		expectedValue      ast.Expression
	}{
		{"pi", &ast.FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "3.14159265359"}, Value: 3.14159265359}},
		{"red", &ast.ArrayExpression{
			Token: token.Token{Type: token.LBRACKET, Literal: "["},
			Elements: []ast.Expression{
				&ast.FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "1.0"}, Value: 1.0},
				&ast.FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0.0"}, Value: 0.0},
				&ast.FloatLiteral{Token: token.Token{Type: token.FLOAT, Literal: "0.0"}, Value: 0.0},
			}}},
		{"negative", &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "pi"}, Value: "pi"}},
	}

	if len(program.Statements) != len(tests) {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", len(tests), len(program.Statements))
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		assignStmt, ok := stmt.(*ast.AssignStatement)
		if !ok {
			t.Errorf("stmt is not *ast.AssignStatement. got=%T", stmt)
			continue
		}

		if assignStmt.Value.String() != tt.expectedValue.String() {
			t.Errorf("assignStmt.Value.String() not %s. got=%s", tt.expectedValue.String(), assignStmt.Value.String())
		}
	}
}

func testAssignStatement(t *testing.T, s ast.Statement, name string, expectedValue ast.Expression) bool {
	if !isObjectType(token.LookupIdent(s.TokenLiteral())) {
		t.Errorf("s.TokenLiteral not an object type. got=%q", s.TokenLiteral())
		return false
	}

	assignStmt, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}
	if assignStmt.Name.Value != name {
		t.Errorf("assignStmt.Name.Value not '%s'. got=%s", name, assignStmt.Name.Value)
		return false
	}
	if assignStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, assignStmt.Name)
		return false
	}

	if !reflect.DeepEqual(assignStmt.Value, expectedValue) {
		t.Errorf("assignStmt.Value not %v. got=%v", expectedValue, assignStmt.Value)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseFile()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input      string
		operator   string
		floatValue float64
	}{
		{"-15", "-", 15},
		{"-5.5", "-", 5.5},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		file := p.ParseFile()
		checkParserErrors(t, p)
		if len(file.Statements) != 1 {
			var statements []string
			for _, s := range file.Statements {
				statements = append(statements, "\""+s.String()+"\"")
			}
			t.Fatalf("file.Statements does not contain %d statements. got=%d: %v\n",
				1, len(file.Statements), statements)
		}
		stmt, ok := file.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("file.Statements[0] is not ast.ExpressionStatement. got=%T",
				file.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testFloatLiteral(t, exp.Right, tt.floatValue) {
			return
		}
	}
}

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
	number, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("fl not *ast.FloatLiteral. got=%T", fl)
		return false
	}
	if number.Value != value {
		t.Errorf("number.Value not %f. got=%f", value, number.Value)
		return false
	}

	// It is difficult to compare float64 values as strings, because of the precision
	// so we will convert the Literal to a float64 and compare that.
	floatTokenLiteral, err := strconv.ParseFloat(number.TokenLiteral(), 64)
	if err != nil {
		t.Errorf("could not parse %q as float64", number.TokenLiteral())
		return false
	}

	if floatTokenLiteral != value {
		t.Errorf("number.TokenLiteral not %f. got=%s", value, number.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", float64(5), "+", float64(5)},
		{"5 - 5", float64(5), "-", float64(5)},
		{"5 * 5", float64(5), "*", float64(5)},
		{"5 / 5", float64(5), "/", float64(5)},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseFile()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)\n",
		},
		{
			"-a + b * c / d",
			"((-a) + ((b * c) / d))\n",
		},
		{
			"(-a + b) * (c / d)",
			"(((-a) + b) * (c / d))\n",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseFile()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T,
	exp ast.Expression,
	expected interface{}) bool {
	switch v := expected.(type) {
	case float64:
		return testFloatLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestArrayExpression(t *testing.T) {
	tests := []struct {
		input            string
		expectedElements []float64
	}{
		{"[]", []float64{}},
		{"[1.0]", []float64{1.0}},
		{`
[
  1.0,
  2.0,
  3.0,
]`, []float64{1.0, 2.0, 3.0}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseFile()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		array, ok := stmt.Expression.(*ast.ArrayExpression)
		if !ok {
			t.Fatalf("exp not *ast.ArrayExpression. got=%T", stmt.Expression)
		}

		if len(array.Elements) != len(tt.expectedElements) {
			t.Fatalf("len(array.Elements) not %d. got=%d", len(tt.expectedElements), len(array.Elements))
		}

		for i, element := range tt.expectedElements {
			if !testLiteralExpression(t, array.Elements[i], element) {
				return
			}
		}
	}
}

func TestParsingPropertiesExpression(t *testing.T) {
	type expectation struct {
		key   string
		value float64
	}

	tests := []struct {
		input    string
		expected []expectation
	}{
		{
			input: `{
prop_one: 1.0,
prop_two: 2.0,
}`,
			expected: []expectation{
				{"prop_one", 1.0},
				{"prop_two", 2.0},
			},
		},
		{
			input: `{prop_one: 1.0, prop_two: 2.0}`,
			expected: []expectation{
				{"prop_one", 1.0},
				{"prop_two", 2.0},
			},
		},
		{
			input:    `{}`,
			expected: []expectation{},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseFile()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseFile() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		propertiesExp, ok := stmt.Expression.(*ast.PropertiesExpression)
		if !ok {
			t.Fatalf("exp not *ast.PropertiesExpression. got=%T", stmt.Expression)
		}

		for _, expected := range tt.expected {
			prop, ok := propertiesExp.Properties[expected.key]
			if !ok {
				t.Errorf("propertiesExp.Properties does not contain key '%s'", expected.key)
				continue
			}

			if !testLiteralExpression(t, prop, expected.value) {
				return
			}
		}
	}
}
