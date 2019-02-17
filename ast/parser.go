package ast

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/zkry/go-sed/lexer"
)

type Parser struct {
	l *lexer.Lexer
	i chan lexer.Item

	curToken  lexer.Item
	peekToken lexer.Item

	lineCt int
	errors []string
	tokens []lexer.Item
}

func New(input string) *Parser {
	p := &Parser{
		errors: []string{},
	}
	p.l, p.i = lexer.New(input)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = <-p.i // TODO: Make the next token always be EOF

	// Specail next token logic here.
	if p.curTokenIs(lexer.ItemNewline) {
		p.lineCt++
	}

	p.tokens = append(p.tokens, p.peekToken)
}

// ParserProgram will parse the program that was initialized in the parer
// and return a program (list of sed commands).
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Labels = make(map[string]int)

	program.Statements = []statement{}

	for p.curToken.Type != lexer.ItemEOF {
		stmt, label := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		if label != "" {
			program.Labels[label] = len(program.Statements)
		}
		p.nextToken()
	}
	program.Tokens = make([]lexer.Item, len(p.tokens))
	copy(program.Tokens, p.tokens)
	return program
}

type ErrorList []string

func (e ErrorList) Error() string {
	buff := bytes.Buffer{}
	for _, err := range e {
		buff.WriteString(err)
	}
	return buff.String()
}

// Errors returns the list of errors encountered during the parsing process.
func (p *Parser) Errors() ErrorList {
	return p.errors
}

