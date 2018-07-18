package main

import (
	"flag"
	"fmt"
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
	fmt.Printf("Calling Set(%s)\n", v)
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
	fmt.Printf("Calling Set(%s)\n", v)
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

	fmt.Println("Running gosed:")
	// fmt.Printf("  commandFiles: %v\n", config.commandFiles.String())
	fmt.Printf("  eCommands:     %v\n", config.eCommands)
	fmt.Printf("  fileCommands:     %v\n", config.fileCommands)
	fmt.Printf("  silenced:     %v\n", config.silenceLine)
}
