package docker

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecCmd will execute a input cmd string.
func ExecCmd(input string, output bool) error {
	if output {
		fmt.Printf("==> running: %s\n", input)
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	parts := strings.Fields(input)
	head := parts[0]
	parts = parts[1:len(parts)]

	cmd := exec.Command(head, parts...)
	cmd.Dir = path

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if output {
		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				fmt.Println("  " + scanner.Text())
			}
		}()
	}

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
