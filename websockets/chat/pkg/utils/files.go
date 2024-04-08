package utils

import "path/filepath"

// Convert relative filepath to absolute.
func (u *Utils) ToAbsolutePath(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(base, path))
}
