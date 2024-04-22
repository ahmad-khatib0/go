package store

import (
	"strings"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/validate"
)

// GetValidator returns registered validator by name.
func (s *Store) GetValidator(name string) validate.Validator {
	return s.validators[strings.ToLower(name)]
}
