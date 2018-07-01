package ast

import (
	"strings"
)

type directives struct {
	nextCmd   bool
	deleteCmd bool
	quitCmd   bool
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
	AllowExec bool
	AutoPrint bool
}

func Run(p *Program, text string, options RuntimeOptions) string {
	// TODO: Setup runtime flags
	// Create runtime
	r := &Runtime{
		program: p,
		lines:   strings.Split(text, "\n"),
	}
	retStr := ""

lineLoop:
	for r.lineNo = 0; r.lineNo < len(r.lines); r.lineNo++ {
		r.output = ""
		r.patternSpace = r.lines[r.lineNo]
		pc := 0
		for pc < len(p.Statements) {
			s := p.Statements[pc]
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
			}
			pc++
		}
		if options.AutoPrint {
			r.output += r.patternSpace + "\n"
		}
		retStr += r.output
		if len(r.appendSpace) > 0 {
			retStr += r.appendSpace
			r.appendSpace = ""
		}
	}
	if len(retStr) > 0 && retStr[len(retStr)-1] == '\n' {
		retStr = retStr[0 : len(retStr)-1]
	}
	return retStr
}
