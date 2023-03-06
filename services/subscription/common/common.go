package common

import "fmt"

type UnknownSecurityError string

func (e UnknownSecurityError) Error() string {
	return fmt.Sprintf("unknown security type: %q", string(e))
}

func IsNone(item string) bool {
	return item == "" || item == "none"
}
