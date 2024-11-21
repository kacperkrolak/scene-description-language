package evaluator

import (
	"fmt"
	"io"
	"strings"

	"github.com/kacperkrolak/scene-description-language/ast"
	"github.com/kacperkrolak/scene-description-language/lexer"
	"github.com/kacperkrolak/scene-description-language/parser"
)

type ObjectType string

const (
	NUMBER_OBJ     ObjectType = "NUMBER"
	COLOR_OBJ      ObjectType = "COLOR"
	MATERIAL_OBJ   ObjectType = "MATERIAL"
	ERROR_OBJ      ObjectType = "ERROR"
	ARRAY_OBJ      ObjectType = "ARRAY"
	DICTIONARY_OBJ ObjectType = "DICTIONARY"
	ENTITY_OBJ     ObjectType = "ENTITY"
)

// Object represent a constant defined in the SDL file.
type Object interface {
	Type() ObjectType
}

// Array represents an array of objects.
type Array struct {
	Elements []Object
}

func (a Array) Type() ObjectType {
	return ARRAY_OBJ
}

// Array represents an array of objects.
type Dictionary struct {
	Properties map[string]Object
}

func (d Dictionary) Type() ObjectType {
	return DICTIONARY_OBJ
}

// Number represents a number constant defined in the SDL file.
type Number struct {
	Value float64
}

func (n Number) Type() ObjectType {
	return NUMBER_OBJ
}

// Entity represents a scene object (like Sphere, Light) with its properties
type Entity struct {
	Class string
	Value Object
}

func (e Entity) Type() ObjectType {
	return ENTITY_OBJ
}

// Error represents an object that could not be evaluated.
type Error struct {
	Message string
}

func (e Error) Type() ObjectType {
	return ERROR_OBJ
}

// Environment represents the environment in which the SDL file is evaluated.
// It contains the contants defined in the SDL file.
type Environment struct {
	store map[string]Entity
}

// EvaluatedValues groups the most high-level objects that can be evaluated.
type EvaluatedValues struct {
	Entities map[string][]Entity // Entities grouped by their class name.
}

type Evaluator struct {
	env *Environment
}

func Eval(node ast.Node) (EvaluatedValues, error) {
	evaluator := NewEvaluator()
	ev := evaluator.Eval(node)
	if isError(ev) {
		return EvaluatedValues{}, fmt.Errorf("failed to evaluate file: %s", ev.(Error).Message)
	}

	return evaluator.ExportValues(), nil
}

func NewEvaluator() *Evaluator {
	return &Evaluator{env: &Environment{store: make(map[string]Entity)}}
}

func (evaluator *Evaluator) ExportValues() EvaluatedValues {
	entities := make(map[string][]Entity)
	for _, entity := range evaluator.env.store {
		entities[entity.Class] = append(entities[entity.Class], entity)
	}

	return EvaluatedValues{Entities: entities}
}

func (evaluator *Evaluator) EvaluateFile(r io.Reader) error {
	fileAst, err := getAst(r)
	if err != nil {
		return err
	}

	obj := evaluator.Eval(fileAst)
	if isError(obj) {
		return fmt.Errorf("failed to evaluate file: %s", obj.(Error).Message)
	}

	return nil
}

func isError(obj Object) bool {
	return obj.Type() == ERROR_OBJ
}

func (evaluator *Evaluator) Eval(node ast.Node) Object {
	switch node := node.(type) {
	case *ast.File:
		return evaluator.evalFile(node)
	case *ast.AssignStatement:
		return evaluator.evalAssignStatement(node)
	case *ast.ExpressionStatement:
		return evaluator.Eval(node.Expression)
	case *ast.FloatLiteral:
		return &Number{Value: node.Value}
	case *ast.ArrayExpression:
		return evaluator.evalArrayExpression(node)
	case *ast.PropertiesExpression:
		return evaluator.evalPropertiesExpression(node)
	case *ast.PrefixExpression:
		right := evaluator.Eval(node.Right)
		if isError(right) {
			return right
		}

		return evaluator.evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := evaluator.Eval(node.Left)
		if isError(left) {
			return left
		}

		right := evaluator.Eval(node.Right)
		if isError(right) {
			return right
		}

		return evaluator.evalInfixExpression(node.Operator, left, right)
	case *ast.Identifier:
		return evaluator.evalIdentifier(node)
	default:
		return Error{Message: fmt.Sprintf("unknown node type: %T", node)}
	}
}

