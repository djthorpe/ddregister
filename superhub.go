package ddregister

////////////////////////////////////////////////////////////////////////////////
// TYPES

type KeyType int

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SUPERHUB_DOWNSTREAM KeyType = iota
	SUPERHUB_UPSTREAM
	SUPERHUB_UPSTREAM_EXT
	SUPERHUB_UPSTREAM_STATUS
	SUPERHUB_SIGNAL_QUALITY
	SUPERHUB_QOS
	SUPERHUB_QOS_FLOWS
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Superhub interface {
	Get(KeyType) error
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (k KeyType) String() string {
	switch k {
	case SUPERHUB_DOWNSTREAM:
		return "SUPERHUB_DOWNSTREAM"
	case SUPERHUB_UPSTREAM:
		return "SUPERHUB_UPSTREAM"
	case SUPERHUB_UPSTREAM_EXT:
		return "SUPERHUB_UPSTREAM_EXT"
	case SUPERHUB_UPSTREAM_STATUS:
		return "SUPERHUB_UPSTREAM_STATUS"
	case SUPERHUB_SIGNAL_QUALITY:
		return "SUPERHUB_SIGNAL_QUALITY"
	case SUPERHUB_QOS:
		return "SUPERHUB_QOS"
	case SUPERHUB_QOS_FLOWS:
		return "SUPERHUB_QOS_FLOWS"
	default:
		return "[?? Invalid KeyType]"
	}
}
