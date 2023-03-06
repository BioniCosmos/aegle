package subscription

import (
	"errors"
	"fmt"
)

type Error struct {
	err error
}

func (e *Error) Error() string {
	return fmt.Sprintf("failed to generate subscription: %v", e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}

var ErrNoProfile = errors.New("empty profile")
var ErrNoOutbound = errors.New("empty outbound setting")
var ErrNoProtocol = errors.New("no protocol specified")

type UnknownProtocolError string

func (e UnknownProtocolError) Error() string {
	return fmt.Sprintf("unknown protocol: %q", string(e))
}

type ParseHeaderError struct {
	protocol string
	err      error
}

func (e *ParseHeaderError) Error() string {
	return fmt.Sprintf("failed to parse the header in %v settings: %v", e.protocol, e.err)
}

func (e *ParseHeaderError) Unwrap() error {
	return e.err
}

var ErrNoQuicSecurity = errors.New("security item not specified for QUIC: specify a security item except `none` or remove the key")
var ErrNoQuicKey = errors.New("no key set for QUIC: set a key or set `security` to `none`")
