package ast

import "regexp"

type Program struct {
	Statements []Statement
}

type Statement struct {
	Addr Address
	Cmd  string

	// Arg1 and Arg2 are the two arguments after s and y commands
	Arg1 string
	Arg2 string

	Flags []string
}

func (s *Statement) Run(isEOF bool, lineNo int, l string) string {
	switch s.Cmd {
	default:
	}
	return l
}

type Address interface {
	Address(isEOF bool, lineNo int, l string) bool
}

type RegexpAddr struct {
	Regexp *regexp.Regexp
}

func (a *RegexpAddr) Address(isEOF bool, lineNo int, l string) bool {
	return len(a.Regexp.FindIndex([]byte(l))) > 0
}

type LineNoAddr struct {
	LineNo int
}

func (a *LineNoAddr) Address(isEOF bool, lineNo int, l string) bool {
	if lineNo == a.LineNo {
		return true
	}
	return false
}

type EOFAddr struct{}

func (a *EOFAddr) Address(isEOF bool, lineNo int, l string) bool {
	return isEOF
}

type RangeAddress struct {
	Addr1 Address
	Addr2 Address
	on    bool
}

func (a *RangeAddress) Address(isEOF bool, lineNo int, l string) bool {
	if a.on {
		if a.Addr2.Address(isEOF, lineNo, l) {
			a.on = false
		}
		return true
	}
	if a.Addr1.Address(isEOF, lineNo, l) {
		a.on = true
		return true
	}
	return false
}

type BlankAddress struct{}

func (a *BlankAddress) Address(isEOF bool, lineNo int, l string) bool {
	return true
}
