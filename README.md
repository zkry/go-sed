# go-sed
A sed implementation in Go. This project aims to be a fully compatable sed clone writen in pure Go with no dependencies. There is also a Go api that you can call to integrate this into other Go programs.

For example, the follwoing will indent the text by five spaces:

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

## Command Line Tool
This also command line tool that you can use in place of the standard sed command. You can install it by running `go get github.com/zkry/go-sed/cmd/gosed`.

## API
The current API for the package is as follows
```
func MustCompile(program string, opt Options) Program
func Compile(program string, opt Options) (Program, error)

func (p Program) Copy(data []byte) Program

func (p Program) Filter(data []byte) []byte
func (p Program) FilterString(data []byte, opt Options) []byte

func (p Program) FilterA(str string) string
func (p Program) FilterStringA(data []byte, opt Options) []byte

```

## Roadmap

- [ ] Thourough test suite for compatability with sed.
- [ ] GNU sed features.
- [ ] Implement more command line flags.
