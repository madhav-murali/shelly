package pipes

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
