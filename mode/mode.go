package mode

// Kind represents the mode an logger should run in -
// for example: Dev or Pro
type Kind int8

// Possible mode Kinds for a logger
const (
	None Kind = iota + 1
	Development
	Production
)
