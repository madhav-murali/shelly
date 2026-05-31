package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	auto "github.com/codecrafters-io/shell-starter-go/internal/autocomplete"
	customDS "github.com/codecrafters-io/shell-starter-go/internal/custom"
	"github.com/codecrafters-io/shell-starter-go/internal/jobs"
	parser "github.com/codecrafters-io/shell-starter-go/internal/parser"
	"github.com/codecrafters-io/shell-starter-go/internal/pipes"

	"golang.org/x/term"
)

// returns a buffer containing stdout and stderr
func runCmd(oldState *term.State, jobManager *jobs.Manager, t *term.Terminal, args []string) ([]byte, []byte) {
	var stdoutBuf, stderrBuf bytes.Buffer
	originalCmd := strings.Join(args, " ")
	isBackground := false
	if len(args) > 0 && args[len(args)-1] == "&" {
		isBackground = true
		args = args[:len(args)-1]
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if isBackground {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		_ = cmd.Start()
		// if err != nil {
		// 	fmt.Fprintf(t, "Failed to start background process: %v\n", err)
		// 	return nil, nil
		// }

		jobId := jobManager.Add(cmd.Process.Pid, originalCmd)

		go func(id int) {
			cmd.Wait()
			jobManager.MarkDone(id)
		}(jobId)

		fmt.Fprintf(t, "[%d] %d\n", jobId, cmd.Process.Pid)
	} else {
		term.Restore(int(os.Stdin.Fd()), oldState)
		_ = cmd.Run()
		oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
	}
	// if stderrBuf.Len() > 0 {
	// 	fmt.Print(stderrBuf.String())
	// }
	//out, _ := cmd.CombinedOutput()
	return stdoutBuf.Bytes(), stderrBuf.Bytes()
}

func writeToFile(fileName string, output []byte, preserve bool) error {
	flags := os.O_CREATE | os.O_WRONLY
	if preserve {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	file, err := os.OpenFile(fileName, flags, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(output)
	if err != nil {
		return err
	}
	return nil
}

//func GetPathAndWalk(tr *auto.auto, )

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to set raw mode: ", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	//builtins := []string{"ECHO"}
	validCmds := customDS.NewCmdSet()
	validCmds.Add("EXIT")
	validCmds.Add("ECHO")
	validCmds.Add("TYPE")
	validCmds.Add("PWD")
	validCmds.Add("CD")
	validCmds.Add("JOBS")

	//jobs manager
	jb := jobs.NewManager()

	//command auto
	tr := auto.NewTrie()
	tr.AddAll(validCmds.ReturnAllLower())
	tr.IndexSystemPath()

	//File auto
	t := term.NewTerminal(os.Stdin, "$ ")

	err = t.SetSize(4096, 80)
	if err != nil {
		fmt.Printf("err : %v", err)
		os.Exit(1)
	}

	var tabCount int16
	var lastLine string
	t.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		if key == '\t' {
			if strings.Contains(line, " ") {
				//means we are searching for filename
				return auto.FileCompletetion(t, line, pos, key, &tabCount, &lastLine)
			}
			commands, ok := tr.HasPrefix(line)
			if !ok || len(commands) == 0 {
				fmt.Fprintf(t, "\x07")
				return line, pos, false
			}

			if len(commands) == 1 {
				tabCount = 0
				completed := commands[0] + " "
				return completed, len(completed), true
			}

			//when more than one commands found to be here

			commonPrefix := tr.LCP(line)

			if len(commonPrefix) > len(line) {
				tabCount = 0
				return commonPrefix, len(commonPrefix), true
			}

			if line != lastLine {
				tabCount = 0
			}

			tabCount++
			lastLine = line

			switch tabCount {
			case 1:
				fmt.Fprintf(t, "\x07")
				return line, pos, false
			case 2:
				tabCount = 0

				slices.Sort(commands)
				fmt.Fprintf(t, "$ %s", line)
				fmt.Fprintf(t, "\r\n%s\r\n", strings.Join(commands, "  "))

				return line, pos, false
			}
		}

		tabCount = 0
		lastLine = ""
		return line, pos, false
	}

	for {
		jb.ReapJobs(t)
		line, err := t.ReadLine()
		if err != nil {
			break
		}

		line = strings.ReplaceAll(line, "\r", "")

		if len(line) == 0 {
			continue
		}

		if strings.Contains(line, "|") {
			parsedCmds := pipes.ParseLineCmds(line)
			term.Restore(int(os.Stdin.Fd()), oldState)

			err := pipes.ExecutePipeline(parsedCmds)
			if err != nil {
				fmt.Fprintf(t, "pipeline error : %v", err)
			}
			oldState, _ = term.MakeRaw(int(os.Stdin.Fd()))
			continue
		}

		var redirects bool
		var redirectErr bool
		var fileName string
		var appends bool
		args := parser.ParseLine(line)

		if len(args) == 0 {
			continue
		}

		for i, arg := range args {
			if arg == ">" || arg == "1>" || arg == ">>" || arg == "1>>" {
				redirects = true
				if arg == ">>" || arg == "1>>" {
					appends = true
				}
				fileName = args[i+1]
				args = args[:i]
				break
			} else if arg == "2>" || arg == "2>>" {
				redirectErr = true
				if arg == "2>>" {
					appends = true
				}
				fileName = args[i+1]
				args = args[:i]
				break
			}
		}
		cmd := args[0]
		switch strings.ToUpper(cmd) {
		case "EXIT":
			term.Restore(int(os.Stdin.Fd()), oldState)
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
				fmt.Fprintf(t, "%s is a shell builtin\n", args[1])
			} else {
				path, err := exec.LookPath(args[1])
				if err != nil {
					fmt.Fprintf(t, "%s: not found\n", args[1])
					break
				}
				fmt.Fprintf(t, "%s is %s\n", args[1], path)
			}
		case "PWD":
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(t, "error: ", err)
				break
			}
			fmt.Fprintln(t, dir)
		case "CD":
			cd := args[1]
			if cd == "~" {
				cd, err = os.UserHomeDir()
				if err != nil {
					fmt.Fprintln(t, err)
					break
				}
			}
			err := os.Chdir(cd)
			if err != nil {
				fmt.Fprintf(t, "cd: %s: No such file or directory\n", args[1])
				//break
			}
		case "JOBS":
			jb.ListJobs(t)
		default:
			_, err := exec.LookPath(args[0])
			if err != nil {
				fmt.Fprintf(t, "%s: not found\n", args[0])
				break
			}

			stdOut, stdErr := runCmd(oldState, jb, t, args)

			if redirectErr {
				if err = writeToFile(fileName, stdErr, appends); err != nil {
					panic(err)
				}
			} else {
				if len(stdErr) > 0 {
					fmt.Fprint(t, string(stdErr))
				}
			}

			if !redirects {
				fmt.Fprint(t, string(stdOut))
			} else {
				// dir := filepath.Dir(fileName)
				// err = os.Mkdir(dir, 0755)
				// if err != nil {
				// 	panic(err)
				// }
				if err = writeToFile(fileName, stdOut, appends); err != nil {
					panic(err)
				}
			}
		}
	}
}
