package parser

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/zkry/go-sed/ast"
	"github.com/zkry/go-sed/lexer"
	"github.com/zkry/go-sed/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		// stmt := p.parseStatement()
		// if stmt != nil {
		// 	program.Statements = append(program.Statements, stmt)
		// }
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() *ast.Statement {
	stmt := &ast.Statement{}

	if p.curTokenIs(token.LBRACE) {
		// END BLOCK
	}

	if p.curTokenIs(token.COLON) {
		// ADD LABEL
	}

	addr := p.parseAddress()
	stmt.Addr = addr

	switch p.curToken.Type {
	case token.RBRACE:
		// Start block
	case token.CMD:
		stmt.Cmd = p.curToken.Literal

		if stmt.Cmd == "s" || stmt.Cmd == "y" {
			p.expectPeek(token.DIV)
			p.expectPeek(token.LIT)
			stmt.Arg1 = p.curToken.Literal
			p.expectPeek(token.DIV)
			p.expectPeek(token.LIT)
			stmt.Arg2 = p.curToken.Literal
			p.expectPeek(token.DIV)

			p.nextToken()
			for p.curTokenIs(token.IDENT) {
				stmt.Flags = append(stmt.Flags, p.curToken.Literal)
				p.nextToken()
			}
		}
	default:
		fmt.Println("Could not parse token:", p.curToken)
	}

	return stmt
}

func (p *Parser) parseAddress() ast.Address {
	if p.curTokenIs(token.CMD) {
		return &ast.BlankAddress{}
	}

	addr1 := p.parseAddressPart()
	if addr1 == nil {
		return nil
	}
	switch p.curToken.Type {
	case token.CMD:
		return addr1
	case token.COMMA:
		p.nextToken()
		addr2 := p.parseAddressPart()
		return &ast.RangeAddress{Addr1: addr1, Addr2: addr2}
	default:
		fmt.Println("parseAddress: Could not parse:", p.curToken)
	}
	return nil
}

func (p *Parser) parseAddressPart() ast.Address {
	var addr ast.Address
	switch p.curToken.Type {
	case token.SLASH:
		if !p.expectPeek(token.LIT) {
			return nil
		}
		lit := p.curToken.Literal
		if !p.expectPeek(token.SLASH) {
			return nil
		}
		p.nextToken()
		// TODO: Should have own form of regexp
		regexp := regexp.MustCompile(lit)
		addr = &ast.RegexpAddr{Regexp: regexp}
	case token.INT:
		i, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			// TODO: Have a better way of doing error handling
			panic(err)
		}
		addr = &ast.LineNoAddr{LineNo: i}
	case token.DOLLAR:
		return &ast.EOFAddr{}
	default:
		fmt.Println("parseAddressPart: Could not parse:", p.curToken)
	}
	p.nextToken()
	return addr
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}
