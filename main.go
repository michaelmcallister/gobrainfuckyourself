package main

import (
	"io"
	"os"

	"github.com/michaelmcallister/gobrainfuckyourself/brainfuck"
)

func main() {
	// Hello World!
	c := "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."
	// Use Stdin and Stdout for program interface.
	var rwc io.ReadWriter = struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	bf := brainfuck.New(c, rwc)
	bf.Run()
}
