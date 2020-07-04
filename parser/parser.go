package parser

import (
	"fmt"
	"xlang/ast"
	"xlang/lexer"
	"xlang/token"
)

// Parser literally parses the code
type Parser struct {
	l         *lexer.Lexer
	errors    []string
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TypeToken]prefixParseFn
	infixParseFns  map[token.TypeToken]infixParseFn
}

// New returns a new xlang parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TypeToken]prefixParseFn)
	p.infixParseFns = make(map[token.TypeToken]infixParseFn)
	p.registerInfix(token.JUMP, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOTEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)

	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	p.nextToken()
	p.nextToken()

	return p

}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses statements and add them to the ast tree
func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if len(p.errors) > 0 && program.Line() == 0 {
			program.SetLine(p.l.Line)
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	if len(p.errors) > 0 {
		fmt.Println(p.Errors())
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.JUMP:
		{
			p.nextToken()
			return nil
		}
	case token.LET:
		{
			// We do this messy stuff because if we returned directly we wouldn't be able to check fast enough if it's nil.
			let := p.parseLetStatement()
			if let == nil {
				return nil
			}
			return let
		}

	case token.RETURN:
		{
			r := p.parseReturn()
			if r == nil {
				return nil
			}
			return r
		}
	default:
		return p.parseExpressionStatement()
	}
}

// Parses a let statement (logic)
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// Save identifier
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Param: What I expect in the peek
// Returns: If the peeked token it's not what is expected.
func (p *Parser) expectPeek(t token.TypeToken) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) curTokenIs(t token.TypeToken) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TypeToken) bool {
	return p.peekToken.Type == t
}

// Errors returns errors lmao
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TypeToken) {
	msg := fmt.Sprintf("Expected next token to be %s but it's %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
