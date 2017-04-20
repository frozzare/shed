package log

import (
	"fmt"
	"os"
)

func Error(err error) {
	fmt.Printf("==> error: %s\n", err.Error())
	os.Exit(1)
}
