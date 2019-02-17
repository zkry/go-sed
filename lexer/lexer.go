package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type state string

var validCommands = [...]rune{
	's', // substitute

	'q', // quit
	'd', // delete
	'p', // print
	'n', // next
	'b', // branch
	't', // test
	'y', // transliterate chars
	'a', // append after line
	'i', // insert before line
	'c', // change
	'=', // current line number
	'l', // list (print unambiguously)
	'r', // read
	'w', // write
	'D', // delete to newline
	'N', // next to newline
	'P', // print to newline
	'h', // hold pattern
	'H', // hold to newline
	'g', // grab hold
	'G', // grab with newline
	'x', // exchange hold and pattern
}

var validGNUCommands = [...]rune{
	'e', // execute
	'F', // filename
	'Q', // quit with no-print
	'R', // queue read file
	'T', // branch if no subs
	'v', // version
	'W', // write till newline
	'z', // empty pattern space
}

type stateFn func(*Lexer) stateFn

type escape struct {
	pos   int
	width int
}

// Lexer represents the state of the lexer object
type Lexer struct {
	name  string    // used only for error reports
	input string    // the input that the lexer will run on
	start int       // the start position of this item
	pos   int       // the current position we are at
	width int       // the width of the last rune read
	items chan Item // the cannel to which we send our output tokens
}

func New(input string) (*Lexer, chan Item) {
	l := &Lexer{
		name:  "Test Lexer",
		input: input,
		start: 0,
		pos:   0,
		width: -1, // we haven't read anything but startState shoudn't go back
		items: make(chan Item),
	}
	go l.run()
	return l, l.items
}

func (l *Lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *Lexer) emit(t ItemType) {
	l.items <- Item{t, l.input[l.start:l.pos], l.pos}
	l.start = l.pos
}

func (l *Lexer) escapePrev() {
	l.input = l.input[:l.pos-l.width] + l.input[l.pos:]
	l.pos -= l.width
}

func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

// backup moves the lexer back one rune. Can only be used once after
// calling next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// peek returns the next rune like next() but does not move
// the position of the lexer.
func (l *Lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

// accept returns true if the next rune is in the set of valid runes.
func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun accepts as much instances of a character in valid as possible.
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{
		Type:  ItemError,
		Value: fmt.Sprintf(format, args...),
	}
	return nil
}

func isCommand(r rune) bool {
	for _, cmd := range validCommands {
		if r == cmd {
			return true
		}
	}
	return false
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isFlag(r rune) bool {
	return isNumeric(r) || strings.ContainsRune("mMiIewpg1234567890r", r)

}

func isNumeric(r rune) bool {
	return unicode.IsDigit(r)
}

func lexStart(l *Lexer) stateFn {
	r := l.next()
	switch {
	case r == 0:
		l.emit(ItemEOF)
		return nil
	case isSpace(r):
		l.ignore()
		return lexStart
	case isNumeric(r):
		l.acceptRun("0123456789")
		l.emit(ItemInt)
		return lexNextAddrOrCommand
	case isCommand(r):
		l.backup()
		return lex2ndAddrDone
	case r == '$':
		l.emit(ItemDollar)
		return lexNextAddrOrCommand
	case r == ':':
		l.emit(ItemColon)
		return lexIdentToEnd
	case r == '/':
		l.emit(ItemSlash)
		return lexInsideAddr('/', lexNextAddrOrCommand)
	case r == '\n':
		l.emit(ItemNewline)
		return lexStart
	case r == '\\':
		l.ignore()
		delim := l.next()
		l.emit(ItemSlash)
		return lexInsideAddr(delim, lexNextAddrOrCommand)
	case r == ';':
		l.emit(ItemSemicolon)
		return lexStart
	case r == '#':
		l.emit(ItemNewline)
		for {
			switch l.next() {
			case '\n':
				l.ignore()
				return lexStart
			case 0:
				l.ignore()
				l.emit(ItemEOF)
				return nil
			}
		}
	case r == '{':
		l.emit(ItemLBrace)
		return lexStart
	case r == '}':
		l.emit(ItemRBrace)
		return lexStart
	}

	return l.errorf("no symbol found in start state")
}

// lexNextAddrOrCommand lexes the portion after the first address.
// There could either be a command or another address portion coming.
func lexNextAddrOrCommand(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == 0:
			l.emit(ItemEOF)
			return nil
		case isSpace(r):
			l.ignore()
		case isCommand(r):
			l.backup()
			return lex2ndAddrDone
		case r == '{':
			l.emit(ItemLBrace)
			return lexStart
		case r == '}':
			l.emit(ItemRBrace)
			return lexStart
		case r == ',':
			l.emit(ItemComma)
			return lex2ndAddr
		case r == '!':
			l.emit(ItemExpMark)
			return lex2ndAddrDone
		}
	}
}

// lex2ndAddr lexes the part after a comma was received.
func lex2ndAddr(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == 0:
			l.emit(ItemEOF)
			return nil
		case isSpace(r):
			l.ignore()
		case isNumeric(r):
			l.acceptRun("0123456789")
			l.emit(ItemInt)
			return lex2ndAddrDone
		case isCommand(r):
			l.backup()
			return lex2ndAddrDone
		case r == '$':
			l.emit(ItemDollar)
			return lex2ndAddrDone
		case r == '\\':
			l.ignore()
			delim := l.next()
			l.emit(ItemSlash)
			return lexInsideAddr(delim, lex2ndAddrDone)
		case r == '/':
			l.emit(ItemSlash)
			return lexInsideAddr('/', lex2ndAddrDone)
		}
	}
}

