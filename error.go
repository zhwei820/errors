// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENSE file for details.

package errors

import (
	"fmt"
	"reflect"
	"runtime"
)

// Err holds a description of an error along with information about
// where the error was created.
//
// It may be embedded in custom error types to add extra information that
// this errors package can understand.
type Err struct {
	// message holds an annotation of the error.
	message string

	// cause holds the cause of the error as returned
	// by the Cause method.
	cause error

	// previous holds the previous error in the error stack, if any.
	previous error

	// file and line hold the source code location where the error was
	// created.
	file string
	line int
}

// Location is the file and line of where the error was most recently
// created or annotated.
func (e *Err) Location() (filename string, line int) {
	return e.file, e.line
}

// Underlying returns the previous error in the error stack, if any. A client
// should not ever really call this method.  It is used to build the error
// stack and should not be introspected by client calls.  Or more
// specifically, clients should not depend on anything but the `Cause` of an
// error.
func (e *Err) Underlying() error {
	return e.previous
}

// Cause returns the most recent error in the error stack that
// meets one of these criteria: the original error that was raised; the new
// error that was passed into the Wrap function; the most recently masked
// error; or nil if the error itself is considered the Cause.  Normally this
// method is not invoked directly, but instead through the Cause stand alone
// function.
func (e *Err) Cause() error {
	return e.cause
}

// Is .
func (e *Err) Is(err error) bool {
	return sameError(e.cause, err)
}

// Unwrap .
func (e *Err) Unwrap() error {
	return e.cause
}

// Message returns the message stored with the most recent location. This is
// the empty string if the most recent call was Trace, or the message stored
// with Annotate or Mask.
func (e *Err) Message() string {
	return e.message
}

// Error implements error.Error.
func (e *Err) Error() string {
	// We want to walk up the stack of errors showing the annotations
	// as long as the cause is the same.
	err := e.previous
	if !sameError(Cause(err), e.cause) && e.cause != nil {
		err = e.cause
	}
	switch {
	case err == nil:
		return e.message
	case e.message == "":
		return err.Error()
	}
	return fmt.Sprintf("%s: %v", e.message, err)
}

// Format implements fmt.Formatter
// When printing errors with %+v it also prints the stack trace.
// %#v unsurprisingly will print the real underlying type.
func (e *Err) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			fmt.Fprintf(s, "%s", ErrorStack(e))
			return
		case s.Flag('#'):
			// avoid infinite recursion by wrapping e into a type
			// that doesn't implement Formatter.
			fmt.Fprintf(s, "%#v", (*unformatter)(e))
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%s", e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	default:
		fmt.Fprintf(s, "%%!%c(%T=%s)", verb, e, e.Error())
	}
}

// helper for Format
type unformatter Err

func (unformatter) Format() { /* break the fmt.Formatter interface */ }

// SetLocation records the source location of the error at callDepth stack
// frames above the call.
func (e *Err) SetLocation(callDepth int) {
	_, file, line, _ := runtime.Caller(callDepth + 1)
	e.file = trimSourcePath(file)
	e.line = line
}

// StackTrace returns one string for each location recorded in the stack of
// errors. The first value is the originating error, with a line for each
// other annotation or tracing of the error.
func (e *Err) StackTrace() []string {
	return errorStack(e)
}

// Ideally we'd have a way to check identity, but deep equals will do.
func sameError(e1, e2 error) bool {
	//nolint
	return reflect.DeepEqual(e1, e2)
}
