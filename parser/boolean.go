package parser

import (
	"xlang/ast"
	"xlang/token"
)

func (p *Parser) parseBoolean() ast.Expression {
	st := &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	st.SetLine(p.curToken.Line)
	return st
}
