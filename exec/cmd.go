package exec

import (
	"bufio"
	"fmt"
	"os"
	goexec "os/exec"
	"regexp"
	"strings"

	"github.com/frozzare/shed/log"
)

// Cmd will execute a input cmd string.
func Cmd(input string, output bool) error {
	if output {
		log.Info("running: %s", input)
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	parts := strings.Fields(input)
	head := parts[0]
	parts = parts[1:len(parts)]

	cmd := goexec.Command(head, parts...)
	cmd.Dir = path

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	r := regexp.MustCompile("^(?:export|)\\s*([^\\d+][\\w_]+)\\s?=\\s?(.+)")
	go func() {
		i := 0
		for scanner.Scan() {
			text := scanner.Text()
			if r.MatchString(text) {
				s := r.FindStringSubmatch(text)
				if len(s) > 1 {
					os.Setenv(s[1], strings.Trim(s[2], "\""))
				}
			} else if output {
				if i == 0 {
					fmt.Println()
				}

				fmt.Println("  " + scanner.Text())

				i++
			}
		}

		if i > 0 {
			fmt.Println()
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

// CmdList executes a list of commands.
func CmdList(cmds []string, output bool) {
	for _, cmd := range cmds {
		if err := Cmd(cmd, output); err != nil {
			log.Error(err)
		}
	}
}
