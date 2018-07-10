package parser

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zkry/go-sed/ast"
	"github.com/zkry/go-sed/lexer"
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
		{program: "b text", isError: false},
		{program: "b", isError: true},
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
		{program: "t label", isError: false},
		{program: "y/abc/def/", isError: false},
		{program: "n", isError: false},
		{program: "N", isError: false},
		{program: "i\\\ntext", isError: false},
		{program: "x", isError: false},
		{program: "x\\", isError: true},
		{program: "=", isError: false},
		{program: "= =", isError: true},
	}

	for i, test := range tests {
		l := lexer.New(test.program)
		p := New(l)
		_ = p.parseStatement()
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
		ast     *ast.Program
	}{
		{program: "s/one/two/", isError: false, ast: &ast.Program{
			Statements: []ast.Statement{
				&ast.SStmt{
					Addresser:   &ast.BlankAddress{},
					FindAddr:    "one",
					ReplaceAddr: "two",
				},
			},
			Labels: map[string]*ast.Statement{},
		}},
		{program: "s/one/two/;s/two/three/;", isError: false},
		{program: "s/one/two/\ns/two/three/", isError: false},
		{program: "s/one/two/\n\ns/two/three/;", isError: false},
		{program: "\ns/one/two/\n\ns/two/three/\n", isError: false},
		{program: "/quit_now/q", isError: false},
	}
	for i, test := range tests {
		l := lexer.New(test.program)
		p := New(l)
		ast := p.ParseProgram()
		if !test.isError && len(p.errors) > 0 {
			t.Errorf("Program [%d] %s expected no errors but got: %v", i, test.program, p.errors)
		}
		if test.ast != nil {
			if !cmp.Equal(ast, test.ast) {
				t.Errorf("Program [%d] %s ast tree not equal to expected result", i, test.program)
			}
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
	}

	opt := ast.RuntimeOptions{
		AllowExec: true,
		AutoPrint: true,
	}

	for i, tt := range runTests {
		fmt.Printf("\n\n=======================\n")
		fmt.Printf("======= test %d =======\n", i)
		fmt.Printf("=======================\n")
		l := lexer.New(tt.program)
		p := New(l)
		program := p.ParseProgram()
		if len(p.errors) > 0 {
			t.Errorf("Program [%d] %s encountered errors %v", i, tt.program, p.errors)
			continue
		}
		out := ast.Run(program, tt.input, opt)
		if out != tt.output {
			t.Errorf("Program [%d] %s produced incorrect output.\n Expected:\n-----\n%s\n-----\n Got:\n-----\n%s\n-----\n", i, tt.program, tt.output, out)
		}
	}
}
