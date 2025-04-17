package cli

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func GetExecutedCommand() string {
	args := os.Args
	ourArgs := make([]string, 0, len(args))
	ourArgs = append([]string{}, path.Base(args[0]))
	ourArgs = append(ourArgs, args[1:]...)
	return strings.Join(ourArgs, " ")
}

func ClearScreen() {
	fmt.Print("\033c")
}

// UsedStdin can be checked: `(sleep 1; echo "some data") | ./main`
func UsedStdin() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func SetCursorPosition(col, line int) {
	fmt.Printf("\033[%d;%dH", line+1, col+1)
}

func SetCursor(posXY [2]int) {
	SetCursorPosition(posXY[1], posXY[0])
}

const (
	NC   = "\033[0m" // reset
	RED  = "\033[0;31m"
	GREY = "\033[1;30m"
	WARN = "\033[3;33m"
)

func ErrorStyle(s string) string {
	return RED + s + NC
}

func WarnStyle(s string) string {
	return WARN + s + NC
}

func ShadowStyle(s string) string {
	return GREY + s + NC
}
