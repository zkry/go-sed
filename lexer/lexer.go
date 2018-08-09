package lexer

import (
	"bytes"
	"errors"

	"github.com/zkry/go-sed/token"
)

type state string

const (
	stateStart state = "START"

	stateLabel state = "LABEL"

	stateAddr         state = "LITERAL"
	stateEndAddr      state = "END_ADDR"
	state2ndAddrStart state = "2ND_ADDR_START"
	state2ndAddr      state = "2ND_ADDR"
	stateEnd2ndAddr   state = "END_2ND_ADDR"

	stateCmd state = "CMD" // The state after reading a command

	stateFindPtn    state = "FIND_PATTERN"
	stateReplacePtn state = "REPLACE_PATTERN"
	stateFlags      state = "FLAGS"
	statePostFlag   state = "POST_FLAG"

	stateReadline state = "READ_LINE"
)

// Lexer represents the state of the lexer object
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	prevCh       byte
	addrDiv      byte
	div          byte
	s            state
	cmd          byte // The command that we have seen
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	l.s = stateStart
	return l
}

func (l *Lexer) rewindChar() {
	if l.readPosition == 0 {
		return
	}
	l.position--
	l.readPosition--
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
		if l.readPosition > 0 {
			l.prevCh = l.input[l.readPosition-1]
		}
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
		if l.readPosition > 0 {
			l.prevCh = l.input[l.readPosition-1]
		}
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken takes the character that the Lexer is at
// and returns a corresponding token, then increments Lexer
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	startPos := l.position

	// Check for the end of file
	if l.ch == 0 {
		tok.Literal = ""
		tok.Type = token.EOF
		return tok
	}

	// fmt.Printf("%c: ", l.ch)

	switch l.s {
	case stateStart:
		// fmt.Println("START")
		tok = l.lexStart()
	case stateAddr:
		// fmt.Println("ADDR")
		tok = l.lexAddr()
	case stateEndAddr:
		// fmt.Println("ADDR_END")
		tok = l.lexEndAddr()
	case stateLabel:
		// fmt.Println("LABEL")
		tok = l.lexLabel()
	case state2ndAddrStart:
		// fmt.Println("ADDR2_START")
		tok = l.lex2ndAddrStart()
	case state2ndAddr:
		// fmt.Println("ADDR2")
		tok = l.lex2ndAddr()
	case stateEnd2ndAddr:
		// fmt.Println("ADDR2_END")
		tok = l.lexEnd2ndAddr()
	case stateCmd:
		// fmt.Println("CMD")
		tok = l.lexCmd()
	case stateFindPtn:
		// fmt.Println("FIND")
		tok = l.lexFind()
	case stateReplacePtn:
		// fmt.Println("REPLACE")
		tok = l.lexReplace()
	case stateFlags:
		// fmt.Println("FLAG")
		tok = l.lexFlag()
	case statePostFlag:
		// fmt.Println("POST_FLAG")
		tok = l.lexPostFlag()
	case stateReadline:
		// fmt.Println("READ_LINE")
		tok = l.lexReadLine()
	default:
		// fmt.Println("NOT COVERED")
	}

	l.readChar()

	tok.Start = startPos
	tok.End = l.position
	return tok
}

func (l *Lexer) lexLabel() token.Token {
	var tok token.Token

	if l.ch == '\n' {
		tok = newToken(token.NEWLINE, l.ch)
		l.s = stateStart
		return tok
	}

	l.readUntil(isSpace)
	tok.Literal = l.readUntil(not(isNewlineOrEOF))
	tok.Type = token.IDENT
	return tok
}

func (l *Lexer) lexReadLine() token.Token {
	var tok token.Token

	if isSpace(l.ch) {
		l.readUntil(isSpace)
	}

	switch l.ch {
	case '\n':
		l.readChar()

		lit, err := l.readLineLiteral()
		tok.Literal = lit
		if err != nil {
			tok.Type = token.ILLEGAL
			return tok
		}
		tok.Type = token.LIT
		l.s = stateStart
	default:
		tok = newToken(token.ILLEGAL, l.ch)
	}
	return tok
}

