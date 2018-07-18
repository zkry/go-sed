package gosed

import "github.com/zkry/go-sed/ast"

type Options struct {
	SupressOutput bool // Prevents program from automatically outputing line.
	AppendFile    bool // Makes the w command append to file.
	ExtendRegexp  bool // Use extended version of regexp
}

type Program struct {
	p ast.Program
}

// MustCompile takes a sed script and compiles it into a program.
// Panics if errors are found in script.
func MustCompile(program string, opt Options) *Program {
	panic("not implemented")
	return nil
}

// Compile compiles a sed script and returns a program upon successfull
// compilation. If unsuccessfull errors are returned.
func Compile(program string, opt Options) (*Program, error) {
	panic("not implemented")
	return nil, nil
}

func (p *Program) Filter(data []byte) []byte {

}

func (p *Program) FilterString(data string) string {

}
