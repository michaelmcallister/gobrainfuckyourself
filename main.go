package main

import (
	"github.com/michaelmcallister/gobrainfuckyourself/brainfuck"
)

func main() {
	// Hello World!
	c := "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."
	bf := brainfuck.New(c)
	bf.Run()
}
