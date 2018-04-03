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
fmt.Println(cmd.Filter(myTxt))

// Prints the above 'if' statement indented by five spaces to the left
```

This is still a work in progress and is in the very early stages of development.

### Progress:
- [x] Lexer
- [ ] Parser
- [ ] Evaluating Program
- [ ] Command line program
- [ ] Go Library
