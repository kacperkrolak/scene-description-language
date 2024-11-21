package evaluator

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kacperkrolak/scene-description-language/lexer"
	"github.com/kacperkrolak/scene-description-language/parser"
)

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"10", 10}}
	for _, tt := range tests {
		evaluator := NewEvaluator()
		evaluated := testEval(evaluator, tt.input)

		testNumberObject(t, evaluated, tt.expected)
	}
}

func testEval(evaluator *Evaluator, input string) Object {
	l := lexer.New(input)
	p := parser.New(l)
	ast := p.ParseFile()

	return evaluator.Eval(ast)
}

func testNumberObject(t *testing.T, obj Object, expected float64) bool {
	result, ok := obj.(*Number)
	if !ok {
		t.Errorf("object is not Number. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}
	return true
}

func testArrayObject(t *testing.T, obj Object, expected interface{}) bool {
	result, ok := obj.(*Array)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}

	var elementValues []interface{}
	for _, element := range result.Elements {
		switch element.(type) {
		case *Number:
			elementValues = append(elementValues, element.(*Number).Value)
		case *Array:
			t.Error("Testing nested arrays is not supported")
			return false
		default:
			t.Errorf("Unexpected type in array: %T", element)
			return false
		}
	}

	for i, element := range elementValues {
		if element != reflect.ValueOf(expected).Index(i).Interface() {
			t.Errorf("Array element %d has wrong value. got=%f, want=%f", i, element, reflect.ValueOf(expected).Index(i).Interface())
			return false
		}
	}

	return true
}

func TestEvalArrayExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected []float64
	}{
		{"[]", []float64{}},
		{"[1]", []float64{1}},
		{"[1, 2, 3]", []float64{1, 2, 3}},
	}

	for _, tt := range tests {
		evaluator := NewEvaluator()
		evaluated := testEval(evaluator, tt.input)

		if isError(evaluated) {
			t.Errorf("error: %v", evaluated)
			continue
		}

		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestEvalPropertiesExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]float64
	}{
		{"{}", map[string]float64{}},
		{"{a: 1,}", map[string]float64{"a": 1}},
		{"{a: 1, b: 2}", map[string]float64{"a": 1, "b": 2}},
	}

	for _, tt := range tests {
		evaluator := NewEvaluator()
		evaluated := testEval(evaluator, tt.input)

		result, ok := evaluated.(*Dictionary)
		if !ok {
			t.Errorf("object is not Dictionary. got=%T (%+v)", evaluated, evaluated)
			return
		}

		if len(result.Properties) != len(tt.expected) {
			t.Errorf("wrong number of properties. got=%d, want=%d", len(result.Properties), len(tt.expected))
			return
		}

		for expectedKey, expectedValue := range tt.expected {
			value, ok := result.Properties[expectedKey]
			if !ok {
				t.Errorf("no key in properties. got=%v", result.Properties)
				return
			}

			testNumberObject(t, value, expectedValue)
		}
	}
}

func TestEvalPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"-5", -5},
		{"--5", 5},
	}

	for _, tt := range tests {
		evaluator := NewEvaluator()
		evaluated := testEval(evaluator, tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5-5", 0},
		{"5+5", 10},
		{"5*5", 25},
		{"5/5", 1},
		{"5-10", -5},
	}

	for _, tt := range tests {
		evaluator := NewEvaluator()
		evaluated := testEval(evaluator, tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5+5*2", 15},
		{"-5+5", 0},
		{"5*2-1", 9},
		{"5*(2-1)", 5},
		{"-5*(2-1)", -5},
		{"5+5/2", 7.5},
		{"-5+5/2", -2.5},
	}

	for _, tt := range tests {
		evaluator := NewEvaluator()
		evaluated := testEval(evaluator, tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func testEnvironmentObject(t *testing.T, env Environment, expected map[string]Entity) bool {
	for expectedKey, expectedValue := range expected {
		evaluatedValue, ok := env.store[expectedKey]
		if !ok {
			t.Errorf("key not found in environment: %s, env=%v", expectedKey, env.store)
			return false
		}

		if !reflect.DeepEqual(evaluatedValue, expectedValue) {
			gotJson, err := json.MarshalIndent(evaluatedValue, "", "  ")
			if err != nil {
				t.Errorf("error marshalling got: %v", err)
				return false
			}

			wantJson, err := json.MarshalIndent(expectedValue, "", "  ")
			if err != nil {
				t.Errorf("error marshalling want: %v", err)
				return false
			}

			t.Errorf("wrong value for key %s. got=%s, want=%s", expectedKey, string(gotJson), string(wantJson))
			return false
		}
	}

	return true
}

func TestEvalAssignStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]Entity
	}{
		{"NUMBER a = 5", map[string]Entity{"a": Entity{Class: "NUMBER", Value: &Number{Value: 5}}}},
		{"COLOR a = [1, 2, 3]", map[string]Entity{"a": Entity{Class: "COLOR", Value: &Array{Elements: []Object{&Number{Value: 1}, &Number{Value: 2}, &Number{Value: 3}}}}}},
		{`
NUMBER a = 5
NUMBER b = -a
COLOR c = [a, a, b]
`, map[string]Entity{
			"a": Entity{Class: "NUMBER", Value: &Number{Value: 5}},
			"b": Entity{Class: "NUMBER", Value: &Number{Value: -5}},
			"c": Entity{Class: "COLOR", Value: &Array{Elements: []Object{&Number{Value: 5}, &Number{Value: 5}, &Number{Value: -5}}}},
		}},
	}

	for _, tt := range tests {
		evaluator := NewEvaluator()
		evaluated := testEval(evaluator, tt.input)
		if isError(evaluated) {
			t.Errorf("error: %v", evaluated)
			continue
		}

		testEnvironmentObject(t, *evaluator.env, tt.expected)
	}
}

func TestExportingValues(t *testing.T) {
	input := `
NUMBER r = 255
NUMBER g = 0
NUMBER b = 0
COLOR red = [r, g, b]
SPHERE sphere = {color: red, radius: 1}
`
	evaluator := NewEvaluator()
	evaluated := testEval(evaluator, input)
	if isError(evaluated) {
		t.Errorf("error: %v", evaluated)
		return
	}

	exported := evaluator.ExportValues()

	expected := map[string][]Entity{
		"NUMBER": []Entity{
			Entity{Class: "NUMBER", Value: &Number{Value: 255}},
			Entity{Class: "NUMBER", Value: &Number{Value: 0}},
			Entity{Class: "NUMBER", Value: &Number{Value: 0}},
		},
		"COLOR": []Entity{
			Entity{Class: "COLOR", Value: &Array{Elements: []Object{
				&Number{Value: 255},
				&Number{Value: 0},
				&Number{Value: 0},
			}}},
		},
		"SPHERE": []Entity{
			Entity{Class: "SPHERE", Value: &Dictionary{Properties: map[string]Object{
				"color": &Array{Elements: []Object{
					&Number{Value: 255},
					&Number{Value: 0},
					&Number{Value: 0},
				}},
				"radius": &Number{Value: 1},
			}}},
		},
	}

	for expectedKey, expectedValue := range expected {
		if !reflect.DeepEqual(exported.Entities[expectedKey], expectedValue) {
			t.Errorf("wrong value for key %s. got=%v, want=%v", expectedKey, exported.Entities[expectedKey], expectedValue)
		}
	}
}
