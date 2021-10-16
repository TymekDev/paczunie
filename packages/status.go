package packages

import "github.com/pkg/errors"

// Status denotes current status of a Pkg.
type Status int

const (
	// Ordered Status means that an order has been placed.
	Ordered Status = iota + 1
	// Shipped Status means that a package has been posted.
	Shipped
	// Delivered Status means that a package has been delivered.
	Delivered
)

// ToStatus converts an integer to a valid Status.
func ToStatus(x int) (Status, error) {
	status := Status(x)
	if status < Ordered || status > Delivered {
		const msg = "Status value (%d) out of range (%d - %d)"
		return 0, errors.Errorf(msg, status, Ordered, Delivered)
	}
	return status, nil
}
