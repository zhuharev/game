package errors

import "encoding/json"

// Error represents app specifics errors
type Error struct {
	code    int
	message string
}

// New return new error
func New(code int, message string) Error {
	return Error{code, message}
}

// Error string of error
func (err Error) Error() string {
	return err.message
}

// Code code of error
func (err Error) Code() int {
	return err.code
}

type errResponseJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MarshallJSON json.Marshaller interface
func (err Error) MarshallJSON() ([]byte, error) {
	return json.Marshal(errResponseJSON{
		Code:    err.code,
		Message: err.message,
	})
}
