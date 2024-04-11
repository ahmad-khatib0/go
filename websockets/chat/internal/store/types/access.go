package types

// AccessMode is a definition of access mode bits.
type AccessMode uint

// DefaultAccess is a per-topic default access modes
type DefaultAccess struct {
	Auth AccessMode
	Anon AccessMode
}
