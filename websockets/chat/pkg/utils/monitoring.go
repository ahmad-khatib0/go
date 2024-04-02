package utils

import (
	"strconv"
	"strings"
	"unicode"
)

// Base10Version() returns base10 app version
func (u *Utils) Base10Version(hex int) int64 {
	major := hex >> 16 & 0xFF
	minor := hex >> 8 & 0xFF
	trailer := hex & 0xFF
	return int64(major*10000 + minor*100 + trailer)
}

// Parses semantic version string in the following formats:
//
//	1.2, 1.2abc, 1.2.3, 1.2.3-abc, v0.12.34-rc5
//
// Unparceable values are replaced with zeros.
func (u *Utils) ParseBuildstampVersion(bs string) int {
	var major, minor, patch int
	bs = strings.TrimPrefix(bs, "v")

	// We can handle 3 parts only.
	parts := strings.SplitN(bs, ".", 3)
	count := len(parts)

	if count > 0 {
		major = u.parseOneSmvPart(parts[0])
	}

	if count > 1 {
		minor = u.parseOneSmvPart(parts[1])
	}

	if count > 2 {
		patch = u.parseOneSmvPart(parts[2])
	}

	return (major << 16) | (minor << 8) | patch
}

// Parse one component of a semantic version string.
func (u *Utils) parseOneSmvPart(vers string) int {
	end := strings.IndexFunc(vers, func(r rune) bool {
		return !unicode.IsDigit(r)
	})

	t := 0
	var err error
	if end > 0 {
		t, err = strconv.Atoi(vers[:end])
	} else if len(vers) > 0 {
		t, err = strconv.Atoi(vers)
	}

	if err != nil || t > 0x1fff || t <= 0 {
		return 0
	}

	return t
}
