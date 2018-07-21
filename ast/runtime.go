package ast

import (
	"strings"
)

type directives struct {
	nextCmd       bool
	deleteCmd     bool
	restartScript bool // Used for the 'D' command
	quitCmd       bool
	quitNoPattern bool
	runBlock      *Program
	jumpTo        string
}

type runtime struct {
	patternSpace string
	holdSpace    string
	appendSpace  string
	lineNo       int
	lines        []string
	program      *Program
	output       string
	directives   directives
	subMade      bool
}

type RuntimeOptions struct {
	AllowExec      bool
	AutoPrint      bool
	AppendFile     bool
	DefaultRuntime *runtime
	IsBlock        bool
	LineNoStart    int
}

func (p *Program) Run(text string, options RuntimeOptions) string {
	// TODO: Setup runtime flags
	// Create runtime
	r := &runtime{
		program: p,
		lines:   strings.Split(text, "\n"),
	}
	if options.DefaultRuntime != nil {
		r = options.DefaultRuntime
		r.directives = directives{}
	}

	retStr := ""

lineLoop:
	for r.lineNo = options.LineNoStart; r.lineNo < len(r.lines); r.lineNo++ {
		// fmt.Printf("\n\nline[%d]=%s\n", r.lineNo, r.lines[r.lineNo])
		r.output = ""
		r.patternSpace = r.lines[r.lineNo]
		r.subMade = false
		pc := 0
		for pc < len(p.Statements) {
			s := p.Statements[pc]
			// fmt.Printf("Running statement[%d]: %T\n", pc, s)
			match := s.Address(r)
			if !match {
				pc++
				continue
			}
			s.Run(r)
			if r.directives.nextCmd {
				// fmt.Println("[D] nextCmd")
				// Proceed to the next line printing out current pattern space.
				r.directives.nextCmd = false
				if options.AutoPrint {
					r.output += r.patternSpace + "\n"
				}
				retStr += r.output
				continue lineLoop
			} else if r.directives.deleteCmd {
				// fmt.Println("[D] deleteCmd")
				// Proceed to the next line not printing out current pattern space.
				r.directives.deleteCmd = false
				retStr += r.output
				continue lineLoop
			} else if r.directives.quitCmd {
				// fmt.Println("[D] quitCmd")
				r.directives.quitCmd = false
				// Quit the program with the rest of the pattern space.
				retStr += r.output
				retStr += r.patternSpace
				return retStr
			} else if r.directives.quitNoPattern {
				// fmt.Println("[D] quit No Pattern")
				r.directives.quitNoPattern = false
				// Quit the program with the rest of the output.
				retStr += r.output
				return retStr
			} else if r.directives.runBlock != nil {
				// fmt.Println("[D] runBlock")
				opt := options
				opt.DefaultRuntime = r
				opt.IsBlock = true
				opt.LineNoStart = r.lineNo

				r.output = r.directives.runBlock.Run("", opt)

				// fmt.Println("exiting block statement")
				// fmt.Println("  opt.dr.output=", opt.DefaultRuntime.output)
				// fmt.Println("  r.output=", r.output)
			} else if r.directives.jumpTo != "" {
				// fmt.Println("[D] jumpTo")
				label := r.directives.jumpTo
				r.directives.jumpTo = ""
				// fmt.Printf("Jumping to %s@%d\n", label, p.Labels[label])
				pc = p.Labels[label]
				continue
			} else if r.directives.restartScript {
				// fmt.Println("[D] restartScript")
				r.directives.restartScript = false
				pc = 0
				continue
			}
			pc++
		}
		if options.AutoPrint && !options.IsBlock {
			r.output += r.patternSpace + "\n"
		}
		// fmt.Printf("  output=%s\n", r.output)
		retStr += r.output
		if len(r.appendSpace) > 0 {
			retStr += r.appendSpace
			r.appendSpace = ""
		}
		if options.IsBlock {
			return retStr
		}
	}
	if len(retStr) > 0 && retStr[len(retStr)-1] == '\n' && !options.IsBlock {
		retStr = retStr[0 : len(retStr)-1]
	}
	// fmt.Println("⏎ returning")
	return retStr
}
