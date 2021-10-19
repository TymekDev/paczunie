package packages

import (
	"strconv"

	"github.com/pkg/errors"
)

// Status denotes current status of a Pkg.
type Status int

const (
	// Ordered Status means that an order has been placed.
	Ordered Status = iota + 1
	// Shipped Status means that a package has been posted.
	Shipped
	// Inpost Status means that a packages has arrived at a parcel locker.
	Inpost
	// Delivered Status means that a package has been delivered.
	Delivered
)

// StatusFromInt converts an integer to a valid Status.
func StatusFromInt(x int) (Status, error) {
	status := Status(x)
	if status < Ordered || status > Delivered {
		const msg = "Status value (%d) out of range (%d - %d)"
		return 0, errors.Errorf(msg, status, Ordered, Delivered)
	}
	return status, nil
}

// StatusFromString converts a string to a valid Status. It is a wrapper on
// strconv.Atoi and StatusFromInt.
func StatusFromString(s string) (Status, error) {
	x, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return StatusFromInt(x)
}

func (s Status) String() string {
	return strconv.Itoa(int(s))
}
