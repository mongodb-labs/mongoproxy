package messages

type Requester interface {

	// Type returns a string that identifies the structure of this particular
	// Requester, which can then be used to cast it into its proper struct
	// type to examine its fields
	Type() string

	// ID returns the request ID of the original request, which can be used to
	// uniquely identify requests, and may correspond to the id used by the wire protocol
	ID() int32
}
