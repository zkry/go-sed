package parser

import (
	"testing"

	"github.com/kr/pretty"
	"github.com/zkry/go-sed/lexer"
)

func TestStatement(t *testing.T) {
	input := `1,2 s/one/two/g`
	l := lexer.New(input)
	p := New(l)

	stmt := p.parseStatement()
	pretty.Println(stmt)
}
