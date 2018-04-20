package goerror

import "fmt"

type ErrorData struct {
	*ErrorData
	error
}

func NewError(err error) *ErrorData {
	return &ErrorData{error: err}
}

func (e *ErrorData) Add(newErr error) {
	prevErr := &ErrorData{
		ErrorData: e.ErrorData,
		error:     e.error,
	}

	e.error = newErr
	e.ErrorData = prevErr
}

func (e *ErrorData) Error() string {
	return e.error.Error()
}

func (e *ErrorData) Cause() string {
	str := fmt.Sprintf("'%s'", e.error.Error())
	nextErr := e.ErrorData

	for nextErr != nil {
		str += fmt.Sprintf(", caused by '%s'", e.ErrorData.error.Error())
		nextErr = nextErr.ErrorData
	}
	return str
}