package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Printf("%s: command not found\n", line)
		}
	}
}
