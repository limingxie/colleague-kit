package errs

import (
	"fmt"
	"go/build"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/labstack/echo"
)

type Err struct {
	status int
	// message holds an annotation of the error.
	message string
	// previous holds the previous error in the error stack, if any.
	previous error
	File     string
	Line     int
}

// New is a drop in replacement for the standard library errors module that records
// the location that the error is created.
//
// For example:
//    return errors.New("validation failed")
//
func New(message string) error {
	err := &Err{message: message}
	err.SetLocation(1)
	return err
}

// Trace adds the location of the Trace call to the stack.  The Cause of the
// resulting error is the same as the error parameter.  If the other error is
// nil, the result will be nil.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       return errors.Trace(err)
//   }
//
func Trace(other error) error {
	if other == nil {
		return nil
	}
	var status int = http.StatusInternalServerError
	if errResult, ok := other.(*echo.HTTPError); ok {
		status = errResult.Code
	} else if errResult, ok := other.(*Err); ok {
		status = errResult.Status()
	}
	err := &Err{previous: other, status: status}
	err.SetLocation(1)
	return err
}

// Annotate is used to add extra context to an existing error. The location of
// the Annotate call is recorded with the annotations. The file, line and
// function are also recorded.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       return errors.Annotate(err, "failed to frombulate")
//   }
//
func Annotate(other error, message string) error {
	if other == nil {
		return nil
	}
	err := &Err{
		previous: other,
		message:  message,
	}
	err.SetLocation(1)
	return err
}

func (e *Err) SetLocation(callDepth int) {
	_, file, line, _ := runtime.Caller(callDepth + 1)
	e.File = trimGoPath(file)
	e.Line = line
}

var goPath = build.Default.GOPATH
var srcDir = filepath.Join(goPath, "src")

func trimGoPath(filename string) string {
	return strings.TrimPrefix(filename, fmt.Sprintf("%s%s", srcDir, string(os.PathSeparator)))
}

func (e *Err) Location() (filename string, line int) {
	return e.File, e.Line
}

// Error implements error.Error.
func (e *Err) Error() string {
	err := e.previous
	switch {
	case err == nil:
		return e.message
	case e.message == "":
		return err.Error()
	}
	return fmt.Sprintf("%s: %v", e.message, err)
}

func (e *Err) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "%s", ErrorStack(e))
	case 's':
		fmt.Fprintf(s, "%s", e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	default:
		fmt.Fprintf(s, "%%!%c(%T=%s)", verb, e, e.Error())
	}
}

func ErrorStack(err error) string {
	return strings.Join(errorStack(err), "\n")
}

func errorStack(err error) []string {
	if err == nil {
		return nil
	}

	// We want the first error first
	var lines []string
	for {
		var buff []byte
		if err, ok := err.(locationer); ok {
			file, line := err.Location()
			// Strip off the leading GOPATH/src path elements.
			file = trimGoPath(file)
			if file != "" {
				buff = append(buff, fmt.Sprintf("%s:%d", file, line)...)
				buff = append(buff, ": "...)
			}
		}

		if cerr, ok := err.(wrapper); ok {
			message := cerr.Message()
			buff = append(buff, message...)

			err = cerr.Underlying()

		} else {
			buff = append(buff, err.Error()...)
			err = nil
		}

		lines = append(lines, string(buff))
		if err == nil {
			break
		}
	}

	// reverse the lines to get the original error, which was at the end of
	// the list, back to the start.[1,2,3] => [3,2,1]
	var result []string
	for i := len(lines); i > 0; i-- {
		result = append(result, lines[i-1])
	}

	return result
}

type locationer interface {
	Location() (string, int)
}

type wrapper interface {
	// Message returns the top level error message,
	// not including the message from the Previous
	// error.
	Message() string

	// Underlying returns the Previous error, or nil
	// if there is none.
	Underlying() error
}

func (e *Err) Underlying() error {
	return e.previous
}

func (e *Err) Message() string {
	return e.message
}

func (e *Err) Status() int {
	return e.status
}
