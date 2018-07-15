package parser

import (
	"bytes"
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
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Labels = make(map[string]int)

	program.Statements = []ast.Statement{}

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

func (p *Parser) parseStatement() (ast.Statement, string) {

	var stmt ast.Statement

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
		block := &ast.Program{}
		block.Statements = []ast.Statement{}
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
		stmt = &ast.BlockStmt{
			Code:      block,
			Addresser: addr,
		}
	case token.CMD:
		switch p.curToken.Literal {
		case "a":
			p.expectPeek(token.BACKSLASH)
			p.expectPeek(token.LIT)
			stmt = &ast.AStmt{
				Addresser:  addr,
				AppendLine: p.curToken.Literal,
			}
		case "b":
			p.expectPeek(token.IDENT)
			stmt = &ast.BStmt{
				Addresser:   addr,
				BranchIdent: p.curToken.Literal,
			}
		case "c":
			p.expectPeek(token.BACKSLASH)
			p.expectPeek(token.LIT)
			stmt = &ast.CStmt{
				Addresser:  addr,
				ChangeLine: p.curToken.Literal,
			}
		case "d":
			stmt = &ast.DStmt{
				Addresser: addr,
			}
		case "D":
			stmt = &ast.D2Stmt{
				Addresser: addr,
			}
		case "e":
			// cmd := ""
			// if p.peekTokenIs(token.LIT) {
			// 	p.expectPeek(token.LIT)
			// 	cmd = p.curToken.Literal
			// }
			// stmt = &ast.EStmt{
			// 	Addresser: addr,
			// 	Command:   cmd,
			// }
		case "F": // Is this even a command
		case "g":
			stmt = &ast.GStmt{
				Addresser: addr,
			}
		case "G":
			stmt = &ast.G2Stmt{
				Addresser: addr,
			}
		case "h":
			stmt = &ast.HStmt{
				Addresser: addr,
			}
		case "H":
			stmt = &ast.H2Stmt{
				Addresser: addr,
			}
		case "i":
			p.expectPeek(token.BACKSLASH)
			p.expectPeek(token.LIT)
			stmt = &ast.IStmt{
				Addresser:  addr,
				InsertLine: p.curToken.Literal,
			}
		case "l":
			stmt = &ast.LStmt{
				Addresser: addr,
			}
		case "n":
			stmt = &ast.NStmt{
				Addresser: addr,
			}
		case "N":
			stmt = &ast.N2Stmt{
				Addresser: addr,
			}
		case "p":
			stmt = &ast.PStmt{
				Addresser: addr,
			}
		case "P":
			stmt = &ast.P2Stmt{
				Addresser: addr,
			}
		case "q":
			stmt = &ast.QStmt{
				Addresser: addr,
			}
		case "r":
			p.expectPeek(token.IDENT)
			stmt = &ast.RStmt{
				Addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "R":
			p.expectPeek(token.IDENT)
			stmt = &ast.R2Stmt{
				Addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "s":
			fa := ""
			ra := ""
			var fl ast.SFlags
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
			stmt = &ast.SStmt{
				Addresser:   addr,
				FindAddr:    fa,
				ReplaceAddr: ra,
				Flags:       fl,
			}
		case "t":
			p.expectPeek(token.IDENT)
			stmt = &ast.TStmt{
				Addresser:   addr,
				BranchIdent: p.curToken.Literal,
			}
		case "T":
			p.expectPeek(token.IDENT)
			stmt = &ast.T2Stmt{
				Addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "v":
		case "w":
			p.expectPeek(token.IDENT)
			stmt = &ast.WStmt{
				Addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "W":
			p.expectPeek(token.IDENT)
			stmt = &ast.W2Stmt{
				Addresser: addr,
				FileName:  p.curToken.Literal,
			}
		case "x":
			stmt = &ast.XStmt{
				Addresser: addr,
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
			stmt, err = ast.NewYStmt(fa, ra, addr)
			if err != nil {
				return nil, ""
			}
		case "z":
			stmt = &ast.ZStmt{
				Addresser: addr,
			}
		case "=":
			stmt = &ast.EquStmt{
				Addresser: addr,
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

func (p *Parser) parseAddress() ast.Addresser {
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
	case token.LBRACE:
		return addr1
	case token.COMMA:
		p.nextToken()
		addr2 := p.parseAddressPart()
		return &ast.RangeAddress{Addr1: addr1, Addr2: addr2}
	default:
		p.unexpectedTokenError()
	}
	return nil
}

func (p *Parser) parseFlags() *ast.SFlags {
	flg := &ast.SFlags{}
	for {
		fmt.Println("Current Token: ", p.curToken.Type, p.curToken.Literal)
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
func (p *Parser) parseAddressPart() ast.Addresser {
	var addr ast.Addresser
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
		addr = &ast.RegexpAddr{Regexp: regexp}
	case token.INT:
		i, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			// TODO: Have a better way of doing error handling
			panic(err)
		}
		addr = &ast.LineNoAddr{LineNo: i}
	case token.DOLLAR:
		addr = &ast.EOFAddr{}
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
