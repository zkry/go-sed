package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"

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
	for {
		if len(conf.fileCommands) == 0 && len(conf.eCommands) == 0 {
			break
		} else if len(conf.fileCommands) == 0 || conf.eCommands[0].order < conf.eCommands[0].order {
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

func displayHelp() {

}

func main() {
	var config Config
	// flag.Var(&config.commandFiles, "f", "")
	flag.Var(&config.fileCommands, "f", "")
	flag.Var(&config.eCommands, "e", "")
	flag.BoolVar(&config.silenceLine, "n", false, "")
	flag.BoolVar(&config.bufferedOutput, "l", false, "")
	flag.BoolVar(&config.appendFile, "a", false, "")
	flag.BoolVar(&config.extendedRegexp, "E", false, "")
	flag.Parse()
	config.commandCt = order

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
			// TODO: Line by line stdin processing
		}
	} else {
		if flag.NArg() > 0 {
			// Use arg[0] as command and arg[1:] as input files. If only one arg,
			// read from stdout
			fname := flag.Arg(0)
			program, err := gosed.Compile(fname, gosed.Options{})
			if err != nil {
				fmt.Printf("gosed: syntax error in %s\n", fname)
				return
			}
			if flag.NArg() == 1 {
				// TODO: Line by line stdin processing
				return
			}
			programInput := combineInputs(flag.Args()[1:])
			out := program.Filter(programInput)
			fmt.Print(string(out))
		} else {
			displayHelp()
			return
		}
	}
}
