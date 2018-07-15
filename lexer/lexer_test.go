package lexer

import (
	"testing"

	"github.com/zkry/go-sed/token"
)

func TestNextToken(t *testing.T) {
	for i, lt := range lexerTests {
		l := New(lt.program)
		for j, et := range lt.expected {
			gotTok := l.NextToken()

			if gotTok.Type != et.Type {
				t.Fatalf("Program[%d]:%s line[%d] - tokentype wrong. expected=%v, got=%v", i, lt.program, j, et.Type, gotTok.Type)
			}

			if gotTok.Literal != et.Literal {
				t.Fatalf("Program[%d]:%s line[%d] - tokenliteral wrong. expected=%v, got=%v", i, lt.program, j, et.Literal, gotTok.Literal)
			}
		}
	}
}

var lexerTests = []struct {
	program  string
	expected []token.Token
}{
	{ // Program 0
		program: "/",
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 1
		program: "/addr/",
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr"},
			token.Token{token.SLASH, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 2
		program: "/addr1/,/addr2/",
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 3
		program: "/addr1/,/addr2/d",
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 4
		program: "/addr1/,/addr2/s/find/replace/",
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "find"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "replace"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 5
		program: "/addr1/,/addr2/s/find/replace/g",
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "find"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "replace"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "g"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 6
		program: `/-> addr1 <-/,/!@#$%\/*+/s/some text/~~~~~~/g`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "-> addr1 <-"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, `!@#$%\/*+`},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "some text"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "~~~~~~"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "g"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 7
		program: `s/one/two/`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 8
		program: `s/one/two/p`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "p"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 9
		program: `y/abc/xyz/`,
		expected: []token.Token{
			token.Token{token.CMD, "y"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "abc"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "xyz"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 10
		program: `/addr/d`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 11
		program: `/addr/ d`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 12
		program: `/addr/     d`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 13
		program: `/addr1/,/addr2/     d`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 14
		program: `/addr1/,/addr2/s/one/two/w outfile.txt`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "w"},
			token.Token{token.IDENT, "outfile.txt"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 15
		program: `/addr1/,/addr2/s/one/two/w      outfile.txt`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "w"},
			token.Token{token.IDENT, "outfile.txt"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 16
		program: `/addr1/,/addr2/s/one/two/woutfile.txt`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "w"},
			token.Token{token.IDENT, "outfile.txt"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 17
		program: `/addr1/,/addr2/s/one/two/woutfile.txt`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "w"},
			token.Token{token.IDENT, "outfile.txt"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 18
		program: `s/one/two/
s/three/four/`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "three"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "four"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 19
		program: `s/one/two/p
s/three/four/
s/five/six/p`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "p"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "three"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "four"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "five"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "six"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "p"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 20
		program: `s/one/two/p
/addr1/,/addr2/s/three/four/
s/five/six/p`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "p"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr1"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "addr2"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "three"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "four"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "five"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "six"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "p"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 21
		program: `$d`,
		expected: []token.Token{
			token.Token{token.DOLLAR, "$"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 22
		program: `5d`,
		expected: []token.Token{
			token.Token{token.INT, "5"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 23
		program: `1,5d`,
		expected: []token.Token{
			token.Token{token.INT, "1"},
			token.Token{token.COMMA, ","},
			token.Token{token.INT, "5"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 24
		program: `5,$d`,
		expected: []token.Token{
			token.Token{token.INT, "5"},
			token.Token{token.COMMA, ","},
			token.Token{token.DOLLAR, "$"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 25
		program: `5,$  d`,
		expected: []token.Token{
			token.Token{token.INT, "5"},
			token.Token{token.COMMA, ","},
			token.Token{token.DOLLAR, "$"},
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 26
		program: `5,$  d
1,2d
3,4d
s/a/b/p`,
		expected: []token.Token{
			token.Token{token.INT, "5"},
			token.Token{token.COMMA, ","},
			token.Token{token.DOLLAR, "$"},
			token.Token{token.CMD, "d"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.INT, "1"},
			token.Token{token.COMMA, ","},
			token.Token{token.INT, "2"},
			token.Token{token.CMD, "d"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.INT, "3"},
			token.Token{token.COMMA, ","},
			token.Token{token.INT, "4"},
			token.Token{token.CMD, "d"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "a"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "b"},
			token.Token{token.DIV, "/"},
			token.Token{token.IDENT, "p"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 27
		program: `s|a|b|`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "|"},
			token.Token{token.LIT, "a"},
			token.Token{token.DIV, "|"},
			token.Token{token.LIT, "b"},
			token.Token{token.DIV, "|"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 28
		program: `s|a|b|p`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "|"},
			token.Token{token.LIT, "a"},
			token.Token{token.DIV, "|"},
			token.Token{token.LIT, "b"},
			token.Token{token.DIV, "|"},
			token.Token{token.IDENT, "p"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 29
		program: `s,a,b,r file.io`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, ","},
			token.Token{token.LIT, "a"},
			token.Token{token.DIV, ","},
			token.Token{token.LIT, "b"},
			token.Token{token.DIV, ","},
			token.Token{token.IDENT, "r"},
			token.Token{token.IDENT, "file.io"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 30
		program: `100,/funny/s,a,b,b`,
		expected: []token.Token{
			token.Token{token.INT, "100"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "funny"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, ","},
			token.Token{token.LIT, "a"},
			token.Token{token.DIV, ","},
			token.Token{token.LIT, "b"},
			token.Token{token.DIV, ","},
			token.Token{token.IDENT, "b"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 31
		program: `s/delete me//`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "delete me"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 32
		program: `s///`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 33
		program: `s////`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.ILLEGAL, "/"},
		},
	},
	{ // Program 34
		program: `$,$,`,
		expected: []token.Token{
			token.Token{token.DOLLAR, "$"},
			token.Token{token.COMMA, ","},
			token.Token{token.DOLLAR, "$"},
			token.Token{token.ILLEGAL, ","},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 35
		program: `/WORD/ i\
Add this line before every line with WORD`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "WORD"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "i"},
			token.Token{token.BACKSLASH, "\\"},
			token.Token{token.LIT, "Add this line before every line with WORD"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 36
		program: `/WORD/ c\
Replace the current line with the line`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "WORD"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "c"},
			token.Token{token.BACKSLASH, "\\"},
			token.Token{token.LIT, "Replace the current line with the line"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 37
		program: `
s/blank/lines/`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "blank"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "lines"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 38
		program: `# This is a comment
s/blank/lines/`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "#"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "blank"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "lines"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 39
		program: `    # This is a comment
s/blank/lines/`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "#"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "blank"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "lines"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 40
		program: `3 s/[0-9][0-9]*//`,
		expected: []token.Token{
			token.Token{token.INT, "3"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "[0-9][0-9]*"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 41
		program: `/^#/ s/[0-9][0-9]*//`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "^#"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "[0-9][0-9]*"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 42
		program: `/^#/ s/[0-9][0-9]*//`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "^#"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "[0-9][0-9]*"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 43
		program: `\_/usr/local/bin_ s_/usr/local_/common/all_`,
		expected: []token.Token{
			token.Token{token.SLASH, "_"},
			token.Token{token.LIT, "/usr/local/bin"},
			token.Token{token.SLASH, "_"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "_"},
			token.Token{token.LIT, "/usr/local"},
			token.Token{token.DIV, "_"},
			token.Token{token.LIT, "/common/all"},
			token.Token{token.DIV, "_"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 44
		program: `/^g/ s_g_s_g`,
		expected: []token.Token{
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "^g"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "_"},
			token.Token{token.LIT, "g"},
			token.Token{token.DIV, "_"},
			token.Token{token.LIT, "s"},
			token.Token{token.DIV, "_"},
			token.Token{token.IDENT, "g"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `1,100 s/A/a/`,
		expected: []token.Token{
			token.Token{token.INT, "1"},
			token.Token{token.COMMA, ","},
			token.Token{token.INT, "100"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "A"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "a"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `p
p
p`,
		expected: []token.Token{
			token.Token{token.CMD, "p"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.CMD, "p"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.CMD, "p"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: "d",
		expected: []token.Token{
			token.Token{token.CMD, "d"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: " p",
		expected: []token.Token{
			token.Token{token.CMD, "p"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: "\tp",
		expected: []token.Token{
			token.Token{token.CMD, "p"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `
	/begin/n
	s/old/new/`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "begin"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "n"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "old"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "new"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `# Testing Grouping
/begin/,/end/ {
s/#.*//
	s/[ ^I]*$//
	/^$/ d
	p
}`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "#"},

			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "begin"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "end"},
			token.Token{token.SLASH, "/"},
			token.Token{token.LBRACE, "{"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "#.*"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "[ ^I]*$"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "^$"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "d"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "p"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.RBRACE, "}"},

			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `
	1,100 {
		/begin/,/end/ {
		     s/#.*//
		     s/[ ^I]*$//
		     /^$/ d
		     p
		}
	}`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.INT, "1"},
			token.Token{token.COMMA, ","},
			token.Token{token.INT, "100"},
			token.Token{token.LBRACE, "{"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "begin"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "end"},
			token.Token{token.SLASH, "/"},
			token.Token{token.LBRACE, "{"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "#.*"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "[ ^I]*$"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "^$"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "d"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "p"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.RBRACE, "}"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.RBRACE, "}"},

			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `
	1,100 !{
		/begin/,/end/ !{
		     s/#.*//
		     s/[ ^I]*$//
		     /^$/ d
		     p
		}
	}`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.INT, "1"},
			token.Token{token.COMMA, ","},
			token.Token{token.INT, "100"},
			token.Token{token.EXPLMARK, "!"},
			token.Token{token.LBRACE, "{"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "begin"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "end"},
			token.Token{token.SLASH, "/"},
			token.Token{token.EXPLMARK, "!"},
			token.Token{token.LBRACE, "{"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "#.*"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "[ ^I]*$"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "^$"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "d"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "p"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.RBRACE, "}"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.RBRACE, "}"},

			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `
	1,100!{
		/begin/,/end/ !{
			/begin/n
			s/old/new/
		}
	}`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.INT, "1"},
			token.Token{token.COMMA, ","},
			token.Token{token.INT, "100"},
			token.Token{token.EXPLMARK, "!"},
			token.Token{token.LBRACE, "{"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"}, // 7
			token.Token{token.LIT, "begin"},
			token.Token{token.SLASH, "/"},
			token.Token{token.COMMA, ","},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "end"},
			token.Token{token.SLASH, "/"},
			token.Token{token.EXPLMARK, "!"},
			token.Token{token.LBRACE, "{"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"}, // 17
			token.Token{token.LIT, "begin"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "n"},
			token.Token{token.NEWLINE, "\n"}, // 21

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "old"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "new"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.RBRACE, "}"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.RBRACE, "}"},

			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `
/^$/ bpara
H
$ bpara
b
:para
x
/'$1'/ p
`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"}, // 7
			token.Token{token.LIT, "^$"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "b"},
			token.Token{token.IDENT, "para"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "H"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.DOLLAR, "$"},
			token.Token{token.CMD, "b"},
			token.Token{token.IDENT, "para"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "b"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.COLON, ":"},
			token.Token{token.IDENT, "para"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "x"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"}, // 7
			token.Token{token.LIT, "'$1'"},
			token.Token{token.SLASH, "/"},
			token.Token{token.CMD, "p"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `
:again
	s/([ ^I]*)//
	tagain
`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.COLON, ":"},
			token.Token{token.IDENT, "again"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "([ ^I]*)"},
			token.Token{token.DIV, "/"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.CMD, "t"},
			token.Token{token.IDENT, "again"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `
/grep/ !{;H;x;s/^.*\n\(.*\n.*\)$/\1/;x;}
/grep/ {;H;n;H;x;p;a\
---
}
`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "grep"},
			token.Token{token.SLASH, "/"},
			token.Token{token.EXPLMARK, "!"},
			token.Token{token.LBRACE, "{"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "H"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "x"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "^.*\\n\\(.*\\n.*\\)$"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "\\1"},
			token.Token{token.DIV, "/"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "x"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.RBRACE, "}"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.SLASH, "/"},
			token.Token{token.LIT, "grep"},
			token.Token{token.SLASH, "/"},
			token.Token{token.LBRACE, "{"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "H"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "n"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "H"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "x"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "p"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "a"},
			token.Token{token.BACKSLASH, "\\"},
			token.Token{token.LIT, "---"},
			token.Token{token.NEWLINE, "\n"},

			token.Token{token.RBRACE, "}"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 45
		program: `
a \
---
`,
		expected: []token.Token{
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.CMD, "a"},
			token.Token{token.BACKSLASH, "\\"},
			token.Token{token.LIT, "---"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 46
		program: `s|a|b|g|`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "|"},
			token.Token{token.LIT, "a"},
			token.Token{token.DIV, "|"},
			token.Token{token.LIT, "b"},
			token.Token{token.DIV, "|"},
			token.Token{token.IDENT, "g"},
			token.Token{token.ILLEGAL, "|"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 47
		program: `s|a|b|gp`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "|"},
			token.Token{token.LIT, "a"},
			token.Token{token.DIV, "|"},
			token.Token{token.LIT, "b"},
			token.Token{token.DIV, "|"},
			token.Token{token.IDENT, "g"},
			token.Token{token.IDENT, "p"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 48
		program: `s/one/two/;s/two/three/;`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "three"},
			token.Token{token.DIV, "/"},
			token.Token{token.SEMICOLON, ";"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 49
		program: "s/one/two/\ns/two/three/",
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.NEWLINE, "\n"},
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "three"},
			token.Token{token.DIV, "/"},
			token.Token{token.EOF, ""},
		},
	},
	{ // Program 48
		program: `s/one/two;`,
		expected: []token.Token{
			token.Token{token.CMD, "s"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "one"},
			token.Token{token.DIV, "/"},
			token.Token{token.LIT, "two;"},
			token.Token{token.EOF, ""},
		},
	},
}
