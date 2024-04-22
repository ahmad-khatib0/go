package types

import (
	"errors"
	"strconv"
)

// Feature is a bitmap of authenticated features, such as validated/not validated.
type Feature uint16

const (
	// FeatureValidated bit is set if user's credentials are already validated (V).
	FeatureValidated Feature = 1 << iota
	// FeatureNoLogin is set if the token should not be used to permanently authenticate a session (L).
	FeatureNoLogin
)

// MarshalText converts Feature to ASCII byte slice.
func (f Feature) MarshalText() ([]byte, error) {
	res := []byte{}
	for i, chr := range []byte{'V', 'L'} {
		if (f & (1 << uint(i))) != 0 {
			res = append(res, chr)
		}
	}
	return res, nil
}

// UnmarshalText parses Feature string as byte slice.
func (f *Feature) UnmarshalText(b []byte) error {
	var f0 int
	var err error

	if len(b) > 0 {
		if b[0] >= '0' && b[0] <= '9' {
			f0, err = strconv.Atoi(string(b))
		} else {
		Loop:
			for i := 0; i < len(b); i++ {
				switch b[i] {
				case 'V', 'v':
					f0 |= int(FeatureValidated)
				case 'L', 'l':
					f0 |= int(FeatureNoLogin)
				default:
					err = errors.New("invalid character '" + string(b[i]) + "' <Feature(UnmarshalText)>")
					break Loop
				}
			}
		}
	}

	*f = Feature(f0)
	return err
}

// String Featureto a string representation.
func (f Feature) String() string {
	res, err := f.MarshalText()
	if err != nil {
		return ""
	}
	return string(res)
}

// MarshalJSON converts Feature to a quoted string.
func (f Feature) MarshalJSON() ([]byte, error) {
	res, err := f.MarshalText()
	if err != nil {
		return nil, err
	}

	return append(append([]byte{'"'}, res...), '"'), nil
}

// UnmarshalJSON reads Feature from a quoted string or an integer.
func (f *Feature) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		return f.UnmarshalText(b[1 : len(b)-1])
	}
	return f.UnmarshalText(b)
}
