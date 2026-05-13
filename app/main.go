package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	customDS "github.com/codecrafters-io/shell-starter-go/internal/custom"
	parser "github.com/codecrafters-io/shell-starter-go/internal/parser"
)

// returns a buffer containing stdout and stderr
func runCmd(name string, args ...string) ([]byte, []byte) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	_ = cmd.Run()
	// if stderrBuf.Len() > 0 {
	// 	fmt.Print(stderrBuf.String())
	// }
	//out, _ := cmd.CombinedOutput()
	return stdoutBuf.Bytes(), stderrBuf.Bytes()
}

func writeToFile(fileName string, output []byte) error {
	file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(output)
	return nil
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
	validCmds.Add("PWD")
	validCmds.Add("CD")

	var redirects bool
	var redirectErr bool
	var fileName string
	// if strings.Contains(line, ">") {
	// 	redirects = true
	// 	idx := strings.Index(line, ">")
	// 	fileName = strings.TrimSpace(line[idx+1:])
	// 	if idx > 0 && line[idx-1] != '1' {
	// 		line = line[:idx]
	// 	} else {
	// 		line = line[:idx-1]
	// 	}
	// }

	args := parser.ParseLine(line)

	for i, arg := range args {
		if arg == ">" || arg == "1>" {
			redirects = true
			fileName = args[i+1]
			args = args[:i]
			break
		} else if arg == "2>" {
			redirectErr = true
			fileName = args[i+1]
			args = args[:i]
			break
		}
	}
	cmd := args[0]
	switch strings.ToUpper(cmd) {
	case "EXIT":
		os.Exit(0)
	// case "ECHO":
	// 	for i := 1; i < len(args); i++ {
	// 		fmt.Print(args[i])
	// 		if i != len(args)-1 {
	// 			fmt.Print(" ")
	// 		}
	// 	}
	// 	fmt.Print("\n")
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
	case "PWD":
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("error: ", err)
			break
		}
		fmt.Println(dir)
	case "CD":
		cd := args[1]
		if cd == "~" {
			cd, err = os.UserHomeDir()
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		err := os.Chdir(cd)
		if err != nil {
			fmt.Printf("cd: %s: No such file or directory\n", args[1])
			break
		}
	// case "CAT": #TODO because a naive impleme
	default:
		_, err := exec.LookPath(args[0])
		if err != nil {
			fmt.Printf("%s: not found\n", args[0])
			break
		}

		// if slices.Contains(args, ">") {
		// 	redirects = true
		// 	idx := slices.Index(args, ">")
		// 	fileName = strings.Join(args[idx+1:], "")
		// 	args = args[:idx]
		// }
		stdOut, stdErr := runCmd(args[0], args[1:]...)

		if redirectErr {
			if err = writeToFile(fileName, stdErr); err != nil {
				panic(err)
			}
		} else {
			if len(stdErr) > 0 {
				fmt.Print(string(stdErr))
			}
		}

		if !redirects {
			fmt.Print(string(stdOut))
		} else {
			// dir := filepath.Dir(fileName)
			// err = os.Mkdir(dir, 0755)
			// if err != nil {
			// 	panic(err)
			// }
			if err = writeToFile(fileName, stdOut); err != nil {
				panic(err)
			}
		}
	}
	main()
}
