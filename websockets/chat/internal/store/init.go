package store

import (
	"errors"
	"strings"
)

// InitAuthLogicalNames() initializes authentication mapping "logical handler name":"actual handler name".
//
// Logical name must not be empty, actual name could be an empty string.
//
// Registered authentication handlers (what are currently supported)
func (s *Store) InitAuthLogicalNames(ln []string) error {
	if len(ln) == 0 {
		return nil
	}

	if s.authHandlerNames == nil {
		s.authHandlerNames = make(map[string]string)
	}

	for _, pair := range ln {
		if parts := strings.Split(pair, ":"); len(parts) == 2 {

			if parts[0] == "" {
				return errors.New("store: empty logical auth name '" + pair + "'")
			}

			parts[0] = strings.ToLower(parts[0])
			if _, ok := s.authHandlerNames[parts[0]]; !ok {
				return errors.New("store: duplicate mapping for logical auth name '" + pair + "'")
			}

			parts[1] = strings.ToLower(parts[1])
			if parts[1] != "" {
				if _, ok := s.authHandlers[parts[1]]; !ok {
					return errors.New("store: unknown handler for logical auth name '" + pair + "'")
				}
			}

			if parts[0] == parts[1] {
				// NOTE: skip useless identity mapping
				continue
			}

			s.authHandlerNames[parts[0]] = parts[1]
		} else {
			return errors.New("store: invalid logical auth mapping '" + pair + "'")
		}
	}

	return nil
}