func (evaluator *Evaluator) evalFile(file *ast.File) Object {
	var result Object
	for _, statement := range file.Statements {
		result = evaluator.Eval(statement)
		if isError(result) {
			return result
		}
	}
	return result
}

func (evaluator *Evaluator) evalArrayExpression(node *ast.ArrayExpression) Object {
	var elements []Object
	for _, element := range node.Elements {
		result := evaluator.Eval(element)
		if isError(result) {
			return result
		}

		elements = append(elements, result)
	}
	return &Array{Elements: elements}
}

func (evaluator *Evaluator) evalPropertiesExpression(node *ast.PropertiesExpression) Object {
	properties := make(map[string]Object)

	for key, element := range node.Properties {
		result := evaluator.Eval(element)
		if isError(result) {
			return result
		}

		properties[key] = result
	}

	return &Dictionary{Properties: properties}
}

// GetAst uses scene-description-language module to parse
// the configuration file into an AST.
func getAst(r io.Reader) (*ast.File, error) {
	configString, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	sdlLexer := lexer.New(string(configString))
	sdlParser := parser.New(sdlLexer)

	ast := sdlParser.ParseFile()
	if len(sdlParser.Errors()) > 0 {
		return nil, fmt.Errorf("failed to parse config: %s", strings.Join(sdlParser.Errors(), ";"))
	}

	return ast, nil
}

func (evaluator *Evaluator) evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "-":
		return evalMinusOperator(right)
	default:
		return Error{Message: fmt.Sprintf("unknown operator: %s", operator)}
	}
}

func (evaluator *Evaluator) evalInfixExpression(operator string, left Object, right Object) Object {
	leftNumber, ok := left.(*Number)
	if !ok {
		return Error{Message: fmt.Sprintf("Infix operator only supports numbers, got: %s", left.Type())}
	}

	rightNumber, ok := right.(*Number)
	if !ok {
		return Error{Message: fmt.Sprintf("Infix operator only supports numbers, got: %s", right.Type())}
	}

	switch operator {
	case "-":
		return &Number{Value: leftNumber.Value - rightNumber.Value}
	case "+":
		return &Number{Value: leftNumber.Value + rightNumber.Value}
	case "*":
		return &Number{Value: leftNumber.Value * rightNumber.Value}
	case "/":
		if rightNumber.Value == 0 {
			return Error{Message: "division by zero"}
		}

		return &Number{Value: leftNumber.Value / rightNumber.Value}
	default:
		return Error{Message: fmt.Sprintf("unknown operator: %s", operator)}
	}
}

func evalMinusOperator(right Object) Object {
	number, ok := right.(*Number)
	if !ok {
		return Error{Message: fmt.Sprintf("unknown operator: -%s", right.Type())}
	}

	return &Number{Value: -number.Value}
}

func (evaluator *Evaluator) evalAssignStatement(s *ast.AssignStatement) Object {
	// Don't allow redefining objects.
	if _, ok := evaluator.env.store[s.Name.Value]; ok {
		return Error{Message: fmt.Sprintf("redefining objects is not allowed: %s", s.Name.Value)}
	}

	evaluatedValue := evaluator.Eval(s.Value)
	if isError(evaluatedValue) {
		return evaluatedValue
	}

	evaluatedEntity := Entity{Class: s.Token.Literal, Value: evaluatedValue}

	evaluator.env.store[s.Name.Value] = evaluatedEntity

	return evaluatedValue
}

func (evaluator *Evaluator) evalIdentifier(node *ast.Identifier) Object {
	entity, ok := evaluator.env.store[node.Value]
	if !ok {
		return Error{Message: fmt.Sprintf("undefined identifier: %s", node.Value)}
	}

	return entity.Value
}
