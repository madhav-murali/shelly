package builtins

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/custom"
)

func ExecuteBuiltin(args []string, in io.Reader, out io.WriteCloser, validCmds *custom.CmdSet) {
	switch strings.ToUpper(args[0]) {
	case "TYPE":
		Type(out, args, validCmds)
	case "PWD":
		PWD(out)
	case "CD":
		CD(out, args)
	case "ECHO":
		ECHO(out, args)
	}
}

func Type(t io.Writer, args []string, validCmds *custom.CmdSet) {
	if validCmds.Find(strings.ToUpper(args[1])) {
		fmt.Fprintf(t, "%s is a shell builtin\n", args[1])
	} else {
		path, err := exec.LookPath(args[1])
		if err != nil {
			fmt.Fprintf(t, "%s: not found\n", args[1])
			return
		}
		fmt.Fprintf(t, "%s is %s\n", args[1], path)
	}
}

func ECHO(w io.Writer, args []string) error {
	fmt.Fprintf(w, "%s\n", strings.Join(args[1:], " "))
	return nil
}

func PWD(t io.Writer) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(t, "error: ", err)
		return
	}
	fmt.Fprintln(t, dir)
}

func CD(t io.Writer, args []string) {
	cd := args[1]
	var err error
	if cd == "~" {
		cd, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(t, err)
			return
		}
	}
	err = os.Chdir(cd)
	if err != nil {
		fmt.Fprintf(t, "cd: %s: No such file or directory\n", args[1])
		//break
	}
}
