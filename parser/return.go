package parser

import (
	"xlang/ast"
	"xlang/token"
)

func (p *Parser) parseReturn() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}
