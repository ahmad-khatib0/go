package store

import (
	"errors"
	"strings"
	"time"

	auth "github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	st "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// GetAuthRecord takes a user ID and a authentication scheme name,
//
// fetches unique scheme-dependent identifier and authentication secret.
func (s *Store) GetAuthRecord(user st.Uid, scheme string) (string, auth.Level, []byte, time.Time, error) {
	unique, authLvl, secret, expires, err := s.adp.Auth().GetRecord(user, scheme)

	if err == nil {
		parts := strings.Split(unique, ":")
		if len(parts) > 1 {
			unique = parts[1]
		} else {
			err = st.ErrInternal
		}
	}

	return unique, authLvl, secret, expires, err
}

// GetAuthUniqueRecord takes a unique identifier and a authentication scheme name, fetches user ID and
// authentication secret.
func (s *Store) AuthGetUniqueRecord(scheme, unique string) (st.Uid, auth.Level, []byte, time.Time, error) {
	return s.adp.Auth().GetUniqueRecord(scheme + ":" + unique)
}

// AddAuthRecord creates a new authentication record for the given user.
func (s *Store) AddAuthRecord(uid st.Uid, authLvl auth.Level, scheme, unique string, secret []byte,
	expires time.Time) error {

	return s.adp.Auth().AddRecord(uid, scheme, scheme+":"+unique, authLvl, secret, expires)
}

// UpdateAuthRecord updates authentication record with a new secret and expiration time.
func (s *Store) UpdateAuthRecord(uid st.Uid, authLvl auth.Level, scheme, unique string,
	secret []byte, expires time.Time) error {

	return s.adp.Auth().UpdRecord(uid, scheme, scheme+":"+unique, authLvl, secret, expires)
}

// DelAuthRecords deletes user's auth records of the given scheme.
func (s *Store) DelAuthRecords(uid st.Uid, scheme string) error {
	return s.adp.Auth().DelScheme(uid, scheme)
}

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

// AuthGetAuthNames returns all addressable auth handler names, logical
//
// and hardcoded excluding those which are disabled like "basic:".
func (s *Store) AuthGetAuthNames() []string {
	if len(s.authHandlerNames) == 0 {
		return nil
	}

	allNames := make(map[string]struct{})

	for name := range s.authHandlers {
		allNames[name] = struct{}{}
	}
	for name := range s.authHandlerNames {
		allNames[name] = struct{}{}
	}

	var names []string
	for name := range allNames {
		if s.AuthGetLogicalAuthHandler(name) != nil {
			names = append(names, name)
		}
	}

	return names
}

// AuthGetLogicalAuthHandler returns an auth handler by logical name. If there is no
//
// handler by that logical name it tries to find one by the hardcoded name.
func (s *Store) AuthGetLogicalAuthHandler(name string) types.AuthHandler {
	name = strings.ToLower(name)
	if len(s.authHandlerNames) != 0 {
		if lname, ok := s.authHandlerNames[name]; ok {
			return s.authHandlers[lname]
		}
	}

	return s.authHandlers[name]
}

func (s *Store) GetAuthHandler(name string) types.AuthHandler {
	return s.authHandlers[strings.ToLower(name)]
}
