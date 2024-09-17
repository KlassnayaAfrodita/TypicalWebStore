package storage

import (
	"math/rand"
	"sync"
)

// var NotFound = errors.New("Not Found")

type Session struct {
	sessions map[string]int //! связываем с юзерами для middleware
	mu       *sync.RWMutex
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string { //* создание ключа сессии
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func NewSession() *Session {
	return &Session{
		sessions: map[string]int{},
		mu:       &sync.RWMutex{},
	}
}

func (s *Session) GetSession(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, ok := s.sessions[key]
	if !ok {
		return 0, NotFound
	}
	return userID, nil
}

func (s *Session) SetSession(userID int) (string, error) {
	SID := RandStringRunes(8)
	_, ok := s.sessions[SID]

	if ok {
		for {
			SID = RandStringRunes(8)
			_, ok := s.sessions[SID]
			if !ok {
				break
			}
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[SID] = userID
	return SID, nil
}

func (s *Session) DeleteSession(key string) {
	delete(s.sessions, key)
}
