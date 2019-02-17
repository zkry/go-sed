package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	gosed "github.com/zkry/go-sed"
)

var order int

type Command struct {
	order int
	cmd   string
}
type FileCommands []Command
type ECommands []Command

func (a *ECommands) String() string {
	if a == nil {
		return "[nil]"
	}
	return ""
}

func (a *ECommands) Set(v string) error {
	newCmd := Command{
		order: order,
		cmd:   v,
	}
	*a = append(*a, newCmd)
	order++
	return nil
}

func (a *FileCommands) String() string {
	if a == nil {
		return "[nil]"
	}
	return ""
}

func (a *FileCommands) Set(v string) error {
	newCmd := Command{
		order: order,
		cmd:   v,
	}
	*a = append(*a, newCmd)
	order++
	return nil
}

// Config specifies all of the pass-in parameters that the
// command can take.
type Config struct {
	fileCommands     FileCommands // Translates to -e flag
	eCommands        ECommands    // Translates to -e flag
	editInplace      bool         // Translates to -i flag
	inplaceExtension string       // Prameter for -i flag
	extendedRegexp   bool         // Translates to -E flag
	appendFile       bool         // Translates to -a flag
	bufferedOutput   bool         // Translates to -l flag
	silenceLine      bool         // Translates to -n flag
	commandCt        int
	interactive      bool
}

func combineInputs(files []string) []byte {
	var buff bytes.Buffer
	for i, f := range files {
		d, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Printf("gosed: %s: %v\n", f, err)
		}
		if i > 0 {
			buff.WriteRune('\n')
		}
		buff.Write(d)
	}
	return buff.Bytes()
}

func programFromConfig(conf Config) (*gosed.Program, error) {
	var programBuff bytes.Buffer
	// Iterate through all of the commands, processing the two slices of commands,
	// (conf.fileCommands and conf.eCommands) in the order that they arrived.
	for {
		if len(conf.fileCommands) == 0 && len(conf.eCommands) == 0 {
			break
		} else if len(conf.fileCommands) == 0 || conf.eCommands[0].order < conf.fileCommands[0].order {
			cmd := conf.eCommands[0].cmd
			conf.eCommands = conf.eCommands[1:]
			programBuff.WriteString("\n" + cmd)
		} else {
			fname := conf.fileCommands[0].cmd
			conf.fileCommands = conf.fileCommands[1:]
			fdata, err := ioutil.ReadFile(fname)
			if err != nil {
				return nil, errors.New("could not read file " + fname)
			}
			programBuff.WriteRune('\n')
			programBuff.Write(fdata)
		}
	}

	program, err := gosed.Compile(programBuff.String(), gosed.Options{})
	if err != nil {
		return nil, errors.New("syntax error: " + err.Error())
	}
	return program, nil
}

func runFromStdin(program *gosed.Program) {
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		out := program.FilterString(line)
		fmt.Print(out)
	}
}

func main() {
	var helpFlag bool
	var config Config
	// flag.Var(&config.commandFiles, "f", "")
	flag.Var(&config.fileCommands, "f", "file with sed commands to run")
	flag.Var(&config.eCommands, "e", "string of command to execute")
	flag.BoolVar(&config.silenceLine, "n", false, "silence the auto-print-line functionality")
	flag.BoolVar(&helpFlag, "h", false, "usage guide")

	// TODO: Implement support for following flags
	//flag.BoolVar(&config.bufferedOutput, "l", false, "")
	//flag.BoolVar(&config.appendFile, "a", false, "")
	//flag.BoolVar(&config.extendedRegexp, "E", false, "")
	//flag.BoolVar(&config.interactive, "i", false, "")
	flag.Parse()
	config.commandCt = order

	if helpFlag || (config.commandCt == 0 && flag.NArg() == 0) {
		displayHelp()
		return
	}
	//if config.interactive {
	//runInteractive()
	//}

	if config.commandCt > 0 {
		program, err := programFromConfig(config)
		if err != nil {
			fmt.Println(err)
			return
		}
		if flag.NArg() > 0 {
			// Read files and send them through commands.
			programInput := combineInputs(flag.Args())
			out := program.Filter(programInput)
			fmt.Print(string(out))
		} else {
			// Read Stdout through commands
			runFromStdin(program)
		}
		return
	}

	if flag.NArg() > 0 {
		// Use arg[0] as command and arg[1:] as input files. If only one arg,
		// read from stdout
		fname := flag.Arg(0)
		program, err := gosed.Compile(fname, gosed.Options{})
		if err != nil {
			fmt.Printf("gosed: syntax error in %s\nerror:%v\n", fname, err.Error())
			return
		}
		if flag.NArg() == 1 {
			runFromStdin(program)
			return
		}
		programInput := combineInputs(flag.Args()[1:])
		out := program.Filter(programInput)
		fmt.Print(string(out))
	}
}

func displayHelp() {
	fmt.Println(`gosed is a Go implementation of the Sed stream editor. 
It seeks to behave as the standard sed command. There is currently
no support for the GNU sed features. The following in an example of
running two commands processing the file 'input.txt'. 

      gosed -e 's/one/ONE/g' -e '/two/d' input.txt

You can also run a file with sed commands as follows:

      gosed -f processor.sed mytext.txt

Flags:`)
	flag.PrintDefaults()
}
