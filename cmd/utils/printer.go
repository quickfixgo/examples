package utils

import (
	"bufio"
	"fmt"
	"io"

	"github.com/fatih/color"
)

func PrintGood(line string) {
	color.Set(color.Bold, color.FgHiGreen, color.Underline)
	fmt.Println(line)
	color.Unset()
}

func PrintBad(line string) {
	color.Set(color.Bold, color.FgHiRed, color.Underline)
	fmt.Println(line)
	color.Unset()
}

func PrintInfo(line string) {
	color.Set(color.Bold, color.FgHiCyan)
	fmt.Println(line)
	color.Unset()
}

func PrintConfig(side string, reader io.Reader) {
	PrintInfo(fmt.Sprintf("starting FIX %s with config:", side))

	scanner := bufio.NewScanner(reader)
	color.Set(color.FgHiBlue)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(" " + line)
	}
	color.Unset()
}
