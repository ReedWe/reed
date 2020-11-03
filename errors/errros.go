package errors

import (
	"errors"
	"fmt"
	"strings"
)

func New(text string) error {
	return errors.New(text)
}

type wrapperError struct {
	msg    string
	detail []string
	data   map[string]interface{}
	stack  []StackFrame
	root   error
}

func (w wrapperError) Error() string {
	return w.msg
}


func wrap(err error, msg string, stackSkip int) error {
	if err == nil {
		return nil
	}

	werr, ok := err.(wrapperError)
	if !ok {
		werr.root = err
		werr.msg = err.Error()
		werr.stack = getStack(stackSkip+2, stackTraceSize)
	}
	if msg != "" {
		werr.msg = msg + ": " + werr.msg
	}

	return werr
}

func Wrap(err error, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return wrap(err, fmt.Sprint(a...), 1)
}

func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return wrap(err, fmt.Sprintf(format, a...), 1)
}


// WithDetail returns a new error that wraps
// err as a chain error messsage containing text
// as its additional context.
// Function Detail will return the given text
// when called on the new error value.
func WithDetail(err error, text string) error {
	if err == nil {
		return nil
	}
	if text == "" {
		return err
	}
	e1 := wrap(err, text, 1).(wrapperError)
	e1.detail = append(e1.detail, text)
	return e1
}

// WithDetailf is like WithDetail, except it formats
// the detail message as in fmt.Printf.
// Function Detail will return the formatted text
// when called on the new error value.
func WithDetailf(err error, format string, v ...interface{}) error {
	if err == nil {
		return nil
	}
	text := fmt.Sprintf(format, v...)
	e1 := wrap(err, text, 1).(wrapperError)
	e1.detail = append(e1.detail, text)
	return e1
}

// Detail returns the detail message contained in err, if any.
// An error has a detail message if it was made by WithDetail
// or WithDetailf.
func Detail(err error) string {
	wrapper, ok := err.(wrapperError)
	if !ok {
		return err.Error()
	}
	return strings.Join(wrapper.detail, "; ")
}