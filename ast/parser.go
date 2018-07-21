package ast

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

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
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParserProgram will parse the program that was initialized in the parer
// and return a program (list of sed commands).
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Labels = make(map[string]int)

	program.Statements = []statement{}

	for p.curToken.Type != token.EOF {
		stmt, label := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		if label != "" {
			program.Labels[label] = len(program.Statements)
		}
		p.nextToken()
	}
	return program
}

// Errors returns the list of errors encountered during the parsing process.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() (statement, string) {

	var stmt statement

	if p.curToken.IsStatementDelim() {
		return nil, ""
	}

	if p.curTokenIs(token.COLON) {
		if !p.expectPeek(token.IDENT) {
			return nil, ""
		}
		lit := p.curToken.Literal
		// Check if valid literal
		if lit == "" {
			p.customError("invalid label name")
			return nil, ""
		}
		return nil, lit
	}

	addr := p.parseAddress()

	switch p.curToken.Type {
	case token.LBRACE:
		// Start block
		p.nextToken()
		block := &Program{}
		block.Statements = []statement{}
		block.Labels = map[string]int{}
		for p.curToken.Type != token.EOF && p.curToken.Type != token.RBRACE {
			stmt, l := p.parseStatement()
			if stmt != nil {
				block.Statements = append(block.Statements, stmt)
			}
			if l != "" {
				block.Labels[l] = len(block.Statements)
			}
			p.nextToken()
		}
		stmt = &blockStmt{
			Code:      block,
			addresser: addr,
		}
	case token.CMD:
		switch p.curToken.Literal {
		case "a":
			p.expectPeek(token.BACKSLASH)
			p.expectPeek(token.LIT)
			stmt = &aStmt{
				addresser:  addr,
				AppendLine: p.curToken.Literal,
			}
		case "b":
			p.expectPeek(token.IDENT)
			stmt = &bStmt{
				addresser:   addr,
				BranchIdent: p.curToken.Literal,
			}
		case "c":
			p.expectPeek(token.BACKSLASH)
			p.expectPeek(token.LIT)
			stmt = &cStmt{
				addresser:  addr,
				ChangeLine: p.curToken.Literal,
			}
		case "d":
			stmt = &dStmt{
				addresser: addr,
			}
		case "D":
			stmt = &d2Stmt{
				addresser: addr,
			}
		case "e":
			// cmd := ""
			// if p.peekTokenIs(token.LIT) {
			// 	p.expectPeek(token.LIT)
			// 	cmd = p.curToken.Literal
			// }
			// stmt = &EStmt{
			// 	Addresser: addr,
			// 	Command:   cmd,
			// }
		case "F": // Is this even a command
		case "g":
			stmt = &gStmt{
				addresser: addr,
			}
		case "G":
			stmt = &g2Stmt{
				addresser: addr,
			}
		case "h":
			stmt = &hStmt{
				addresser: addr,
			}
		case "H":
			stmt = &h2Stmt{
				addresser: addr,
			}
		case "i":
			p.expectPeek(token.BACKSLASH)
			p.expectPeek(token.LIT)
			stmt = &iStmt{
				addresser:  addr,
				InsertLine: p.curToken.Literal,
			}
		case "l":
			stmt = &lStmt{
				addresser: addr,
			}
		case "n":
			stmt = &nStmt{
				addresser: addr,
			}
		case "N":
			stmt = &n2Stmt{
				addresser: addr,
			}
		case "p":
			stmt = &pStmt{
				addresser: addr,
			}
		case "P":
			stmt = &p2Stmt{
				addresser: addr,
			}
		case "q":
			stmt = &qStmt{
				addresser: addr,
			}
		case "r":
			p.expectPeek(token.IDENT)
			stmt = &rStmt{
				addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "R":
			p.expectPeek(token.IDENT)
			stmt = &r2Stmt{
				addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "s":
			fa := ""
			ra := ""
			var fl sFlags
			p.expectPeek(token.DIV)
			if p.peekTokenIs(token.LIT) {
				p.expectPeek(token.LIT)
				fa = p.curToken.Literal
			}
			p.expectPeek(token.DIV)
			if p.peekTokenIs(token.LIT) {
				p.expectPeek(token.LIT)
				ra = p.curToken.Literal
			}
			p.expectPeek(token.DIV)
			if p.peekTokenIs(token.IDENT) {
				p.expectPeek(token.IDENT)
				fl = *p.parseFlags()
			}
			stmt = &sStmt{
				addresser:   addr,
				FindAddr:    fa,
				ReplaceAddr: ra,
				Flags:       fl,
			}
		case "t":
			p.expectPeek(token.IDENT)
			stmt = &tStmt{
				addresser:   addr,
				BranchIdent: p.curToken.Literal,
			}
		case "T":
			p.expectPeek(token.IDENT)
			stmt = &t2Stmt{
				addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "v":
		case "w":
			p.expectPeek(token.IDENT)
			stmt = &wStmt{
				addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "W":
			p.expectPeek(token.IDENT)
			stmt = &w2Stmt{
				addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "x":
			stmt = &xStmt{
				addresser: addr,
			}
		case "y":
			p.expectPeek(token.DIV)
			p.expectPeek(token.LIT)
			fa := p.curToken.Literal
			p.expectPeek(token.DIV)
			p.expectPeek(token.LIT)
			ra := p.curToken.Literal
			p.expectPeek(token.DIV)

			var err error
			stmt, err = newYStmt(fa, ra, addr)
			if err != nil {
				return nil, ""
			}
		case "z":
			stmt = &zStmt{
				addresser: addr,
			}
		case "=":
			stmt = &equStmt{
				addresser: addr,
			}
		}
	default:
		p.unexpectedTokenError()
	}

	p.nextToken()
	if !p.curToken.IsStatementDelim() {
		p.unexpectedTokenError()
	}

	return stmt, ""
}

func (p *Parser) parseAddress() addresser {
	if p.curTokenIs(token.CMD) {
		return &blankAddress{}
	}

	addr1 := p.parseAddressPart()
	if addr1 == nil {
		return nil
	}
	switch p.curToken.Type {
	case token.CMD:
		return addr1
	case token.LBRACE:
		return addr1
	case token.COMMA:
		p.nextToken()
		addr2 := p.parseAddressPart()
		return &rangeAddress{Addr1: addr1, Addr2: addr2}
	default:
		p.unexpectedTokenError()
	}
	return nil
}

func (p *Parser) parseFlags() *sFlags {
	flg := &sFlags{}
	for {
		if p.curToken.Literal[0] >= '1' && p.curToken.Literal[0] <= '9' {
			flg.NFlag = int(p.curToken.Literal[0] - '0')
		} else if p.curToken.Literal == "g" {
			flg.GFlag = true
		} else if p.curToken.Literal == "p" {
			flg.PFlag = true
		} else if p.curToken.Literal == "w" {
			if p.expectPeek(token.IDENT) {
				flg.WFile = p.curToken.Literal
			} else {
				p.unexpectedTokenError()
			}
			return flg // No more flags after this.
		}
		if !p.peekTokenIs(token.IDENT) {
			return flg
		}
		p.nextToken()
	}
}

// translateLiteral performs the translation from user input
// to the string that shoudl be processed in the regexp. This
// incluedes processing escape characters.
func translateLiteral(l string) string {
	var retData bytes.Buffer
	var escState bool
	for _, r := range l {
		if escState {
			switch {
			case r == '\\':
				retData.WriteRune('\\')
			case r == 'n':
				retData.WriteRune('\n')
				// TODO: Check case of the literal that was used to divide
			default:
				// TODO: Error
			}
			escState = false
		} else {
			switch {
			case r == '\\':
				escState = true
			default:
				retData.WriteRune(r)
			}
		}
	}
	return retData.String()
}
func (p *Parser) parseAddressPart() addresser {
	var addr addresser
	switch p.curToken.Type {
	case token.SLASH:
		if !p.expectPeek(token.LIT) {
			return nil
		}
		lit := p.curToken.Literal
		if !p.expectPeek(token.SLASH) {
			return nil
		}
		// TODO: Should have own form of regexp
		lit = translateLiteral(lit)
		regexp := regexp.MustCompile(lit)
		addr = &regexpAddr{Regexp: regexp}
	case token.INT:
		i, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			// TODO: Have a better way of doing error handling
			panic(err)
		}
		addr = &lineNoAddr{LineNo: i}
	case token.DOLLAR:
		addr = &eofAddr{}
	default:
		p.unexpectedTokenError()
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
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) unexpectedTokenError() {
	msg := fmt.Sprintf("unexpected token type %s", p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) unexpectedFlagError(f rune) {
	msg := fmt.Sprintf("unexpected flag type %v", f)
	p.errors = append(p.errors, msg)
}

func (p *Parser) customError(f string) {
	p.errors = append(p.errors, f)
}
