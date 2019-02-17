package ast

import (
	"fmt"
	"testing"
)

func TestStatement(t *testing.T) {
	tests := []struct {
		program string
		isError bool
	}{
		{program: "s/1/2//", isError: true},
		{program: "/adsf/s//2//", isError: true},
		{program: "/adsfs//2//", isError: true},
		{program: "s/1/2/", isError: false},
		{program: "s//2/", isError: false},
		{program: "s/2//", isError: false},
		{program: "a\\\ntext", isError: false},
		{program: "a text", isError: true},
		{program: "btext", isError: false},
		{program: "b", isError: false},
		{program: "c\\\ntext", isError: false},
		{program: "/addr/d", isError: false},
		{program: "/a/,/b/ d", isError: false},
		{program: "$d", isError: false},
		{program: "1,5 d", isError: false},
		{program: "1,5,5 d", isError: true},
		{program: "D", isError: false},
		{program: "p", isError: false},
		{program: "1p", isError: false},
		{program: "1P", isError: false},
		{program: "h", isError: false},
		{program: "H", isError: false},
		{program: "/what/q", isError: false},
		{program: "tlabel", isError: false},
		{program: "y/abc/def/", isError: false},
		{program: "n", isError: false},
		{program: "N", isError: false},
		{program: "i\\\ntext", isError: false},
		{program: "x", isError: false},
		{program: "x\\", isError: true},
		{program: "=", isError: false},
		{program: "= =", isError: true},
		{program: ":label1", isError: false},
	}

	for i, test := range tests {
		p := New(test.program)
		_, _ = p.parseStatement()
		if len(p.Errors()) > 0 && !test.isError {
			t.Errorf("Stmt [%d] %s failed: expected no error, got %v\n", i, test.program, p.Errors())
		} else if len(p.Errors()) == 0 && test.isError {
			t.Errorf("Stmt [%d] %s failed: Expected an error and got none.", i, test.program)
		}
	}
	// l := lexer.New(input)
	// p := New(l)

	// stmt := p.parseStatement()
	// pretty.Println(stmt)
}

func TestParse(t *testing.T) {
	tests := []struct {
		program string
		isError bool
		ast     *Program
	}{
		{program: "s/one/two/", isError: false, ast: &Program{
			Statements: []statement{
				&sStmt{
					addresser:   &blankAddress{},
					FindAddr:    "one",
					ReplaceAddr: "two",
				},
			},
			Labels: map[string]int{},
		}},
		{program: "s/one/two/;s/two/three/;", isError: false},
		{program: "s/one/two/\ns/two/three/", isError: false},
		{program: "s/one/two/\n\ns/two/three/;", isError: false},
		{program: "\ns/one/two/\n\ns/two/three/\n", isError: false},
		{program: "/quit_now/q", isError: false},
	}
	for i, test := range tests {
		p := New(test.program)
		_ = p.ParseProgram()
		if !test.isError && len(p.errors) > 0 {
			t.Errorf("Program [%d] %s expected no errors but got: %v", i, test.program, p.errors)
		}
	}
}

