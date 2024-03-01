package errors

import (
	"errors"
	"fmt"
	"github.com/user823/Sophie/pkg/ds"
	"io"
)

type fundamental struct {
	msg   string
	stack *ds.Stack
}

// 新建一个error对象
func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: ds.Callers(),
	}
}

func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: ds.Callers(),
	}
}

func (f *fundamental) Error() string { return f.msg }
func (f *fundamental) Unwrap() error { return nil }
func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			fmt.Fprintf(s, "%+v", f.stack)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

type withStack struct {
	error
	stack *ds.Stack
}

// 对error 封装调用栈信息
// 返回新的error
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		error: err,
		stack: ds.Callers(),
	}
}

func (w *withStack) Unwrap() error { return w.error }
func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Unwrap())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	e := &withMessage{err, message}
	return &withStack{e, ds.Callers()}
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	e := &withMessage{err, fmt.Sprintf(format, args...)}
	return &withStack{e, ds.Callers()}
}

func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{err, message}
}

func WithMessagef(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &withMessage{err, fmt.Sprintf(format, args...)}
}

// 对error 封装额外的msg 信息用于拼接
// 返回新的error
type withMessage struct {
	error
	msg string
}

func (w *withMessage) Error() string { return fmt.Sprintf("%s: %s", w.msg, w.error.Error()) }
func (w *withMessage) Unwrap() error { return w.error }
func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.error)
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func WithCodeMessage(err error, code int, message string) error {
	if err == nil {
		return nil
	}
	return &withCodeMessage{err, code, message, ds.Callers()}
}

func WithCodeMessagef(err error, code int, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &withCodeMessage{err, code, fmt.Sprintf(format, args...), ds.Callers()}
}

// 对error 封装code 和 msg 信息
// 返回新的error
type withCodeMessage struct {
	error
	code  int
	msg   string
	stack *ds.Stack
}

func CodeMessage(code int, msg string) error {
	return &withCodeMessage{
		code:  code,
		msg:   msg,
		stack: ds.Callers(),
	}
}

func (w *withCodeMessage) Error() string {
	return fmt.Sprintf("CodeMessageError{code %d, message %s}: %s", w.code, w.msg, w.error.Error())
}
func (w *withCodeMessage) Unwrap() error { return w.error }
func (w *withCodeMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\nCodeMessageError{code %d, message %s}", w.Unwrap(), w.code, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

// 找到最初的error
func Cause(err error) error {
	for err != nil {
		w, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}

		e := w.Unwrap()
		if e == nil {
			break
		}

		err = e
	}
	return err
}

func Is(err, target error) bool             { return errors.Is(err, target) }
func As(err error, target interface{}) bool { return errors.As(err, target) }
