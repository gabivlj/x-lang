package parser

import (
	"xlang/ast"
	"xlang/token"
)

func (p *Parser) parsePrefixExpression() ast.Expression {
	if p.curToken.Type == token.JUMP {
		return nil
	}
	expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	expression.SetLine(p.l.Line)
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}
