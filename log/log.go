package log

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	spaces = 7
)

// Error outputs a error row and exits with code 1 if no second bool parameter is provided.
func Error(err error, exit ...bool) {
	red := color.New(color.FgRed).SprintFunc()
	msg := space(fmt.Sprintf("error: %s", err.Error()))
	fmt.Printf("%s %s\n", red("==>"), msg)
	if len(exit) == 0 {
		os.Exit(1)
	}
}

// Info outputs a info row.
func Info(str string, a ...interface{}) {
	green := color.New(color.FgGreen).SprintFunc()
	msg := space(fmt.Sprintf(str, a...))
	fmt.Printf("%s %s\n", green("==>"), msg)
}

func space(s string) string {
	r := regexp.MustCompile("(\\w+)\\:")
	m := r.FindStringSubmatch(s)

	if len(m) < 2 {
		return s
	}

	for i := 0; i < spaces-len(m[1]); i++ {
		s = " " + s
	}

	return strings.Replace(s, "\n", "", -1)
}
