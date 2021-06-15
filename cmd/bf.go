package main

import (
	"io"
	"os"

	"github.com/michaelmcallister/gobrainfuckyourself/brainfuck"
)

func main() {
	// Use Stdin and Stdout for program interface.
	var rwc io.ReadWriter = struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	if len(os.Args) > 1 {
		bf := brainfuck.New(os.Args[1], rwc)
		bf.Run()
	}
}
