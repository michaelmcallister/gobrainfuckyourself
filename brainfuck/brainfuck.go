package brainfuck

import (
	"errors"
	"fmt"
	"io"
)

// BF contains the interal BF state.
type BF struct {
	RW    io.ReadWriter // ReadWriter used for program input and output.
	ptr   int           // current data pointer.
	pc    int           // program counter. Current location in the tape.
	depth int           // Used for bracket matching.
	tape  []OpCode      // Parsed operators.
	cells [30000]byte   // Memory layout.
}

// ErrUnmatchedBracket will be returned if the supplied instructions contain an
// unmatched bracket.
var ErrUnmatchedBracket = errors.New("unmatched bracket")

// OpCode represents each available command in brainfuck.
// See: https://en.wikipedia.org/wiki/Brainfuck#Commands
type OpCode rune

const (
	// IncrementPointer shifts the ptr right one.
	IncrementPointer OpCode = '>'
	// DecrementPointer shifts the ptr left one.
	DecrementPointer OpCode = '<'
	// IncrementByte increments the byte at ptr by one.
	IncrementByte OpCode = '+'
	// DecrementByte decrements the byte at ptr by one.
	DecrementByte OpCode = '-'
	// OutputByte will echo the byte at ptr.
	OutputByte OpCode = '.'
	// StoreByte will accept a single byte from stdin and store it at ptr.
	StoreByte OpCode = ','
	// JumpForward will move the ptr to the next ']' if the byte at ptr is 0.
	JumpForward OpCode = '['
	// JumpBackward will move the ptr to the preceeding '[' if the byte at ptr
	// is 0.
	JumpBackward OpCode = ']'
)

// New initialises a new brainfuck interpreter.
func New(input string, readwriter io.ReadWriter) *BF {
	var ins []OpCode
	for _, op := range input {
		switch OpCode(op) {
		case IncrementPointer:
			ins = append(ins, IncrementPointer)
		case DecrementPointer:
			ins = append(ins, DecrementPointer)
		case IncrementByte:
			ins = append(ins, IncrementByte)
		case DecrementByte:
			ins = append(ins, DecrementByte)
		case OutputByte:
			ins = append(ins, OutputByte)
		case StoreByte:
			ins = append(ins, StoreByte)
		case JumpForward:
			ins = append(ins, JumpForward)
		case JumpBackward:
			ins = append(ins, JumpBackward)
		}
	}
	return &BF{tape: ins, RW: readwriter}
}

// Step will evaluate the tape at the current program counter. Incrementing the
// program counter after a successful evaluation.
func (s *BF) Step() error {
	switch s.tape[s.pc] {
	case IncrementPointer:
		s.ptr++
	case DecrementPointer:
		s.ptr--
	case IncrementByte:
		s.cells[s.ptr]++
	case DecrementByte:
		s.cells[s.ptr]--
	case OutputByte:
		fmt.Fprintf(s.RW, "%c", s.cells[s.ptr])
	case StoreByte:
		var b = make([]byte, 1)
		s.RW.Read(b)
		s.cells[s.ptr] = b[0]
	case JumpForward:
		if s.cells[s.ptr] == 0 {
			c, err := s.findMatchingRightBracket()
			s.pc = c
			if err != nil {
				return err
			}
		} else {
			s.depth++
		}
	case JumpBackward:
		if s.cells[s.ptr] != 0 {
			c, err := s.findMatchingLeftBracket()
			s.pc = c
			if err != nil {
				return err
			}
		} else {
			s.depth--
		}
	}
	s.pc++
	return nil
}

func (s *BF) findMatchingRightBracket() (int, error) {
	found := false
	currentDepth := s.depth
	i := s.pc
	for i > 0 && !found {
		i++
		switch s.tape[i] {
		case JumpForward:
			currentDepth++
		case JumpBackward:
			if currentDepth == s.depth {
				found = true
			} else {
				currentDepth--
			}
		}
	}
	if found {
		return i, nil
	}
	return 0, ErrUnmatchedBracket
}

func (s *BF) findMatchingLeftBracket() (int, error) {
	found := false
	currentDepth := s.depth
	i := s.pc
	for i > 0 && !found {
		i--
		switch s.tape[i] {
		case JumpForward:
			if currentDepth == s.depth {
				found = true
			} else {
				currentDepth--
			}
		case JumpBackward:
			currentDepth++
		}
	}
	if found {
		return i, nil
	}
	return 0, ErrUnmatchedBracket
}

// Run will execute the tape that is loaded into memory.
func (s *BF) Run() error {
	for s.pc < len(s.tape) {
		if err := s.Step(); err != nil {
			return err
		}
	}
	return nil
}
