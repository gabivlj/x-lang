package parser

import (
	"xlang/ast"
	"xlang/token"
)

func (p *Parser) parseBoolean() ast.Expression {
	st := &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	return st
}
