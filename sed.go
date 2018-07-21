package gosed

import (
	"errors"
	"fmt"

	"github.com/zkry/go-sed/ast"
	"github.com/zkry/go-sed/lexer"
)

type Options struct {
	SupressOutput bool // Prevents program from automatically outputing line.
	AppendFile    bool // Makes the w command append to file.
	ExtendRegexp  bool // Use extended version of regexp
}

type Program struct {
	p   *ast.Program
	opt Options
}

// MustCompile takes a sed script and compiles it into a program.
// Panics if errors are found in script.
func MustCompile(program string, opt Options) *Program {
	l := lexer.New(program)
	p := ast.New(l)
	prg := p.ParseProgram()
	errs := p.Errors()
	if len(errs) > 0 {
		panic("program could not compile: " + fmt.Sprintf("%v", errs))
	}
	return &Program{p: prg, opt: opt}
}

// Compile compiles a sed script and returns a program upon successfull
// compilation. If unsuccessfull errors are returned.
func Compile(program string, opt Options) (*Program, error) {
	l := lexer.New(program)
	p := ast.New(l)
	prg := p.ParseProgram()
	errs := p.Errors()
	if len(errs) > 0 {
		return nil, errors.New("Program didn not compile")
	}
	return &Program{p: prg, opt: opt}, nil
}

func (p *Program) Filter(data []byte) []byte {
	ro := ast.RuntimeOptions{
		AllowExec:  false,
		AutoPrint:  !p.opt.SupressOutput,
		AppendFile: p.opt.AppendFile,
	}
	return []byte(p.p.Run(string(data), ro))
}

func (p *Program) FilterString(data string) string {
	ro := ast.RuntimeOptions{
		AllowExec:  false,
		AutoPrint:  !p.opt.SupressOutput,
		AppendFile: p.opt.AppendFile,
	}
	return p.p.Run(data, ro)
}