// lexInsideAddr lexes an address component (/xxx/) that is delimited
// by 'div' and when done will return to the onComplete state.
func lexInsideAddr(div rune, onComplete stateFn) stateFn {
	return func(l *Lexer) stateFn {
		for {
			switch r := l.next(); {
			case r == div:
				l.backup()
				l.emit(ItemLit)
				l.next()
				l.emit(ItemSlash)
				return onComplete
			case r == '\\':
				// Escape the div
				if l.peek() == div {
					l.escapePrev()
					l.next()
				}
				// l.escapePrev()
				// r = l.next()
				// if r != div {
				// 	l.errorf("unrecognised escape character %c", r)
				// 	return nil
				// }
			case r == 0:
				l.emit(ItemEOF)
				return nil
			}
		}
	}
}

// lex2ndAddrDone lexes the portion that is after the 2nd address.
// We know that no more addresses can come in
func lex2ndAddrDone(l *Lexer) stateFn {
	l.acceptRun(" \t")
	l.ignore()
	if l.accept("!") {
		l.emit(ItemExpMark)
	}
	r := l.next()
	if !isCommand(r) {
		// r isn't a command but could be a brace
		switch r {
		case '{':
			l.emit(ItemLBrace)
			return lexEnd
		case '}':
			l.emit(ItemRBrace)
			return lexEnd
		default:
			return l.errorf("expected next rune to be a command, got %c", r)
		}
	}
	l.emit(ItemCmd)

	switch r {
	case 's', 'y':
		// Get divider character
		// return clorure to that func
		div := l.next()
		if isSpace(div) || div == '\n' {
			return l.errorf("can not use %c as divider for s/y commands", div)
		}
		l.emit(ItemDiv)
		return parseDivExp(div)
	case 'r', 'w', 'b', 't':
		// get identifier, stop and ; or \n
		l.acceptRun(" ")
		l.ignore()
		if r = l.next(); r == '\n' || r == ';' || r == '0' {
			l.backup()
			return lexEnd
		}
		return lexIdentToEnd
	case 'c', 'i', 'a':
		// commands that can take a backslash
		l.acceptRun(" ")
		l.ignore()
		if l.next() != '\\' {
			return l.errorf("c, i, and a cmds must be followed by \\")
		}
		l.emit(ItemBackslash)
		if l.next() != '\n' {
			return l.errorf("newline expected after \\ for c,i,a cmds")
		}
		l.ignore()
		return lexLiteralLine
	}
	return lexEnd
}

func lexEnd(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == 0:
			l.emit(ItemEOF)
			return nil
		case r == ';':
			l.emit(ItemSemicolon)
			return lexStart
		case r == '\n':
			l.emit(ItemNewline)
			return lexStart
		case isSpace(r):
			l.ignore()
		default:
			return l.errorf("unexpected token after command '%c'", r)
		}
	}
}

// parseDivExp parses the first term, divider, second term, divider.
// For example with 's/123/456/', this func would handle the '123/456/' part.
func parseDivExp(div rune) stateFn {
	return func(l *Lexer) stateFn {
		i := 2
		for {
			switch r := l.next(); {
			case r == 0:
				l.errorf("unexpected EOF while parsing div exp")
			case r == '\\':
				switch r := l.peek(); {
				// Items that the parser wants to escape. If not, defer
				// to regex handler
				case r == div:
					l.escapePrev()
					l.next() // dont have this be seen as sending delimiter
				case r == '\n':
					l.escapePrev()
				case r == '\\':
					l.next()
				}
			case r == div:
				l.backup()
				l.emit(ItemLit)
				l.next()
				l.emit(ItemDiv)
				i--
				if i == 0 {
					// we collected both parts, look for flags
				endSY:
					switch r := l.next(); {
					case r == 0:
						l.emit(ItemEOF)
						return nil
					case r == '\n':
						l.emit(ItemNewline)
						return lexStart
					case r == ';':
						l.emit(ItemSemicolon)
						return lexStart
					case isSpace(r):
						l.ignore()
						goto endSY // this feels so good :)
					case isFlag(r):
						l.emit(ItemIdent)
						l.acceptRun(" ") // optional space before additional arg
						l.ignore()
						switch r = l.next(); {
						case r == 0:
							l.emit(ItemEOF)
							return nil
						case r == '\n':
							// no flags
							l.emit(ItemNewline)
							return lexStart
						default:
							// process until end of line as arg
							for {
								switch r = l.next(); {
								case r == 0:
									l.emit(ItemIdent)
									l.emit(ItemEOF)
									return nil
								case r == '\n':
									l.backup()
									l.emit(ItemIdent)
									l.emit(ItemNewline)
									return lexStart
								}
							}
						}
					default:
						l.emit(ItemError)
						return nil
					}
				}
			}
		}
	}
}

func lexLiteralLine(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == 0:
			l.emit(ItemLit)
			l.emit(ItemEOF)
			return nil
		case r == '\\':
			if l.next() == '\n' {
				l.escapePrev()
			}
		case r == '\n':
			l.backup()
			l.emit(ItemLit)
			l.next()
			l.emit(ItemNewline)
			return lexStart
		default:
		}
	}
}

func lexIdentToEnd(l *Lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == 0:
			l.emit(ItemIdent)
			l.emit(ItemEOF)
			return nil
		case r == ';':
			l.backup()
			l.emit(ItemIdent)
			l.next()
			l.emit(ItemSemicolon)
			return lexStart
		case r == '\n':
			l.backup()
			l.emit(ItemIdent)
			l.next()
			l.emit(ItemNewline)
			return lexStart
		}

	}
}
