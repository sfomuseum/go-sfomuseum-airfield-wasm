package airlines

import (
	"fmt"
)

type NotFound struct{ Code string }

func (e NotFound) Error() string {
	return fmt.Sprintf("Airline '%s' not found", e.Code)
}

func (e NotFound) String() string {
	return e.Error()
}

type MultipleCandidates struct{ Code string }

func (e MultipleCandidates) Error() string {
	return fmt.Sprintf("Multiple candidates for airline '%s'", e.Code)
}

func (e MultipleCandidates) String() string {
	return e.Error()
}

func IsNotFound(e error) bool {

	switch e.(type) {
	case NotFound, *NotFound:
		return true
	default:
		return false
	}
}

func IsMultipleCandidates(e error) bool {

	switch e.(type) {
	case MultipleCandidates, *MultipleCandidates:
		return true
	default:
		return false
	}
}
