package packages

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
