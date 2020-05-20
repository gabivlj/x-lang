package parser

import "xlang/ast"

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	exp.SetLine(p.l.Line)
	precedence := p.currPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}
