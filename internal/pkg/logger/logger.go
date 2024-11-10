package logger

import (
	"fmt"
	"os"
)

func Fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}
