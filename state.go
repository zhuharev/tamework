package tamework

import (
	"context"
	"sync"
)

type State string

// StateStorage represent storage for chat state. It's useful for fill form.
type StateStorage interface {
	GetState(ctx context.Context, chatID, userID int) (State, error)
	SetState(ctx context.Context, chatID, userID int, state State) error
}

type memStateStorage struct {
	data map[int]map[int]State
	mu   sync.RWMutex
}

func newMemStateStorage() *memStateStorage {
	return &memStateStorage{
		data: map[int]map[int]State{},
	}
}

// GetState return state stored in memory
func (s *memStateStorage) GetState(ctx context.Context, chatID, userID int) (state State, err error) {
	s.mu.RLock()
	chatMap := s.data[chatID]
	if chatMap != nil {
		state = chatMap[userID]
	}
	s.mu.RUnlock()

	return
}

// SetState set state to memory
func (s *memStateStorage) SetState(ctx context.Context, chatID, userID int, state State) error {
	s.mu.Lock()
	chatMap := s.data[chatID]
	if chatMap == nil {
		s.data[chatID] = map[int]State{}
	}
	s.data[chatID][userID] = state
	s.mu.Unlock()
	return nil
}
