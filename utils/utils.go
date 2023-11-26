package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func PrintErr(text string) {
	red := "\033[31m"
	reset := "\033[0m"
	fmt.Println()
	fmt.Println(red + text + reset)
}

func ParsePosition(pos string) (int, int) {
	col := int(pos[0] - 'A')
	row, _ := strconv.Atoi(string(pos[1]))
	return col, row - 1
}

func ClearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
