package ast

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Program struct {
	Statements []Statement
	Labels     map[string]*Statement
}

type Addresser interface {
	Address(r *Runtime) bool
}

type Statement interface {
	Addresser
	Run(r *Runtime)
}

type AStmt struct {
	Addresser
	AppendLine string
}

func (s *AStmt) Run(r *Runtime) {
	r.appendSpace += s.AppendLine + "\n"
}

type BStmt struct {
	Addresser
	BranchIdent string
}

func (s *BStmt) Run(r *Runtime) {
}

type CStmt struct {
	Addresser
	ChangeLine string
}

func (s *CStmt) Run(r *Runtime) {
}

// SFlags represents the various options that can be passed to the s command.
// The zero value means the flag is not set.
type SFlags struct {
	NFlag int    // N - Make the substitution only for the Nth occurence of regexp
	GFlag bool   // g - Make the substitution for all non-overlapping matches
	PFlag bool   // p - Write the pattern space to stdout
	WFile string // w file  - append pattern space to file if a replacement made.
}

type SStmt struct {
	Addresser
	FindAddr    string
	ReplaceAddr string
	Flags       SFlags
}

func (s *SStmt) Run(r *Runtime) {
	var rgxp *regexp.Regexp
	var err error
	if s.Flags.GFlag {
		// Replace for all occurences.
		rgxp, err = regexp.Compile(s.FindAddr)
		if err != nil {
			return
		}
		r.patternSpace = rgxp.ReplaceAllString(r.patternSpace, s.ReplaceAddr)
	} else {
		matchIdx := 0
		if s.Flags.NFlag != 0 {
			// Replace for nth occurence.
			matchIdx = s.Flags.NFlag - 1
		}

		// } else {
		// Replace for first occurence.
		rgxp, err = regexp.Compile(s.FindAddr)
		if err != nil {
			return
		}
		loc := rgxp.FindAllStringIndex(r.patternSpace, -1)
		if loc == nil {
			return
		}
		r.patternSpace = r.patternSpace[0:loc[matchIdx][0]] + rgxp.ReplaceAllString(r.patternSpace[loc[matchIdx][0]:loc[matchIdx][1]], s.ReplaceAddr) + r.patternSpace[loc[matchIdx][1]:len(r.patternSpace)]
	}
	if s.Flags.PFlag {
		r.output += r.patternSpace + "\n"
	}
}

type DStmt struct {
	Addresser
}

func (s *DStmt) Run(r *Runtime) {
	r.directives.deleteCmd = true
}

type D2Stmt struct {
	Addresser
}

func (s *D2Stmt) Run(r *Runtime) {
}

type EStmt struct {
	Addresser
	Command string
}

func (s *EStmt) Run(r *Runtime) {
}

type GStmt struct {
	Addresser
}

func (s *GStmt) Run(r *Runtime) {
	r.patternSpace = r.holdSpace
}

type G2Stmt struct {
	Addresser
}

func (s *G2Stmt) Run(r *Runtime) {
}

type HStmt struct {
	Addresser
}

func (s *HStmt) Run(r *Runtime) {
	r.holdSpace = r.patternSpace
}

type H2Stmt struct {
	Addresser
}

func (s *H2Stmt) Run(r *Runtime) {
}

type IStmt struct {
	Addresser
	InsertLine string
}

func (s *IStmt) Run(r *Runtime) {
	r.output += s.InsertLine + "\n"
}

type LStmt struct {
	Addresser
}

func (s *LStmt) Run(r *Runtime) {
}

type NStmt struct {
	Addresser
}

func (s *NStmt) Run(r *Runtime) {
	r.directives.nextCmd = true
}

type N2Stmt struct {
	Addresser
}

func (s *N2Stmt) Run(r *Runtime) {
}

type PStmt struct {
	Addresser
}

func (s *PStmt) Run(r *Runtime) {
	r.output += r.patternSpace + "\n"
}

type P2Stmt struct {
	Addresser
}

func (s *P2Stmt) Run(r *Runtime) {
}

type QStmt struct {
	Addresser
}

func (s *QStmt) Run(r *Runtime) {
	r.directives.quitCmd = true
}

type RStmt struct {
	Addresser
	FileName string
}

func (s *RStmt) Run(r *Runtime) {
}

type R2Stmt struct {
	Addresser
	FileName string
}

func (s *R2Stmt) Run(r *Runtime) {
}

type TStmt struct {
	Addresser
	FileName string
}

func (s *TStmt) Run(r *Runtime) {
}

type T2Stmt struct {
	Addresser
	FileName string
}

func (s *T2Stmt) Run(r *Runtime) {
}

type WStmt struct {
	Addresser
	FileName string
}

func (s *WStmt) Run(r *Runtime) {
}

type W2Stmt struct {
	Addresser
	FileName string
}

func (s *W2Stmt) Run(r *Runtime) {
}

type XStmt struct {
	Addresser
}

func (s *XStmt) Run(r *Runtime) {
	r.patternSpace, r.holdSpace = r.holdSpace, r.patternSpace
}

type YStmt struct {
	Addresser
	charMap map[rune]rune
}

func (s *YStmt) Run(r *Runtime) {
	var newPS string
	for _, r := range r.patternSpace {
		if nr, ok := s.charMap[r]; ok {
			newPS += string(nr)
		} else {
			newPS += string(r)
		}
	}
	r.patternSpace = newPS
}

func NewYStmt(find, replace string, addr Addresser) (*YStmt, error) {
	fRunes := []rune{}
	rRunes := []rune{}
	for _, r := range find {
		fRunes = append(fRunes, r)
	}
	for _, r := range replace {
		rRunes = append(rRunes, r)
	}

	if len(fRunes) != len(rRunes) {
		return nil, errors.New("transform strings are not the same length")
	}

	cm := make(map[rune]rune)
	for i := range fRunes {
		cm[fRunes[i]] = rRunes[i]
	}

	return &YStmt{charMap: cm}, nil
}

type ZStmt struct {
	Addresser
}

func (s *ZStmt) Run(r *Runtime) {
}

type EquStmt struct {
	Addresser
}

func (s *EquStmt) Run(r *Runtime) {
	r.output += strconv.Itoa(r.lineNo+1) + "\n"
}

type BlockStmt struct {
	Code *Program
	Addresser
}

func (s *BlockStmt) Run(r *Runtime) {
	fmt.Println("↓ Running BlockStmt")
	r.directives.runBlock = s.Code
}

// func (s *Statement) Run(r *Runtime) {
// 	switch s.Cmd {
// 	default:
// 	}
// 	return l
// }

type RegexpAddr struct {
	Regexp *regexp.Regexp
}

func (a *RegexpAddr) Address(r *Runtime) bool {
	return len(a.Regexp.FindIndex([]byte(r.patternSpace))) > 0
}

type LineNoAddr struct {
	LineNo int
}

func (a *LineNoAddr) Address(r *Runtime) bool {
	if r.lineNo+1 == a.LineNo {
		return true
	}
	return false
}

type EOFAddr struct{}

func (a *EOFAddr) Address(r *Runtime) bool {
	return r.lineNo == len(r.lines)-1
}

type RangeAddress struct {
	Addr1 Addresser
	Addr2 Addresser
	on    bool
}

func (a *RangeAddress) Address(r *Runtime) bool {
	if a.on {
		if a.Addr2.Address(r) {
			a.on = false
		}
		return true
	}
	if a.Addr1.Address(r) {
		a.on = true
		return true
	}
	return false
}

type BlankAddress struct{}

func (a *BlankAddress) Address(r *Runtime) bool {
	return true
}
