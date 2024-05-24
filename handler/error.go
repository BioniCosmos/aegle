package handler

import "fmt"

type ParseError struct {
	err error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("ParseError: %v", e.err)
}
