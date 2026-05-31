package pipes

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/internal/builtins"
	"github.com/codecrafters-io/shell-starter-go/internal/custom"
	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

// This function returns the commands in the pipe to be executed in order
func ParseLineCmds(line string) [][]string {
	cmdSlices := []string{}
	cmdSlices = strings.Split(line, "|")
	ret := [][]string{}
	for _, cmd := range cmdSlices {
		ret = append(ret, parser.ParseLine(cmd))
	}
	return ret
}

func ExecutePipeline(parsedCmds [][]string) error {
	var cmds []*exec.Cmd

	for _, args := range parsedCmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmds = append(cmds, cmd)
	}

	for idx := 0; idx < len(cmds)-1; idx++ {
		stdOutPipe, err := cmds[idx].StdoutPipe()
		if err != nil {
			return fmt.Errorf("encountered error: %v", err)
		}

		cmds[idx+1].Stdin = stdOutPipe
	}

	cmds[0].Stdin = os.Stdin
	cmds[len(cmds)-1].Stdout = os.Stdout

	for _, cmd := range cmds {
		cmd.Stderr = os.Stderr
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command: %v", err)
		}
	}

	for _, cmd := range cmds {
		cmd.Wait()
	}
	return nil
}

func ExecPipeline(parserCmds [][]string, validCmds *custom.CmdSet) error {
	var wg sync.WaitGroup

	var prevReader io.Reader = os.Stdin
	for i, args := range parserCmds {
		isLast := i == len(parserCmds)-1

		var nextWriter io.WriteCloser
		var currentReader = prevReader

		if isLast {
			nextWriter = os.Stdout
		} else {
			pr, pw := io.Pipe()
			nextWriter = pw
			prevReader = pr
		}

		wg.Add(1)

		go func(args []string, in io.Reader, out io.WriteCloser, isLast bool) error {
			defer wg.Done()

			if !isLast {
				defer out.Close()
			}

			if closer, ok := in.(io.Closer); ok && in != os.Stdin {
				defer closer.Close()
			}

			if validCmds.Find(strings.ToUpper(args[0])) {
				builtins.ExecuteBuiltin(args, in, out, validCmds)
				return nil
			}
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdin = in
			cmd.Stdout = out
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				//handle the error; outside too
				return err
			}
			return nil
		}(args, currentReader, nextWriter, isLast)
	}

	wg.Wait()
	return nil
}
