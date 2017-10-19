package tamework

import (
	"sync"
	"time"
)

// Waiter used for wait user inputs
type Waiter struct {
	mu sync.RWMutex
	m  map[int64]chan Update

	waitTimeout time.Duration
}

// NewWaiter returns Waiter
func NewWaiter(waitTimeout time.Duration) *Waiter {
	return &Waiter{
		m:           make(map[int64]chan Update),
		waitTimeout: waitTimeout,
	}
}

// NeedNext check is Waiter waiting for input
func (w *Waiter) NeedNext(chatID int64, update Update) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if ch, ok := w.m[chatID]; ok {
		ch <- update
		delete(w.m, chatID)
		return false
	}
	return true
}

// Wait block the input and redirect one update to Waiter
func (w *Waiter) Wait(chatID int64, stopWord string, durations ...time.Duration) (Update, bool) {
	w.mu.Lock()
	w.m[chatID] = make(chan Update, 1)
	w.mu.Unlock()

	waitTimeout := w.waitTimeout
	if len(durations) > 0 {
		waitTimeout = durations[0]
	}

	select {
	case <-time.After(waitTimeout):
		w.mu.Lock()
		delete(w.m, chatID)
		w.mu.Unlock()
		return Update{}, false
	case u := <-w.m[chatID]:
		if u.Text() != "" && u.Text() == stopWord {
			return u, false
		}
		return u, true
	}
}

// Waiterer middleware for *Tamework
func Waiterer() Handler {
	return func(c *Context) {
		if !c.tamework.waiter.NeedNext(c.ChatID, c.update) {
			c.exited = true
			return
		}
	}
}
