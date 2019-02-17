package ast

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/zkry/go-sed/lexer"
)

type Program struct {
	Statements []statement
	Labels     map[string]int
	Tokens     []lexer.Item
}

type addresser interface {
	Address(r *runtime) bool
}

type statement interface {
	addresser
	Run(r *runtime)
}

type aStmt struct {
	addresser
	AppendLine string
}

func (s *aStmt) Run(r *runtime) {
	r.appendSpace += s.AppendLine + "\n"
}

type bStmt struct {
	addresser
	BranchIdent string
}

func (s *bStmt) Run(r *runtime) {
	r.directives.jumpTo = s.BranchIdent
}

type cStmt struct {
	addresser
	ChangeLine string
	prevResult bool
}

func (s *cStmt) Run(r *runtime) {
	// TODO: To be implemented
	switch s.addresser.(type) {
	case *rangeAddress:
		match := s.Address(r)
		if !match {
			if s.prevResult {
				// output the ChangeLine
				r.directives.deleteCmd = true
				r.output += s.ChangeLine + "\n"
			}
			s.prevResult = false
			return
		}
		r.directives.deleteCmd = true
		s.prevResult = true
	default:
		match := s.Address(r)
		if !match {
			return
		}
		r.directives.deleteCmd = true
		r.output += s.ChangeLine + "\n"
	}
}

// SFlags represents the various options that can be passed to the s command.
// The zero value means the flag is not set.
type sFlags struct {
	NFlag int    // N - Make the substitution only for the Nth occurence of regexp
	GFlag bool   // g - Make the substitution for all non-overlapping matches
	PFlag bool   // p - Write the pattern space to stdout
	WFile string // w file  - append pattern space to file if a replacement made.
}

type sStmt struct {
	addresser
	FindAddr    string
	ReplaceAddr string
	Flags       sFlags
}

func (s *sStmt) Run(r *runtime) {
	var rgxp *regexp.Regexp
	var err error
	if s.Flags.GFlag {
		// Replace for all occurences.
		rgxp, err = regexp.Compile(s.FindAddr)
		if err != nil {
			return
		}
		r.subMade = true
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
		r.subMade = true
		r.patternSpace = r.patternSpace[0:loc[matchIdx][0]] + rgxp.ReplaceAllString(r.patternSpace[loc[matchIdx][0]:loc[matchIdx][1]], s.ReplaceAddr) + r.patternSpace[loc[matchIdx][1]:len(r.patternSpace)]
	}
	if s.Flags.PFlag {
		r.output += r.patternSpace + "\n"
	}
}

type dStmt struct {
	addresser
}

func (s *dStmt) Run(r *runtime) {
	r.directives.deleteCmd = true
}

type d2Stmt struct {
	addresser
}

func (s *d2Stmt) Run(r *runtime) {
	idx := strings.IndexRune(r.patternSpace, '\n')
	if idx == -1 {
		r.directives.deleteCmd = true
	}
	r.patternSpace = r.patternSpace[idx+1:]
	r.directives.restartScript = true
}

type eStmt struct {
	addresser
	Command string
}

func (s *eStmt) Run(r *runtime) {
	// TODO: To be implemented
}

type gStmt struct {
	addresser
}

func (s *gStmt) Run(r *runtime) {
	r.patternSpace = r.holdSpace
}

type g2Stmt struct {
	addresser
}

func (s *g2Stmt) Run(r *runtime) {
	r.patternSpace += "\n" + r.holdSpace
}

type hStmt struct {
	addresser
}

func (s *hStmt) Run(r *runtime) {
	r.holdSpace = r.patternSpace
}

type h2Stmt struct {
	addresser
}

func (s *h2Stmt) Run(r *runtime) {
	r.holdSpace += "\n" + r.patternSpace
}

type iStmt struct {
	addresser
	InsertLine string
}

func (s *iStmt) Run(r *runtime) {
	r.output += s.InsertLine + "\n"
}

type lStmt struct {
	addresser
}

func (s *lStmt) Run(r *runtime) {
}

type nStmt struct {
	addresser
}

func (s *nStmt) Run(r *runtime) {
	r.directives.nextCmd = true
}

type n2Stmt struct {
	addresser
}

func (s *n2Stmt) Run(r *runtime) {
	r.lineNo++
	if r.lineNo >= len(r.lines) {
		r.directives.quitNoPattern = true
		return
	}
	r.patternSpace += "\n" + r.lines[r.lineNo]
}

type pStmt struct {
	addresser
}

func (s *pStmt) Run(r *runtime) {
	r.output += r.patternSpace + "\n"
}

type p2Stmt struct {
	addresser
}

func (s *p2Stmt) Run(r *runtime) {
	idx := strings.IndexRune(r.patternSpace, '\n')
	if idx == -1 {
		r.output += r.patternSpace
		return
	}
	r.output += r.patternSpace[:idx]
}

type qStmt struct {
	addresser
}

func (s *qStmt) Run(r *runtime) {
	r.directives.quitCmd = true
}

type rStmt struct {
	addresser
	FileName string
}

func (s *rStmt) Run(r *runtime) {
}

type r2Stmt struct {
	addresser
	FileName string
}

func (s *r2Stmt) Run(r *runtime) {
	// To be implemented
}

type tStmt struct {
	addresser
	BranchIdent string
}

func (s *tStmt) Run(r *runtime) {
	if r.subMade {
		r.subMade = false
		r.directives.jumpTo = s.BranchIdent
	}
}

type t2Stmt struct {
	addresser
	FileName string
}

func (s *t2Stmt) Run(r *runtime) {
	// To be implemented
}

type wStmt struct {
	addresser
	FileName string
}

func (s *wStmt) Run(r *runtime) {
	// To be implemented
}

type w2Stmt struct {
	addresser
	FileName string
}

func (s *w2Stmt) Run(r *runtime) {
	// To be implemented
}

type xStmt struct {
	addresser
}

func (s *xStmt) Run(r *runtime) {
	r.patternSpace, r.holdSpace = r.holdSpace, r.patternSpace
}

type yStmt struct {
	addresser
	charMap map[rune]rune
}

func (s *yStmt) Run(r *runtime) {
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

func newYStmt(find, replace string, addr addresser) (*yStmt, error) {
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

	return &yStmt{charMap: cm}, nil
}

type zStmt struct {
	addresser
}

func (s *zStmt) Run(r *runtime) {
	// To be implemented
}

type equStmt struct {
	addresser
}

func (s *equStmt) Run(r *runtime) {
	r.output += strconv.Itoa(r.lineNo+1) + "\n"
}

type blockStmt struct {
	Code *Program
	addresser
}

func (s *blockStmt) Run(r *runtime) {
	r.directives.runBlock = s.Code
}

type regexpAddr struct {
	Regexp *regexp.Regexp
}

func (a *regexpAddr) Address(r *runtime) bool {
	return len(a.Regexp.FindIndex([]byte(r.patternSpace))) > 0
}

type lineNoAddr struct {
	LineNo int
}

func (a *lineNoAddr) Address(r *runtime) bool {
	return r.lineNo+1 == a.LineNo
}

type eofAddr struct{}

func (a *eofAddr) Address(r *runtime) bool {
	return r.lineNo == len(r.lines)-1
}

type notAddr struct {
	Addr addresser
}

func (a *notAddr) Address(r *runtime) bool {
	return !a.Address(r)
}

type rangeAddress struct {
	Addr1 addresser
	Addr2 addresser
	on    bool
}

func (a *rangeAddress) Address(r *runtime) bool {
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

type blankAddress struct{}

func (a *blankAddress) Address(r *runtime) bool {
	return true
}
