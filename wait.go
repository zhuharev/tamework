package tamework

import (
	"sync"
	"time"
)

type Waiter struct {
	mu sync.RWMutex
	m  map[int64]chan string

	waitTimeout time.Duration
}

func NewWaiter(waitTimeout time.Duration) *Waiter {
	return &Waiter{
		m:           make(map[int64]chan string),
		waitTimeout: waitTimeout,
	}
}

func (w *Waiter) NeedNext(chatID int64, text string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if ch, ok := w.m[chatID]; ok {
		ch <- text
		delete(w.m, chatID)
		return false
	}
	return true
}

// TODO: wait document photo video audio etcs
func (w *Waiter) Wait(chatID int64) (string, bool) {
	w.mu.Lock()
	w.m[chatID] = make(chan string, 1)
	w.mu.Unlock()

	select {
	case <-time.After(w.waitTimeout):
		w.mu.Lock()
		delete(w.m, chatID)
		w.mu.Unlock()
		return "", false
	case text := <-w.m[chatID]:
		return text, true
	}
}
