package auth

import "errors"

// ParseAuthLevel parses authentication level from a string.
func ParseAuthLevel(name string) Level {
	switch name {
	case "anon", "ANON":
		return LevelAnon
	case "auth", "AUTH":
		return LevelAuth
	case "root", "ROOT":
		return LevelRoot
	default:
		return LevelNone
	}
}

// String implements Stringer interface: gets human-readable name for a numeric authentication level.
func (l Level) String() string {
	s, err := l.MarshalText()
	if err != nil {
		return "unknown"
	}
	return string(s)
}

// MarshalText converts Level to a slice of bytes with the name of the level.
func (a Level) MarshalText() ([]byte, error) {
	switch a {
	case LevelNone:
		return []byte(""), nil
	case LevelAnon:
		return []byte("anon"), nil
	case LevelAuth:
		return []byte("auth"), nil
	case LevelRoot:
		return []byte("root"), nil
	default:
		return nil, errors.New("auth.Level: invalid level value")
	}
}