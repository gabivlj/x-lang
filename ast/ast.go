package ast

import (
	"bytes"
	"fmt"
	"strings"
	"xlang/token"
)

// Node for the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement ...
type Statement interface {
	Node
	statementNode()
}

// Expression produces a value
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST
type Program struct {
	Statements []Statement
}

// TokenLiteral Root
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns the program as a strings
func (p *Program) String() string {
	var str strings.Builder = strings.Builder{}
	for _, s := range p.Statements {
		str.WriteString("Statement: " + s.String() + "\n")
	}
	return str.String()
}

// Identifier is a variable name
type Identifier struct {
	Token token.Token // IDENT
	Value string      // value
}

func (i *Identifier) expressionNode() {}

// TokenLiteral .
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }

// LetStatement represents things like let x = 5;
type LetStatement struct {
	Token token.Token // let
	Name  *Identifier // name
	Value Expression  // exp
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral .
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	out := strings.Builder{}

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents a return statement duh
type ReturnStatement struct {
	Token       token.Token // the 'return'
	ReturnValue Expression  // what is returning
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral .
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	out := strings.Builder{}
	out.WriteString("return ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents things like x + 10;
type ExpressionStatement struct {
	Token      token.Token // the first token of the exp.
	Expression Expression
}

func (exst *ExpressionStatement) statementNode() {}

// TokenLiteral .
func (exst *ExpressionStatement) TokenLiteral() string { return exst.Token.Literal }

func (exst *ExpressionStatement) String() string {
	if exst.Expression != nil {
		return exst.Expression.String()
	}
	return ""
}

// IntegerLiteral represents a number
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral ..
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

func (il *IntegerLiteral) String() string { return il.Token.Literal }

// PrefixExpression like -165 or !true
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral ..
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) String() string {
	out := strings.Builder{}

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents operations like 5 * 5
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) String() string {
	out := strings.Builder{}
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral returns token literal
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// Boolean value
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral is the literal string of the boolean
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// BlockStatement is a snippet of code inside a statement
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral returns the token literal of the token blockstatement
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	out := strings.Builder{}
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// IfExpression represents a if (<condition>) <consequence> else <consecuence>
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral returns the token literal
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IfExpression) String() string {
	out := strings.Builder{}
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" {\n")
	out.WriteString(ie.Consequence.String() + "\n} ")
	if ie.Alternative != nil {
		out.WriteString("else {\n")
		out.WriteString(ie.Alternative.String())
		out.WriteByte('\n')
		out.WriteByte('}')
		out.WriteByte('\n')
	}
	return out.String()
}

// FunctionLiteral fn(x, y, ...) {}
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral returns token literal
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FunctionLiteral) String() string {
	out := strings.Builder{}
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteByte('(')
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression represents a call(...)
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral .
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	out := strings.Builder{}
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteByte('(')
	out.WriteString(strings.Join(args, ", "))
	out.WriteByte(')')
	return out.String()
}

// StringLiteral represents a string
type StringLiteral struct {
	Value string
	Token token.Token
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral .
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// String .
func (sl *StringLiteral) String() string { return sl.Token.Literal }

// ArrayLiteral represents an array [...]
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral .
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

// String .
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// IndexExpression is left[right]
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral .
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

// String .
func (ie *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", ie.Left.String(), ie.Right.String())
}
