package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Print("$ ")
	line, err := (bufio.NewReader(os.Stdin).ReadString('\n'))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
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
	default:
		fmt.Print(cmd + ": command not found\n")
	}
	main()
}