func TestRun(t *testing.T) {
	runTests := []struct {
		program string
		input   string
		output  string
	}{
		{
			program: "p",
			input:   "hello",
			output:  "hello\nhello",
		},
		{
			program: "p",
			input:   "hello\nworld",
			output:  "hello\nhello\nworld\nworld",
		},
		{
			program: "a\\\nXXX",
			input:   "1\n2\n3",
			output:  "1\nXXX\n2\nXXX\n3\nXXX",
		},
		{
			program: "i\\\nXXX",
			input:   "1\n2\n3",
			output:  "XXX\n1\nXXX\n2\nXXX\n3",
		},
		{
			program: "a\\\nafter\ni\\\ninsert",
			input:   "1\n2\n3",
			output:  "insert\n1\nafter\ninsert\n2\nafter\ninsert\n3\nafter",
		},
		{
			program: "a\\\nafter\ni\\\ninsert",
			input:   "1\n2\n3",
			output:  "insert\n1\nafter\ninsert\n2\nafter\ninsert\n3\nafter",
		},
		{
			program: "a\\\na1\na\\\na2\na\\\na3",
			input:   "1\n2\n3",
			output:  "1\na1\na2\na3\n2\na1\na2\na3\n3\na1\na2\na3",
		},
		{
			program: "d",
			input:   "line1\nline2\nline3\nline4",
			output:  "",
		},
		{
			program: "n",
			input:   "line1\nline2\nline3\nline4",
			output:  "line1\nline2\nline3\nline4",
		},
		{
			program: "=",
			input:   "hello\nworld",
			output:  "1\nhello\n2\nworld",
		},
		{
			program: "s/a/b/",
			input:   "This is a word.",
			output:  "This is b word.",
		},
		{ // TODO: Change to traditional Sed syntax
			program: "s/This is a (.*)\\./$1/",
			input:   "This is a word.",
			output:  "word",
		},
		{
			program: "s/a/b/",
			input:   "aaaaa",
			output:  "baaaa",
		},
		{
			program: "s/a/b/g",
			input:   "aaaaa",
			output:  "bbbbb",
		},
		{
			program: "s/a/b/2",
			input:   "aaaaa",
			output:  "abaaa",
		},
		{
			program: "s/a/b/;s/This/That/;s:word::;",
			input:   "This is a word.",
			output:  "That is b .",
		},
		{
			program: "q",
			input:   "a\nb\nc\nd\ne\nf\ng\nh",
			output:  "a",
		},
		{
			program: "/e/q",
			input:   "a\nb\nc\nd\ne\nf\ng\nh",
			output:  "a\nb\nc\nd\ne",
		},
		{
			program: "3q",
			input:   "a\nb\nc\nd\ne\nf\ng\nh",
			output:  "a\nb\nc",
		},
		{
			program: "$s/h/-/",
			input:   "a\nb\nc\nd\ne\nf\ng\nh",
			output:  "a\nb\nc\nd\ne\nf\ng\n-",
		},
		{
			program: "s/-/X/g\n/one/,/two/s/.*//",
			input:   "---\n---\none\n+++\n+++\ntwo\n---\n---",
			output:  "XXX\nXXX\n\n\n\n\nXXX\nXXX",
		},
		{
			program: `
/here/ {
	s/here/HERE/
	s/E/X/g
}`,
			input:  "---\nhere1\n---\nhere2",
			output: "---\nHXRX1\n---\nHXRX2",
		},
		{
			program: `s/here/HERE/p`,
			input:   "here",
			output:  "HERE\nHERE",
		},
		{
			program: `
/here/ {
	s/here/HERE/p
	s/E/X/gp
}`,
			input:  "---\nhere1\n---\nhere2",
			output: "---\nHERE1\nHXRX1\nHXRX1\n---\nHERE2\nHXRX2\nHXRX2",
		},
		{
			program: `
:label1
/xxxxxxx/blabel2
s/x/xx/p
blabel1
:label2
s/x/=/g
p
`,
			input:  "x",
			output: "xx\nxxx\nxxxx\nxxxxx\nxxxxxx\nxxxxxxx\n=======\n=======",
		},
		{
			program: `G`,
			input:   "one\ntwo\nthree\nfour",
			output:  "one\n\ntwo\n\nthree\n\nfour\n",
		},
		{
			program: `\xtwoxd`,
			input:   "one\n\ntwo\n\n\nthree\n\n\n\nfour\n\n\n\n\nend",
			output:  "one\n\n\n\nthree\n\n\n\nfour\n\n\n\n\nend",
		},
		{
			program: `
N
/one\ntwo/s/one/ONE/`,
			input:  "one\ntwo\nthree\nfour",
			output: "ONE\ntwo\nthree\nfour",
		},
		{
			program: `
/^$/ {
	N
	/^\n$/D
}
`,
			input:  "one\n\ntwo\n\n\nthree\n\n\n\nfour\n\n\n\n\nend",
			output: "one\n\ntwo\n\nthree\n\nfour\n\nend",
		},
		{
			program: `
N
/1.*\n.*2/P
D
`,
			input:  "line1\nline2\nline3\nline4",
			output: "line1",
		},
		{
			program: `
:beginning
s/1/ one/
s/on/ON/
p
tbeginning
s/here/HERE/
`,
			input:  "here1\nhere2\nhere3\nhere4\nThis is the end.",
			output: "here ONe\nhere ONe\nHERE ONe\nhere2\nHERE2\nhere3\nHERE3\nhere4\nHERE4\nThis is the end.\nThis is the end.",
		},
		{
			program: `H;G`,
			input:   "here1\nhere2",
			output:  "here1\n\nhere1\nhere2\n\nhere1\nhere2",
		},
		{
			program: "c\\\nCHANGE",
			input:   "here1\nhere2\nhere3\nhere4\nhere5",
			output:  "CHANGE\nCHANGE\nCHANGE\nCHANGE\nCHANGE",
		},
		{
			program: "/START/,/END/c\\\nCHANGE",
			input:   "START\nhere2\nhere3\nhere4\nEND",
			output:  "CHANGE",
		},
	}

	opt := RuntimeOptions{
		AllowExec: true,
		AutoPrint: true,
	}

	for i, tt := range runTests {
		fmt.Printf("\n\n=======================\n")
		fmt.Printf("======= test %d =======\n", i)
		fmt.Printf("=======================\n")
		fmt.Println("program:", tt.program)
		fmt.Println("-----------------------")

		p := New(tt.program)
		program := p.ParseProgram()
		if len(p.errors) > 0 {
			fmt.Printf("ERROR: %v\n", p.errors)
			t.Errorf("Program [%d] %s encountered errors %v", i, tt.program, p.errors)
			continue
		}
		out := program.Run(tt.input, opt)
		if out != tt.output {
			t.Errorf("Program [%d] %s produced incorrect output.\n Expected:\n-----\n%s\n-----\n Got:\n-----\n%s\n-----\n", i, tt.program, tt.output, out)
		}
	}
}
