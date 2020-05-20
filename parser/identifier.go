package parser

import (
	"xlang/ast"
)

func (p *Parser) parseIdentifier() ast.Expression {
	i := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	i.SetLine(p.l.Line)
	return i
}