func (p *Parser) parseStatement() (statement, string) {
	var stmt statement

	if isStatementDelim(p.curToken.Type) {
		return nil, ""
	}

	if p.curTokenIs(lexer.ItemColon) {
		if !p.expectPeek(lexer.ItemIdent) {
			return nil, ""
		}
		lit := p.curToken.Value
		// Check if valid literal
		if lit == "" {
			p.customError("invalid label name")
			return nil, ""
		}
		return nil, lit
	}

	addr := p.parseAddress()

	switch p.curToken.Type {
	case lexer.ItemLBrace:
		// Start block
		p.nextToken()
		block := &Program{}
		block.Statements = []statement{}
		block.Labels = map[string]int{}
		for p.curToken.Type != lexer.ItemEOF && p.curToken.Type != lexer.ItemRBrace {
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
	case lexer.ItemCmd:
		switch p.curToken.Value {
		case "a":
			p.expectPeek(lexer.ItemBackslash)
			p.expectPeek(lexer.ItemLit)
			stmt = &aStmt{
				addresser:  addr,
				AppendLine: p.curToken.Value,
			}
		case "b":
			branchIdent := "$" // TODO: Find better way to signify end.
			if p.peekTokenIs(lexer.ItemIdent) {
				p.nextToken()
				branchIdent = p.curToken.Value
			}
			stmt = &bStmt{
				addresser:   addr,
				BranchIdent: branchIdent,
			}
		case "c":
			p.expectPeek(lexer.ItemBackslash)
			p.expectPeek(lexer.ItemLit)
			stmt = &cStmt{
				addresser:  addr,
				ChangeLine: p.curToken.Value,
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
			// 	cmd = p.curToken.Value
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
			p.expectPeek(lexer.ItemBackslash)
			p.expectPeek(lexer.ItemLit)
			stmt = &iStmt{
				addresser:  addr,
				InsertLine: p.curToken.Value,
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
			p.expectPeek(lexer.ItemIdent)
			stmt = &rStmt{
				addresser: addr,
				FileName:  p.curToken.Value,
			}
		case "R":
			p.expectPeek(lexer.ItemIdent)
			stmt = &r2Stmt{
				addresser: addr,
				FileName:  p.curToken.Value,
			}
		case "s":
			fa := ""
			ra := ""
			var fl sFlags
			p.expectPeek(lexer.ItemDiv)
			if p.peekTokenIs(lexer.ItemLit) {
				p.expectPeek(lexer.ItemLit)
				fa = p.curToken.Value
			}
			p.expectPeek(lexer.ItemDiv)
			if p.peekTokenIs(lexer.ItemLit) {
				p.expectPeek(lexer.ItemLit)
				ra = p.curToken.Value
			}
			p.expectPeek(lexer.ItemDiv)
			if p.peekTokenIs(lexer.ItemIdent) {
				p.expectPeek(lexer.ItemIdent)
				fl = *p.parseFlags()
			}
			stmt = &sStmt{
				addresser:   addr,
				FindAddr:    fa,
				ReplaceAddr: ra,
				Flags:       fl,
			}
		case "t":
			p.expectPeek(lexer.ItemIdent)
			stmt = &tStmt{
				addresser:   addr,
				BranchIdent: p.curToken.Value,
			}
		case "T":
			p.expectPeek(lexer.ItemIdent)
			stmt = &t2Stmt{
				addresser: addr,
				FileName:  p.curToken.Value,
			}
		case "v":
		case "w":
			p.expectPeek(lexer.ItemIdent)
			stmt = &wStmt{
				addresser: addr,
				FileName:  p.curToken.Value,
			}
		case "W":
			p.expectPeek(lexer.ItemIdent)
			stmt = &w2Stmt{
				addresser: addr,
				FileName:  p.curToken.Value,
			}
		case "x":
			stmt = &xStmt{
				addresser: addr,
			}
		case "y":
			p.expectPeek(lexer.ItemDiv)
			p.expectPeek(lexer.ItemLit)
			fa := p.curToken.Value
			p.expectPeek(lexer.ItemDiv)
			p.expectPeek(lexer.ItemLit)
			ra := p.curToken.Value
			p.expectPeek(lexer.ItemDiv)

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
	if !isStatementDelim(p.curToken.Type) {
		p.unexpectedTokenError()
	}

	return stmt, ""
}

func (p *Parser) parseAddress() addresser {
	if p.curTokenIs(lexer.ItemCmd) {
		return &blankAddress{}
	}

	addr1 := p.parseAddressPart()
	if addr1 == nil {
		return nil
	}
	switch p.curToken.Type {
	case lexer.ItemCmd:
		return addr1
	case lexer.ItemLBrace:
		return addr1
	case lexer.ItemComma:
		p.nextToken()
		addr2 := p.parseAddressPart()
		rangeAddr := &rangeAddress{Addr1: addr1, Addr2: addr2}
		if p.curToken.Type == lexer.ItemExpMark {
			p.nextToken()
			return &notAddr{Addr: rangeAddr}
		}
		return rangeAddr
	case lexer.ItemExpMark:
		p.nextToken()
		return &notAddr{Addr: addr1}
	default:
		p.unexpectedTokenError()
	}

	return nil
}

func (p *Parser) parseFlags() *sFlags {
	flg := &sFlags{}
	for {
		if p.curToken.Value[0] >= '1' && p.curToken.Value[0] <= '9' {
			flg.NFlag = int(p.curToken.Value[0] - '0')
		} else if p.curToken.Value == "g" {
			flg.GFlag = true
		} else if p.curToken.Value == "p" {
			flg.PFlag = true
		} else if p.curToken.Value == "w" {
			if p.expectPeek(lexer.ItemIdent) {
				flg.WFile = p.curToken.Value
			} else {
				p.unexpectedTokenError()
			}
			return flg // No more flags after this.
		}
		if !p.peekTokenIs(lexer.ItemIdent) {
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
	case lexer.ItemSlash:
		if !p.peekTokenIs(lexer.ItemLit) {
			// Could be a blank literal
			if p.peekTokenIs(lexer.ItemSlash) {
				p.nextToken()
				regex := regexp.MustCompile("")
				addr = &regexpAddr{Regexp: regex}
				break
			}
			return nil
		}
		p.nextToken()

		lit := p.curToken.Value
		if !p.expectPeek(lexer.ItemSlash) {
			return nil
		}
		// TODO: Should have own form of regexp
		lit = translateLiteral(lit)
		regex, err := regexp.Compile(lit)
		if err != nil {
			regex = regexp.MustCompile(".*")
		}

		addr = &regexpAddr{Regexp: regex}
	case lexer.ItemInt:
		i, err := strconv.Atoi(p.curToken.Value)
		if err != nil {
			// TODO: Have a better way of doing error handling
			panic(err)
		}
		addr = &lineNoAddr{LineNo: i}
	case lexer.ItemDollar:
		addr = &eofAddr{}
	default:
		p.unexpectedTokenError()
	}
	p.nextToken()
	return addr
}

func (p *Parser) curTokenIs(t lexer.ItemType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.ItemType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.ItemType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) lineNumber() int {
	return p.lineCt + 1
}

func (p *Parser) peekError(t lexer.ItemType) {
	msg := fmt.Sprintf("line %d: expected next token to be %s, got %s instead", p.lineNumber(), t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) unexpectedTokenError() {
	msg := fmt.Sprintf("line %d: unexpected token type %s", p.lineNumber(), p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) unexpectedFlagError(f rune) {
	msg := fmt.Sprintf("line %d: unexpected flag type %v", p.lineNumber(), f)
	p.errors = append(p.errors, msg)
}

func (p *Parser) customError(f string) {
	p.errors = append(p.errors, f)
}

func isStatementDelim(t lexer.ItemType) bool {
	return t == lexer.ItemNewline || t == lexer.ItemEOF || t == lexer.ItemSemicolon
}
