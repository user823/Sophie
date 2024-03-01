package ds

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func formatFrame(f *runtime.Frame, s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.Function)
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.File)
		default:
			io.WriteString(s, path.Base(f.File))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.Line))
	case 'n':
		io.WriteString(s, funcname(f.Function))
	case 'v':
		formatFrame(f, s, 's')
		io.WriteString(s, ":")
		formatFrame(f, s, 'd')
	}
}

type Stack struct {
	fs *runtime.Frames
}

func (s *Stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for {
				f, more := s.fs.Next()
				io.WriteString(st, "\n")
				formatFrame(&f, st, verb)
				if !more {
					break
				}
			}
		}
	}
}

func Callers() *Stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return &Stack{
		fs: runtime.CallersFrames(pcs[:n]),
	}
}

func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
