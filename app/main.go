package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Print("$ ")
	cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
	cmd = cmd[:len(cmd)-1]
	if strings.ToUpper(cmd) == "EXIT" {
		os.Exit(0)
	}
	fmt.Print(cmd + ": command not found\n")
	main()
}
