package store

import "github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"

func (s *Store) DBGetAdapterName() string {
	return s.adp.DB().GetName()
}

func (s *Store) DBGetAdapterVersion() int {
	return s.adp.DB().Version()
}

func (s *Store) DBClose() error {
	if !s.adp.DB().IsOpen() {
		return nil
	}
	return s.adp.DB().Close()
}

// DBStats() returns a callback returning db connection stats object.
func (s *Store) DBStats() func() interface{} {
	return s.adp.DB().Stats
}

func (s Store) Adp() types.Adapter {
	return s.adp
}
