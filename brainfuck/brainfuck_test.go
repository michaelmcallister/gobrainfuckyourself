package brainfuck

import (
	"testing"
)

type fakeReadWriter struct {
	Got               string
	ReadErr, WriteErr error
}

func (f *fakeReadWriter) Read(p []byte) (n int, err error) {
	return len(p), f.ReadErr
}

func (f *fakeReadWriter) Write(p []byte) (n int, err error) {
	f.Got += string(p)
	return len(p), f.WriteErr
}

func Test(t *testing.T) {
	testCases := []struct {
		desc         string
		instructions string
		readErr      error
		writeErr     error
		wantErr      bool
		want         string
	}{
		{
			desc:         "test Hello World!",
			instructions: "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.",
			wantErr:      false,
			want:         "Hello World!\n",
		},
		{
			desc:         "test unmatched left brace causes error",
			instructions: "+++[[[[[[[[[[[[[",
			wantErr:      true,
			want:         "",
		},
		{
			desc:         "test unmatched right brace causes error",
			instructions: "+++[[[",
			wantErr:      true,
			want:         "",
		},
		{
			desc:         "test jumps",
			instructions: ">>[-]<<[->>+<<]",
			wantErr:      false,
			want:         "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			frw := &fakeReadWriter{
				ReadErr:  tC.readErr,
				WriteErr: tC.writeErr,
			}
			bf := New(tC.instructions, frw)
			if err := bf.Run(); err != nil && !tC.wantErr {
				t.Errorf("brainfuck.Run(%s) = %v, want error = %T", tC.instructions, err, tC.wantErr)
			}
			if frw.Got != tC.want {
				t.Errorf("brainfuck.Run(%s) = %s, want=%s", tC.instructions, frw.Got, tC.want)
			}
		})
	}
}
