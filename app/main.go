package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	set "github.com/codecrafters-io/shell-starter-go/internal/custom"
)

func main() {
	fmt.Print("$ ")
	line, err := (bufio.NewReader(os.Stdin).ReadString('\n'))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
	validCmds := set.NewCmdSet()
	validCmds.Add("EXIT")
	validCmds.Add("ECHO")
	validCmds.Add("TYPE")
	line = line[:len(line)-1]
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
			fmt.Printf("%s: not found\n", args[1])
		}
	default:
		fmt.Print(cmd + ": command not found\n")
	}
	main()
}
