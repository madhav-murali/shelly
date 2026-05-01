package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		return
	}
	fmt.Print("$ " + cmd[:len(cmd)-1] + ": command not found")
	main()
}
