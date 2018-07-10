package ast

import (
	"fmt"
	"strings"
)

type directives struct {
	nextCmd   bool
	deleteCmd bool
	quitCmd   bool
	runBlock  *Program
}

type Runtime struct {
	patternSpace string
	holdSpace    string
	appendSpace  string
	lineNo       int
	lines        []string
	program      *Program
	output       string
	directives   directives
}

type RuntimeOptions struct {
	AllowExec      bool
	AutoPrint      bool
	DefaultRuntime *Runtime
	IsBlock        bool
}

func Run(p *Program, text string, options RuntimeOptions) string {
	// TODO: Setup runtime flags
	// Create runtime
	r := &Runtime{
		program: p,
		lines:   strings.Split(text, "\n"),
	}
	if options.DefaultRuntime != nil {
		r = options.DefaultRuntime
		r.directives = directives{}
	}

	retStr := ""

lineLoop:
	for r.lineNo = 0; r.lineNo < len(r.lines); r.lineNo++ {
		fmt.Printf("line[%d]=%s\n", r.lineNo, r.lines[r.lineNo])
		r.output = ""
		r.patternSpace = r.lines[r.lineNo]
		pc := 0
		for pc < len(p.Statements) {
			s := p.Statements[pc]
			fmt.Printf("Running statement: %T\n", s)
			match := s.Address(r)
			if !match {
				pc++
				continue
			}
			s.Run(r)
			if r.directives.nextCmd {
				// Proceed to the next line printing out current pattern space.
				r.directives.nextCmd = false
				if options.AutoPrint {
					r.output += r.patternSpace + "\n"
				}
				retStr += r.output
				continue lineLoop
			} else if r.directives.deleteCmd {
				// Proceed to the next line not printing out current pattern space.
				r.directives.deleteCmd = false
				continue lineLoop
			} else if r.directives.quitCmd {
				// Quit the program with the rest of the pattern space.
				retStr += r.patternSpace
				return retStr
			} else if r.directives.runBlock != nil {
				opt := options
				opt.DefaultRuntime = r
				opt.IsBlock = true
				prevLines, prevLineNo, _ := opt.DefaultRuntime.lines, opt.DefaultRuntime.lineNo, opt.DefaultRuntime.output

				opt.DefaultRuntime.lines = opt.DefaultRuntime.lines[r.lineNo : r.lineNo+1]

				r.output = Run(r.directives.runBlock, "", opt)
				opt.DefaultRuntime.lines = prevLines
				opt.DefaultRuntime.lineNo = prevLineNo
				// opt.DefaultRuntime.output = prevOut
				fmt.Println("exiting block statement")
				fmt.Println("  opt.dr.output=", opt.DefaultRuntime.output)
				fmt.Println("  r.output=", r.output)
			}
			pc++
		}
		if options.AutoPrint && !options.IsBlock {
			r.output += r.patternSpace + "\n"
		}
		fmt.Printf("  output=%s\n", r.output)
		retStr += r.output
		if len(r.appendSpace) > 0 {
			retStr += r.appendSpace
			r.appendSpace = ""
		}
	}
	if len(retStr) > 0 && retStr[len(retStr)-1] == '\n' && !options.IsBlock {
		retStr = retStr[0 : len(retStr)-1]
	}
	fmt.Println("⏎ returning")
	return retStr
}