func (l *Lexer) lexPostFlag() token.Token {
	var tok token.Token

	switch l.ch {
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '\n':
		tok = newToken(token.NEWLINE, l.ch)
		l.s = stateStart
	case ';':
		// TODO: ; separated commands
	default:
		if l.ch == ' ' {
			l.readUntil(isSpace)
			l.readChar()
		}
		tok.Literal = l.readUntil(not(isNewlineOrEOF))
		tok.Type = token.IDENT
		l.s = stateStart
	}

	return tok
}

// lexFlag extracts the flag portion of the s/1/2/f  pattern
func (l *Lexer) lexFlag() token.Token {
	var tok token.Token

	l.readUntil(isSpace)
	switch l.ch {
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
		l.s = stateStart
	case '\n':
		tok = newToken(token.NEWLINE, l.ch)
		l.s = stateStart
	case ' ':
		l.s = statePostFlag
	default:
		if isLetter(l.ch) || isNumber(l.ch) {
			tok = newToken(token.IDENT, l.ch)
			if l.ch == 'r' || l.ch == 'w' {
				l.s = statePostFlag
			}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	return tok
}

// lexReplace extracts the second part of the s/1/2/f pattern
func (l *Lexer) lexReplace() token.Token {
	var tok token.Token
	switch l.ch {
	case l.div:
		tok = newToken(token.DIV, l.ch)
		l.s = stateFlags
	default:
		tok.Literal = l.readUntil(func(b byte) bool { return b != l.div })
		tok.Type = token.LIT
	}
	return tok
}

func (l *Lexer) lexFind() token.Token {
	var tok token.Token

	switch l.ch {
	case l.div:
		tok = newToken(token.DIV, l.ch)
		l.s = stateReplacePtn
	default:
		tok.Literal = l.readUntil(func(b byte) bool { return b != l.div })
		tok.Type = token.LIT
	}
	return tok
}

func (l *Lexer) lexCmd() token.Token {
	var tok token.Token

	switch l.cmd {
	case 's', 'y':
		tok = newToken(token.DIV, l.ch)
		l.div = l.ch
		l.s = stateFindPtn
	case 'c', 'a', 'i':
		if isSpace(l.ch) {
			l.readUntil(isSpace)
			l.readChar()
		}
		switch l.ch {
		case '\\':
			tok = newToken(token.BACKSLASH, l.ch)
			l.s = stateReadline
		default:
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case 'r', 'w':
		l.readUntil(isSpace)
		if isLetter(l.ch) {
			tok.Literal = l.readUntil(isNewline)
			tok.Type = token.IDENT
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	case 'b', 't':
		if l.ch == '\n' {
			tok = newToken(token.NEWLINE, '\n')
			l.s = stateStart
			return tok
		}
		if l.ch == ' ' {
			l.readChar()
		}
		if isLetter(l.ch) {
			tok.Literal = l.readUntil(not(isNewlineOrEOF))
			tok.Type = token.IDENT
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	default:
		switch l.ch {
		case ';':
			tok = newToken(token.SEMICOLON, l.ch)
			l.s = stateStart
		case '\n':
			tok = newToken(token.NEWLINE, l.ch)
			l.s = stateStart
		default:
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	return tok
}

func (l *Lexer) lexEnd2ndAddr() token.Token {
	var tok token.Token
	// TODO: abstract readUntil+readChar to readF
	// and have other version readT. This is simmilar
	// to the vim t and f commands. For one
	if l.ch == ' ' {
		l.readUntil(isSpace)
		l.readChar()
	}

	switch l.ch {
	case '!':
		tok = newToken(token.EXPLMARK, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
		l.readUntil(not(isNewlineOrEOF))
		l.s = stateStart
	default:
		if isCmd(l.ch) {
			tok = newToken(token.CMD, l.ch)
			l.cmd = l.ch
			l.s = stateCmd
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	return tok
}

func (l *Lexer) lex2ndAddrStart() token.Token {
	tok := l.lexStart()
	if l.s == stateAddr {
		l.s = state2ndAddr
	} else if l.s == stateEndAddr {
		l.s = stateEnd2ndAddr
	}
	return tok
}

func (l *Lexer) lex2ndAddr() token.Token {
	tok := l.lexAddr()
	if l.s == stateEndAddr {
		l.s = stateEnd2ndAddr
	}
	return tok
}

func (l *Lexer) lexEndAddr() token.Token {
	var tok token.Token

	if isSpace(l.ch) {
		l.readUntil(isSpace)
		defer l.readChar()
	}
	switch l.ch {
	case '!':
		tok = newToken(token.EXPLMARK, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
		l.s = state2ndAddrStart
	case '{':
		tok = newToken(token.LBRACE, l.ch)
		l.s = stateStart
	default:
		if isCmd(l.ch) {
			tok = newToken(token.CMD, l.ch)
			l.cmd = l.ch
			l.s = stateCmd
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	return tok
}

func (l *Lexer) lexAddr() token.Token {
	var tok token.Token

	switch l.ch {
	case l.addrDiv:
		tok = newToken(token.SLASH, l.ch)
		l.s = stateEndAddr
	default:
		// TODO: /this is \\/ my address/ s/this/that/  will fail!
		tok.Literal = l.readUntilWithPrev(func(p, a byte) bool { return a != l.addrDiv || p == '\\' })
		tok.Type = token.LIT
	}
	return tok
}

func (l *Lexer) lexStart() token.Token {
	var tok token.Token

	// Reset state
	l.div = '/'
	l.addrDiv = '/'

	if isSpace(l.ch) {
		l.readUntil(isSpace)
		l.readChar()
	}
	switch l.ch {
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '#':
		tok = newToken(token.NEWLINE, l.ch)
		l.readUntil(not(isNewlineOrEOF))
		l.readChar()
	case '\n':
		tok = newToken(token.NEWLINE, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
		l.s = stateAddr
	case ':':
		tok = newToken(token.COLON, l.ch)
		l.s = stateLabel
	case '$':
		tok = newToken(token.DOLLAR, l.ch)
		l.s = stateEndAddr
	case '}':
		tok = newToken(token.RBRACE, l.ch)
		l.readUntil(not(isNewlineOrEOF))
	case '\\':
		l.readChar()
		l.addrDiv = l.ch
		l.s = stateAddr
		tok = newToken(token.SLASH, l.ch)
	default:
		if isNumber(l.ch) {
			tok.Literal = l.readUntil(isNumber)
			tok.Type = token.INT
			l.s = stateEndAddr
		} else if isCmd(l.ch) {
			tok = newToken(token.CMD, l.ch)
			l.cmd = l.ch
			l.s = stateCmd
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	return tok

}

func (l *Lexer) readUntilWithPrev(toFunc func(prev, at byte) bool) string {
	position := l.position
	for toFunc(l.prevCh, l.ch) && l.ch != 0 {
		l.readChar()
	}
	ret := l.input[position:l.position]
	if position < l.position {
		l.rewindChar()
	}
	return ret
}

func (l *Lexer) readUntil(toFunc func(byte) bool) string {
	position := l.position
	for toFunc(l.ch) && l.ch != 0 {
		l.readChar()
	}
	ret := l.input[position:l.position]
	if position < l.position {
		l.rewindChar()
	}
	return ret
}

func (l *Lexer) readLineLiteral() (string, error) {
	buf := new(bytes.Buffer)
	for l.ch != '\n' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			if l.ch != '\n' && l.ch != '\\' {
				return buf.String(), errors.New("invalid escape sequence")
			}
		}
		buf.WriteByte(l.ch)
		l.readChar()
	}
	l.rewindChar() // leave function off at a \n so lexStart will give us newline token
	return buf.String(), nil
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isCmd(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '='
}

func isAlnum(ch byte) bool {
	return isLetter(ch) || isNumber(ch)
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t'
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

func isNewline(ch byte) bool {
	return ch == '\n'
}

func isNewlineOrEOF(ch byte) bool {
	return ch == '\n' || ch == 0
}

func not(f func(byte) bool) func(byte) bool {
	return func(b byte) bool {
		return !f(b)
	}
}

func newToken(t token.Type, ch byte) token.Token {
	return token.Token{Type: t, Literal: string(ch)}
}
