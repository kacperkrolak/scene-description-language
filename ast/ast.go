package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kacperkrolak/scene-description-language/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type File struct {
	Statements []Statement
}

func (p *File) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *File) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String() + "\n")
	}
	return out.String()
}

type AssignStatement struct {
	Token token.Token // the type token
	Name  *Identifier
	Value Expression
}

func (ls *AssignStatement) statementNode() {

}
func (ls *AssignStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	return out.String()
}

type ModifyStatement struct {
	Token token.Token // the type token
	Name  *Identifier
	Value Expression
}

func (ms *ModifyStatement) statementNode() {

}
func (ms *ModifyStatement) TokenLiteral() string {
	return ms.Token.Literal
}

func (ms *ModifyStatement) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("%s %s ", ms.Token.Literal, ms.Name.String()))
	if ms.Value != nil {
		out.WriteString(ms.Value.String())
	}
	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string { return i.Value }

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type FloatLiteral struct {
	Token token.Token // the token.FLOAT token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

type ArrayExpression struct {
	Token    token.Token // The token.ARRAY token
	Elements []Expression
}

func (ae *ArrayExpression) expressionNode()      {}
func (ae *ArrayExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *ArrayExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ae.Elements {
		args = append(args, a.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString("]")
	return out.String()
}

type PropertiesExpression struct {
	Token      token.Token // The token.PROPERTIES token
	Properties map[string]Expression
}

func (pe *PropertiesExpression) expressionNode()      {}
func (pe *PropertiesExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PropertiesExpression) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for k, v := range pe.Properties {
		out.WriteString(fmt.Sprintf("%s: %s,\n", k, v.String()))
	}
	out.WriteString("}")
	return out.String()
}
