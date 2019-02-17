package gosed

import (
	"fmt"

	"github.com/zkry/go-sed/ast"
	"github.com/zkry/go-sed/lexer"
)

type Options struct {
	SupressOutput     bool // Prevents program from automatically outputing line.
	AppendFile        bool // Makes the w command append to file.
	ExtendRegexp      bool // Use extended version of regexp
	PreviousLinesRead int
}

func (opt *Options) baseRuntimeOptions() ast.RuntimeOptions {
	return ast.RuntimeOptions{
		AllowExec:  false,
		AutoPrint:  !opt.SupressOutput,
		AppendFile: opt.AppendFile,
	}
}

type state struct {
	linesRead int
}

type Program struct {
	p   *ast.Program
	opt Options
	s   state
}

// MustCompile takes a sed script and compiles it into a program.
// Panics if errors are found in script.
func MustCompile(program string, opt Options) *Program {
	p := ast.New(program)
	prg := p.ParseProgram()
	errs := p.Errors()
	if len(errs) > 0 {
		panic("program could not compile: " + fmt.Sprintf("%v", errs))
	}
	return &Program{p: prg, opt: opt}
}

// Compile compiles a sed script and returns a program upon successfull
// compilation. If unsuccessfull errors are returned.
func Compile(program string, opt Options) (*Program, ast.ErrorList) {
	p := ast.New(program)
	prg := p.ParseProgram()
	errs := p.Errors()
	if len(errs) > 0 {
		fmt.Println("Compile: length of error: ", len(p.Errors()))
		return nil, errs
	}
	return &Program{p: prg, opt: opt}, nil
}

func (p *Program) Filter(data []byte) []byte {
	ro := p.opt.baseRuntimeOptions()
	return []byte(p.p.Run(string(data), ro))
}

func (p *Program) FilterString(data string) string {
	ro := p.opt.baseRuntimeOptions()
	return p.p.Run(data, ro)
}

// FilterA performs a normal filter operation but does not reset the state
// after completion. You can repeatedly call FilterA to process input line
// by line.
func (p *Program) FilterA(data []byte) []byte {
	ro := p.opt.baseRuntimeOptions()
	ro.LineNoStart = p.s.linesRead + 1
	res := []byte(p.p.Run(string(data), ro))
	p.s.linesRead += countLines(string(data)) // TODO: Think of more elegant way to do this.
	return res
}

// FilterStringA performs a normal filter operation but does not reset the state
// after completion. Operation is performed on string and returns a string.
// You can repeatedly call FilterA to process input line by line.
func (p *Program) FilterStringA(data string) string {
	ro := p.opt.baseRuntimeOptions()
	ro.LineNoStart = p.s.linesRead + 1
	res := p.p.Run(data, ro)
	p.s.linesRead += countLines(data) // TODO: Think of more elegant way to do this.
	return res
}

func countLines(d string) int {
	var ct int
	for _, r := range d {
		if r == '\n' {
			ct++
		}
	}
	return ct
}

func Info(program string) []lexer.Item {
	p := ast.New(program)
	prg := p.ParseProgram()
	return prg.Tokens
}
