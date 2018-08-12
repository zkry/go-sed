package lexer

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/zkry/go-sed/token"
)

// I should probably rewrite this.

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
	input    []rune
	position int

	ch     rune
	prevCh rune

	addrDiv rune
	div     rune
	s       state
	cmd     rune

	row, col int
}

func New(input string) *Lexer {
	rr := []rune{}

	for _, r := range input {
		rr = append(rr, r)
	}

	l := &Lexer{input: rr}
	l.readChar()
	l.s = stateStart
	return l
}

func (l *Lexer) readChar() {
	if l.ch == '\n' {
		l.row++
		l.col = 0
	} else {
		l.col++
	}

	l.prevCh = l.ch
	if l.position >= len(l.input) {
		l.ch = 0
		return
	}
	l.ch = l.input[l.position]
	l.position++
}

// readUntil w
func (l *Lexer) readUntil(toFunc func(rune) bool) string {
	buf := bytes.Buffer{}
	fmt.Printf("l.ch = '%c'\n", l.ch)
	for l.ch != 0 && !toFunc(l.ch) {
		fmt.Printf("l.ch > '%c'\n", l.ch)
		buf.WriteRune(l.ch)
		l.readChar()
	}
	fmt.Printf("l.ch = '%c'\n", l.ch)
	return buf.String()
}

func (l *Lexer) readUntilEscape(toFunc func(rune) bool) string {
	buf := bytes.Buffer{}
	for l.ch != 0 && !toFunc(l.ch) {
		// Allow \ to escape the toFunc.
		// Example, if toFunc is unicode.isDigit, allow \1 to bypass this.
		if l.ch == '\\' {
			l.readChar()
			if !toFunc(l.ch) {
				buf.WriteRune('\\')
			}
			if l.ch == 0 {
				break
			}
		}
		if l.ch != 0 {
			buf.WriteRune(l.ch)
		}
		l.readChar()
	}
	return buf.String()
}

// func (l *Lexer) readWhileWithPrev(toFunc func(prev, at rune) bool) string {
// 	start := l.position
// 	for toFunc(l.prevCh, l.ch) && l.ch != 0 {
// 		l.readChar()
// 	}
// 	return string(l.input[start:l.position])
// }

func (l *Lexer) readWhileEscape(toFunc func(rune) bool) string {
	return l.readUntilEscape(not(toFunc))
}

func (l *Lexer) readWhile(toFunc func(rune) bool) string {
	return l.readUntil(not(toFunc))
}

func isCmd(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '='
}

func isAlnum(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r)
}

// isASpace is different than unicode.IsSpace in that isASpace
func isASpace(r rune) bool {
	return r == ' '
}

func isNewline(r rune) bool {
	return r == '\n'
}

func isNewCommand(r rune) bool {
	return isNewlineOrEOF(r) || r == ';'
}

func isNewlineOrEOF(r rune) bool {
	return r == '\n' || r == 0
}

func not(f func(r rune) bool) func(rune) bool {
	return func(r rune) bool {
		return !f(r)
	}
}

func newToken(t token.Type, r rune) token.Token {
	return token.Token{Type: t, Literal: string(r)}
}

const debug = false

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

	if debug {
		fmt.Printf("%c: %v\n", l.ch, l.s)
	}

	switch l.s {
	case stateStart:
		tok = l.lexStart()
	case stateAddr:
		tok = l.lexAddr()
	case stateEndAddr:
		tok = l.lexEndAddr()
	case stateLabel:
		tok = l.lexLabel()
	case state2ndAddrStart:
		tok = l.lex2ndAddrStart()
	case state2ndAddr:
		tok = l.lex2ndAddr()
	case stateEnd2ndAddr:
		tok = l.lexEnd2ndAddr()
	case stateCmd:
		tok = l.lexCmd()
	case stateFindPtn:
		tok = l.lexFind()
	case stateReplacePtn:
		tok = l.lexReplace()
	case stateFlags:
		tok = l.lexFlag()
	case statePostFlag:
		tok = l.lexPostFlag()
	case stateReadline:
		tok = l.lexReadLine()
	default:
		panic("no handler for state '" + l.s + "'")
	}

	tok.Start = startPos
	tok.End = l.position
	return tok
}

func (l *Lexer) lexStart() token.Token {
	// Reset state
	l.div = '/'
	l.addrDiv = '/'

	l.readWhile(isASpace)

	tok := token.Token{}
	switch l.ch {
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
		l.readChar()
	case '#':
		tok = newToken(token.NEWLINE, l.ch)
		l.readUntil(isNewlineOrEOF)
	case '\n':
		tok = newToken(token.NEWLINE, l.ch)
		l.readChar()
	case '/':
		tok = newToken(token.SLASH, l.ch)
		l.s = stateAddr
		l.readChar()
	case ':':
		tok = newToken(token.COLON, l.ch)
		l.s = stateLabel
		l.readChar()
	case '$':
		tok = newToken(token.DOLLAR, l.ch)
		l.s = stateEndAddr
		l.readChar()
	case '}':
		tok = newToken(token.RBRACE, l.ch)
		l.readUntil(isNewCommand)
		l.readChar()
	case '\\':
		l.readChar()
		l.addrDiv = l.ch
		l.s = stateAddr
		tok = newToken(token.SLASH, l.ch)
		l.readChar()
	default:
		if unicode.IsDigit(l.ch) {
			tok.Literal = l.readWhile(unicode.IsDigit)
			tok.Type = token.INT
			l.s = stateEndAddr
		} else if isCmd(l.ch) {
			tok = newToken(token.CMD, l.ch)
			l.cmd = l.ch
			l.s = stateCmd
			l.readChar()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
			l.readChar()
		}
	}
	fmt.Printf("->'%c'\n", l.ch)
	return tok
}

func (l *Lexer) lexLabel() token.Token {
	return token.Token{}
}

func (l *Lexer) lexReadLine() token.Token {
	return token.Token{}
}

func (l *Lexer) lexPostFlag() token.Token {
	return token.Token{}
}

// lexFlag extracts the flag portion of the s/1/2/f  pattern
func (l *Lexer) lexFlag() token.Token {
	return token.Token{}
}

// lexReplace extracts the second part of the s/1/2/f pattern
func (l *Lexer) lexReplace() token.Token {
	return token.Token{}
}

func (l *Lexer) lexFind() token.Token {
	return token.Token{}
}

func (l *Lexer) lexCmd() token.Token {
	return token.Token{}
}

func (l *Lexer) lexEnd2ndAddr() token.Token {
	return token.Token{}
}

func (l *Lexer) lex2ndAddrStart() token.Token {
	return token.Token{}
}

func (l *Lexer) lex2ndAddr() token.Token {
	return token.Token{}
}

func (l *Lexer) lexEndAddr() token.Token {
	return token.Token{}
}

func (l *Lexer) lexAddr() token.Token {
	return token.Token{}
}
