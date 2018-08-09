package gosed

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

// TestInfo tests to see if the ending positions returned from Info
// are as expected.
func TestInfo(t *testing.T) {
	cases := []struct {
		program   string
		positions []int
	}{
		{"s/one/two/g", []int{1, 2, 5, 6, 9, 10, 11}},
	}

infoTests:
	for _, c := range cases {
		tokens := Info(c.program)
		for i := range tokens {
			if len(c.positions) == 0 {
				break
			}
			if tokens[i].End != c.positions[0] {
				t.Errorf("Info produced wrong positions:\n  Got: %v\n  Expected: %v\n", tokens[i].End, c.positions[i])
				continue infoTests
			}
			c.positions = c.positions[1:]
		}
	}
}

var update = flag.Bool("update", false, "update golden files")

// The -u flag may be used to generate the golden files from the gsed command. It should
// also be noted that the first line on the sed script determins with what options and
// against which files the sed program should be compared. Some of the test programs are
// genaral file editing which can work with any input while some are specially designed for
// certain files. The first items on the comment determins whith which flags the commadn
// should be run (beginnign with a '-' character) while the last parameter determines the
// specific file. If no files are mentiond, it will be run against all of the test files.
//
// There are three directories dealing with this test case:
//
//   programs/
//   inputs/
//   expected/
//
// Program contains the test programs. Tests by default are expected to compile sucessfully
// and any compile error will fail the test. If the programs file name contains a _fail.sed
// then it should be expected to fail.
//
// Unless otherwise stated in the programs first comment, all programs will be tested agains
// all inputs. The expected will be generated into the expected folder with the name of the
// file being the program name and the input name with a '_' dividor. Thus, reverse-lines.sed
// ran with input input.txt will produce the file reverse-liens_input.txt.
func TestSed(t *testing.T) {
	const programDir = "./testdata/programs"
	const inputDir = "./testdata/inputs"
	const expectDir = "./testdata/expected"

	prgFiles, err := ioutil.ReadDir(programDir)
	if err != nil {
		panic(err)
	}

	for _, prgF := range prgFiles {
		path := path.Join(programDir, prgF.Name())
		prgData, err := ioutil.ReadFile(path)
		if err != nil {
			panic("could not open file: " + path)
		}

		// var inputFile string
		opt := Options{}

		// Get the options from the first line.
		firstLine := bytes.Split(prgData, []byte("\n"))[0]
		for _, setting := range bytes.Split(firstLine, []byte(" ")) {
			switch s := string(setting); s {
			case "#":
			case "-n":
				opt.SupressOutput = true
			default:
				if strings.HasSuffix(s, ".txt") {
					// inputFile = s
				}
			}
		}

		// Attemp to compile the program, returning any errors if compilation failed.
		_, errs := Compile(string(prgData), opt)
		if len(errs) != 0 {
			errDescBuff := bytes.Buffer{}
			for _, errStr := range errs {
				errDescBuff.WriteString("\n" + errStr)
			}
			t.Errorf("Program %s did not compile. %s", path, errDescBuff.String())
		}
	}
}
