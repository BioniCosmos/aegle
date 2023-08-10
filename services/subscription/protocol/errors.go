package protocol

import (
	"errors"
	"fmt"
)

type ParseSettingsError struct {
	protocol string
	err      error
}

func (e *ParseSettingsError) Error() string {
	return fmt.Sprintf("failed to parse %v settings: %v", e.protocol, e.err)
}

func (e *ParseSettingsError) Unwrap() error {
	return e.err
}

type ParseAccountError struct {
	protocol string
	err      error
}

func (e *ParseAccountError) Error() string {
	return fmt.Sprintf("failed to parse the %v account: %v", e.protocol, e.err)
}

func (e *ParseAccountError) Unwrap() error {
	return e.err
}

type IllegalPortError uint16

func (e IllegalPortError) Error() string {
	return fmt.Sprintf("a port should be an integer between 1 and 65535, but the port is %v now", uint16(e))
}

var ErrNoId = errors.New("no id or password specified")
var ErrNoHost = errors.New("no address or port specified")
