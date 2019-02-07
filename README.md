### NOTE: This is a work in progress and is currently not in a functional state as the lexer is going through a complete rewrite.

# go-sed
A sed implementation in Go.

This is a work in progress project that aims to be a fully compatable sed clone writen in pure Go. I will also have a library component that you can use to filter text. It should end up looking as follows:

```go
myTxt := `
if (!the_program) {
  if (optind < argc) {
	char *arg = argv[optind++];
    the_program = compile_string(the_program, arg, strlen(arg));
  } else {
      usage(4);
  }
}`

cmd := sed.MustCompileProgram('s/^/     /')

fmt.Println(cmd.FilterString(myTxt))
```

## Usage
Here is a possible usage API... still trying to work this out.
```
func MustCompileProgram(program string, opt Options) Program
func CompileProgram(program string, opt Options) (Program, error)

func (p Program) Copy(data []byte) Program

func (p Program) Filter(data []byte) []byte
func (p Program) FilterWithOptions(data []byte, opt Options) []byte

func (p Program) FilterString(str string) string
func (p Program) FilterStringWithOptions(data []byte, opt Options) []byte

func (p Program) Stream(chan []byte) chan []byte
func (p Program) StreamWithOptions(chan []byte, opt Options) chan []byte
```

This is still a work in progress and is in the very early stages of development.

### Progress:
- [x] Lexer
- [x] Parser
- [x] Evaluating Program
- [ ] Command line program
- [ ] Go Library
