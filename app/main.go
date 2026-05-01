package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	for {
		fmt.Print("$ ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("%s: command not found\n", line)
		}
	}
}
