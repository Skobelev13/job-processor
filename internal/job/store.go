package job

import "sync"

type Store struct {
	mu     sync.RWMutex
	status map[int]Status
}

func NewStore() *Store {
	return &Store{
		status: make(map[int]Status),
	}
}

func (s *Store) Set(id int, st Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status[id] = st
}

func (s *Store) Get(id int) (Status, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.status[id]
	return st, ok
}