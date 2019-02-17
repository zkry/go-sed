package lexer

import (
	"testing"
)

var lexerTests = []struct {
	program  string
	expected []Item
}{
	{ // Program
		program: "/",
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 1
		program: "/addr/",
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 2
		program: "/addr1/,/addr2/",
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemError, Value: ""},
		},
	},
	{ // Program 3
		program: "/addr1/,/addr2/d",
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 4
		program: "/addr1/,/addr2/s/find/replace/",
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "find"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "replace"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 5
		program: "/addr1/,/addr2/s/find/replace/g",
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "find"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "replace"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "g"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 6
		program: `/-> addr1 <-/,/!@#$%\/*+/s/some text/~~~~~~/g`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "-> addr1 <-"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: `!@#$%/*+`},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "some text"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "~~~~~~"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "g"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 7
		program: `s/one/two/`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 8
		program: `s/one/two/p`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 9
		program: `y/abc/xyz/`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "y"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "abc"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "xyz"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 10
		program: `/addr/d`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 11
		program: `/addr/ d`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 12
		program: `/addr/     d`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 13
		program: `/addr1/,/addr2/     d`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 14
		program: `/addr1/,/addr2/s/one/two/w outfile.txt`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "w"},
			Item{Type: ItemIdent, Value: "outfile.txt"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 15
		program: `/addr1/,/addr2/s/one/two/w      outfile.txt`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "w"},
			Item{Type: ItemIdent, Value: "outfile.txt"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 16
		program: `/addr1/,/addr2/s/one/two/woutfile.txt`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "w"},
			Item{Type: ItemIdent, Value: "outfile.txt"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 17
		program: `/addr1/,/addr2/s/one/two/woutfile.txt`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "w"},
			Item{Type: ItemIdent, Value: "outfile.txt"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 18
		program: `s/one/two/
	s/three/four/`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "three"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "four"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 19
		program: `s/one/two/p
	s/three/four/
	s/five/six/p`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "p"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "three"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "four"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "five"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "six"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 20
		program: `s/one/two/p
	/addr1/,/addr2/s/three/four/
	s/five/six/p`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "p"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr1"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "addr2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "three"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "four"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "five"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "six"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 21
		program: `$d`,
		expected: []Item{
			Item{Type: ItemDollar, Value: "$"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 22
		program: `5d`,
		expected: []Item{
			Item{Type: ItemInt, Value: "5"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 23
		program: `1,5d`,
		expected: []Item{
			Item{Type: ItemInt, Value: "1"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemInt, Value: "5"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 24
		program: `5,$d`,
		expected: []Item{
			Item{Type: ItemInt, Value: "5"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemDollar, Value: "$"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 25
		program: `5,$  d`,
		expected: []Item{
			Item{Type: ItemInt, Value: "5"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemDollar, Value: "$"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 26
		program: `5,$  d
	1,2d
	3,4d
	s/a/b/p`,
		expected: []Item{
			Item{Type: ItemInt, Value: "5"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemDollar, Value: "$"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemInt, Value: "1"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemInt, Value: "2"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemInt, Value: "3"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemInt, Value: "4"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "a"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "b"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 27
		program: `s|a|b|`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemLit, Value: "a"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemLit, Value: "b"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 28
		program: `s|a|b|p`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemLit, Value: "a"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemLit, Value: "b"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemIdent, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 29
		program: `s,a,b,r file.io`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: ","},
			Item{Type: ItemLit, Value: "a"},
			Item{Type: ItemDiv, Value: ","},
			Item{Type: ItemLit, Value: "b"},
			Item{Type: ItemDiv, Value: ","},
			Item{Type: ItemIdent, Value: "r"},
			Item{Type: ItemIdent, Value: "file.io"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 30
		program: `100,/funny/s,a,b,p`,
		expected: []Item{
			Item{Type: ItemInt, Value: "100"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "funny"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: ","},
			Item{Type: ItemLit, Value: "a"},
			Item{Type: ItemDiv, Value: ","},
			Item{Type: ItemLit, Value: "b"},
			Item{Type: ItemDiv, Value: ","},
			Item{Type: ItemIdent, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 31
		program: `s/delete me//`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "delete me"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 32
		program: `s///`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 33
		program: `s////`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemError, Value: "/"},
		},
	},
	{ // Program 34
		program: `$,$,`,
		expected: []Item{
			Item{Type: ItemDollar, Value: "$"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemDollar, Value: "$"},
			Item{Type: ItemError, Value: ","},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 35
		program: `/WORD/ i\
Add this line before every line with WORD`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "WORD"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "i"},
			Item{Type: ItemBackslash, Value: "\\"},
			Item{Type: ItemLit, Value: "Add this line before every line with WORD"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 36
		program: `/WORD/ c\
Replace the current line with the line`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "WORD"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "c"},
			Item{Type: ItemBackslash, Value: "\\"},
			Item{Type: ItemLit, Value: "Replace the current line with the line"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 37
		program: `
	s/blank/lines/`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "blank"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "lines"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 38
		program: `# This is a comment
	s/blank/lines/`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "#"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "blank"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "lines"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 39
		program: `    # This is a comment
	s/blank/lines/`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "#"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "blank"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "lines"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 40
		program: `3 s/[0-9][0-9]*//`,
		expected: []Item{
			Item{Type: ItemInt, Value: "3"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "[0-9][0-9]*"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 41
		program: `/^#/ s/[0-9][0-9]*//`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "^#"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "[0-9][0-9]*"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 42
		program: `/^#/ s/[0-9][0-9]*//`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "^#"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "[0-9][0-9]*"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 43
		program: `\_/usr/local/bin_ s_/usr/local_/common/all_`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "_"},
			Item{Type: ItemLit, Value: "/usr/local/bin"},
			Item{Type: ItemSlash, Value: "_"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "_"},
			Item{Type: ItemLit, Value: "/usr/local"},
			Item{Type: ItemDiv, Value: "_"},
			Item{Type: ItemLit, Value: "/common/all"},
			Item{Type: ItemDiv, Value: "_"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 44
		program: `/^g/ s_g_s_g`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "^g"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "_"},
			Item{Type: ItemLit, Value: "g"},
			Item{Type: ItemDiv, Value: "_"},
			Item{Type: ItemLit, Value: "s"},
			Item{Type: ItemDiv, Value: "_"},
			Item{Type: ItemIdent, Value: "g"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 45
		program: `1,100 s/A/a/`,
		expected: []Item{
			Item{Type: ItemInt, Value: "1"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemInt, Value: "100"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "A"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "a"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `p
	p
	p`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: "d",
		expected: []Item{
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: " p",
		expected: []Item{
			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: "\tp",
		expected: []Item{
			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
		/begin/n
		s/old/new/`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "begin"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "n"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "old"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "new"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `# Testing Grouping
	/begin/,/end/ {
	s/#.*//
		s/[ ^I]*$//
		/^$/ d
		p
	}`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "#"},

			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "begin"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "end"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "#.*"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "[ ^I]*$"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "^$"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemRBrace, Value: "}"},

			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
		1,100 {
			/begin/,/end/ {
			     s/#.*//
			     s/[ ^I]*$//
			     /^$/ d
			     p
			}
		}`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemInt, Value: "1"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemInt, Value: "100"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "begin"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "end"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "#.*"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "[ ^I]*$"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "^$"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemRBrace, Value: "}"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemRBrace, Value: "}"},

			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
		1,100 !{
			/begin/,/end/ !{
			     s/#.*//
			     s/[ ^I]*$//
			     /^$/ d
			     p
			}
		}`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemInt, Value: "1"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemInt, Value: "100"},
			Item{Type: ItemExpMark, Value: "!"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "begin"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "end"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemExpMark, Value: "!"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "#.*"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "[ ^I]*$"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "^$"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemRBrace, Value: "}"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemRBrace, Value: "}"},

			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
		1,100!{
			/begin/,/end/ !{
				/begin/n
				s/old/new/
			}
		}`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemInt, Value: "1"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemInt, Value: "100"},
			Item{Type: ItemExpMark, Value: "!"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"}, // 7
			Item{Type: ItemLit, Value: "begin"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "end"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemExpMark, Value: "!"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"}, // 17
			Item{Type: ItemLit, Value: "begin"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "n"},
			Item{Type: ItemNewline, Value: "\n"}, // 21

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "old"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "new"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemRBrace, Value: "}"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemRBrace, Value: "}"},

			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
	bx
	`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "b"},
			Item{Type: ItemIdent, Value: "x"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
	/^$/ bpara
	H
	$ bpara
	b
	:para
	x
	/'$1'/ p
	`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"}, // 7
			Item{Type: ItemLit, Value: "^$"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "b"},
			Item{Type: ItemIdent, Value: "para"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "H"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemDollar, Value: "$"},
			Item{Type: ItemCmd, Value: "b"},
			Item{Type: ItemIdent, Value: "para"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "b"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemColon, Value: ":"},
			Item{Type: ItemIdent, Value: "para"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "x"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"}, // 7
			Item{Type: ItemLit, Value: "'$1'"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
	:again
		s/([ ^I]*)//
		tagain
	`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemColon, Value: ":"},
			Item{Type: ItemIdent, Value: "again"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "([ ^I]*)"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ""},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemCmd, Value: "t"},
			Item{Type: ItemIdent, Value: "again"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
	/grep/ !{;H;x;s/^.*\n\(.*\n.*\)$/\1/;x;}
	/grep/ {;H;n;H;x;p;a\
---
	}
	`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "grep"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemExpMark, Value: "!"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "H"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "x"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "^.*\\n\\(.*\\n.*\\)$"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "\\1"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "x"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemRBrace, Value: "}"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "grep"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "H"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "n"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "H"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "x"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "p"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "a"},
			Item{Type: ItemBackslash, Value: "\\"},
			Item{Type: ItemLit, Value: "---"},
			Item{Type: ItemNewline, Value: "\n"},

			Item{Type: ItemRBrace, Value: "}"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
	a \
---
	`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemCmd, Value: "a"},
			Item{Type: ItemBackslash, Value: "\\"},
			Item{Type: ItemLit, Value: "---"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `s|a|b|g|`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemLit, Value: "a"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemLit, Value: "b"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemIdent, Value: "g"},
			Item{Type: ItemIdent, Value: "|"}, // Let the parser handle this...
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 47
		program: `s|a|b|gp`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemLit, Value: "a"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemLit, Value: "b"},
			Item{Type: ItemDiv, Value: "|"},
			Item{Type: ItemIdent, Value: "g"},
			Item{Type: ItemIdent, Value: "p"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 48
		program: `s/one/two/;s/two/three/;`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "three"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 49
		program: "s/one/two/\ns/two/three/",
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "two"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "three"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{ // Program 48
		program: `s/one/two;`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "one"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemError, Value: "two;"},
		},
	},
	{
		program: `$s/$/}/

	/./!d`,
		expected: []Item{
			Item{Type: ItemDollar, Value: "$"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "$"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "}"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "."},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemExpMark, Value: "!"},
			Item{Type: ItemCmd, Value: "d"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `/a/b branch`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "a"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "b"},
			Item{Type: ItemIdent, Value: "branch"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `s/.*/\
|--------|\
|        |\
|        |\
|--------|/`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ".*"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: "\n|--------|\n|        |\n|        |\n|--------|"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `
	  /^t3$/{ s/.*/\
	 TEST 3 - 3\
	      _____________ \
	     |     ==      |\
	     |     ==      |\
	     |    ==  =    |\
	     |     = ==    |\
	     |  =o         |\
	     |  ==         |\
	     |             |\
	     |.____________|\
/ ; b endmap
	  }`,
		expected: []Item{
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "^t3$"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLBrace, Value: "{"},
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: ".*"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: `
	 TEST 3 - 3
	      _____________ 
	     |     ==      |
	     |     ==      |
	     |    ==  =    |
	     |     = ==    |
	     |  =o         |
	     |  ==         |
	     |             |
	     |.____________|
`},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemSemicolon, Value: ";"},
			Item{Type: ItemCmd, Value: "b"},
			Item{Type: ItemIdent, Value: "endmap"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemRBrace, Value: "}"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `s/\\\\/\\/g`,
		expected: []Item{
			Item{Type: ItemCmd, Value: "s"},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: `\\\\`},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemLit, Value: `\\`},
			Item{Type: ItemDiv, Value: "/"},
			Item{Type: ItemIdent, Value: "g"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: `\s\ss,\s\sssx\xx\x\xxwsss`,
		expected: []Item{
			Item{Type: ItemSlash, Value: "s"},
			Item{Type: ItemLit, Value: "s"},
			Item{Type: ItemSlash, Value: "s"},
			Item{Type: ItemComma, Value: ","},
			Item{Type: ItemSlash, Value: "s"},
			Item{Type: ItemLit, Value: "s"},
			Item{Type: ItemSlash, Value: "s"},

			Item{Type: ItemCmd, Value: "s"},

			Item{Type: ItemDiv, Value: "x"},
			Item{Type: ItemLit, Value: "x"},
			Item{Type: ItemDiv, Value: "x"},
			Item{Type: ItemLit, Value: "xx"},
			Item{Type: ItemDiv, Value: "x"},

			Item{Type: ItemIdent, Value: "w"},
			Item{Type: ItemIdent, Value: "sss"},

			Item{Type: ItemEOF, Value: ""},
		},
	},
	{
		program: "N\n/1.*\n.*2/P",
		expected: []Item{
			Item{Type: ItemCmd, Value: "N"},
			Item{Type: ItemNewline, Value: "\n"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemLit, Value: "1.*\n.*2"},
			Item{Type: ItemSlash, Value: "/"},
			Item{Type: ItemCmd, Value: "P"},
			Item{Type: ItemEOF, Value: ""},
		},
	},
}

func TestNextTokens(t *testing.T) {
next_test:
	for i, lt := range lexerTests {
		_, items := New(lt.program)
		for j, et := range lt.expected {

			gotTok := <-items

			if gotTok.Type != et.Type {
				t.Errorf("Program[%d]:%s line[%d] - tokentype wrong. expected=%v, got=%v:%v", i, lt.program, j, et.Type, gotTok.Type, gotTok.Value)
			}

			if gotTok.Type == ItemError || gotTok.Type == ItemEOF {
				continue next_test
			}

			if gotTok.Value != et.Value {
				t.Fatalf("Program[%d]:%s line[%d] - tokenliteral wrong. expected='%v', got='%v'", i, lt.program, j, et.Value, gotTok.Value)
			}
		}
	}
}
