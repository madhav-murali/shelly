package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	customDS "github.com/codecrafters-io/shell-starter-go/internal/custom"
)

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, _ := cmd.CombinedOutput()
	fmt.Print(string(out))
}

func main() {
	fmt.Print("$ ")
	line, err := (bufio.NewReader(os.Stdin).ReadString('\n'))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
	line = line[:len(line)-1]
	if len(line) == 0 {
		main()
	}

	validCmds := customDS.NewCmdSet()
	validCmds.Add("EXIT")
	validCmds.Add("ECHO")
	validCmds.Add("TYPE")

	args := strings.Split(line, " ")
	cmd := args[0]
	switch strings.ToUpper(cmd) {
	case "EXIT":
		os.Exit(0)
	case "ECHO":
		for i := 1; i < len(args); i++ {
			fmt.Print(args[i])
			if i != len(args)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	case "TYPE":
		if validCmds.Find(strings.ToUpper(args[1])) {
			fmt.Printf("%s is a shell builtin\n", args[1])
		} else {
			path, err := exec.LookPath(args[1])
			if err != nil {
				fmt.Printf("%s: not found\n", args[1])
				break
			}
			fmt.Printf("%s is %s\n", args[1], path)
		}
	default:
		_, err := exec.LookPath(args[0])
		if err != nil {
			fmt.Printf("%s: not found\n", args[1])
			break
		}
		runCmd(args[0], args[1:]...)
	}
	main()
}
