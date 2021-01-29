package tamework

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/xid"
)

var (
	ErrNotFound = errors.New("not found")
)

type FormHandler func(*Context, []Answer)

type FormStore interface {
	// GetActiveForm return current user form is exist.
	GetActiveForm(ctx context.Context, chatID, userID int) (*Form, error)
	// SaveForm save anser to store.
	SaveForm(ctx context.Context, chatID, userID int, form *Form) error
	// DeleteForm remove form.
	DeleteForm(ctx context.Context, chatID, userID int) error
}

var (
	// build check
	_ FormStore = (*memFormStore)(nil)
)

type memFormStore struct {
	data map[int]map[int]*Form
	mu   sync.RWMutex
}

func newMemFormStore() *memFormStore {
	return &memFormStore{
		data: make(map[int]map[int]*Form),
	}
}

func (s *memFormStore) GetActiveForm(ctx context.Context, chatID, userID int) (f *Form, _ error) {
	s.mu.RLock()
	if m := s.data[chatID]; m == nil {
		// create map if not exists
		s.data[chatID] = map[int]*Form{}
	}
	f = s.data[chatID][userID]
	s.mu.RUnlock()

	if f == nil {
		return nil, ErrNotFound
	}

	return f, nil
}

func (s *memFormStore) SaveForm(ctx context.Context, chatID, userID int, form *Form) error {
	s.mu.Lock()
	if m := s.data[chatID]; m == nil {
		// create map if not exists
		s.data[chatID] = map[int]*Form{}
	}
	s.data[chatID][userID] = form
	s.mu.Unlock()

	return nil
}

func (s *memFormStore) DeleteForm(ctx context.Context, chatID, userID int) error {
	s.mu.Lock()
	if s.data[chatID] != nil && s.data[chatID][userID] != nil {
		delete(s.data[chatID], userID)
	}
	s.mu.Unlock()

	return nil
}

type Form struct {
	// Keyword is used in router for initial handler
	Keyword           string
	IsDone            bool
	CurrentQuestionID string
	Questions         []*Question
	Answers           Answers
}

func NewForm() *Form {
	return &Form{}
}

func (f *Form) AddQuestion(q *Question) *Form {
	f.Questions = append(f.Questions, q)
	return f
}

// MakeHandler create form handler.
func (f *Form) MakeHandler(store FormStore, handler FormHandler) Handler {
	return func(ctx *Context) {
		question := f.GetNextQuestion()

		// save answer for last question
		if question != nil {
			answer := Answer{
				Answer:    ctx.Text,
				Answered:  true,
				CreatedAt: time.Now(),
				Question:  *question,
			}
			f.Answers = append(f.Answers, answer)
			err := store.SaveForm(ctx.Context, int(ctx.ChatID), int(ctx.UserID), f)
			if err != nil {
				ctx.Send("some database error")
				return
			}
		}
		// get current question
		question = f.GetNextQuestion()
		if question == nil {
			// form is done, call answers handler
			handler(ctx, f.Answers)
			err := store.DeleteForm(ctx.Context, int(ctx.ChatID), int(ctx.UserID))
			if err != nil {
				ctx.Send("some database error")
				return
			}
			return
		}

		// build buttons
		kb := ctx.NewKeyboard(nil)
		for _, v := range question.Answers {
			kb.AddCallbackButton(v)
		}
		_, err := ctx.Markdown(question.Text)
		if err != nil {
			ctx.Send("network error")
			return
		}
	}
}

type Question struct {
	ID             string
	Text           string
	Answers        []string
	allowFreeInput bool
}

func NewQuestion(text string, answers []string) *Question {
	return &Question{
		ID:      xid.New().String(),
		Text:    text,
		Answers: answers,
	}
}

func NewQuestionWithFreeInput(text string, answers []string) *Question {
	return &Question{
		ID:             xid.New().String(),
		Text:           text,
		Answers:        answers,
		allowFreeInput: true,
	}
}

type Answer struct {
	Question  Question
	Answer    string
	Answered  bool
	CreatedAt time.Time
}

func (r *Router) makeFormHandler(handler FormHandler) Handler {
	return func(ctx *Context) {

	}
}

type Answers []Answer

func (a Answers) IsAnswered(questionID string) bool {
	answer, has := a.GetByQuestionID(questionID)
	if !has {
		return false
	}
	return answer.Answered
}

func (a Answers) GetByQuestionID(id string) (Answer, bool) {
	for _, v := range a {
		if v.Question.ID == id {
			return v, true
		}
	}
	return Answer{}, false
}

// NoOneAnswered all answers is empty
func (a Answers) NoOneAnswered() (res bool) {
	for _, v := range a {
		if v.Answered {
			return false
		}
	}
	return true
}

// AllAnswered all answers filled by user
func (a Answers) AllAnswered() (res bool) {
	for _, v := range a {
		if !v.Answered {
			return false
		}
	}
	return true
}

// GetNextQuestion check answered questions and get next question with empty answer
func (f *Form) GetNextQuestion() *Question {
	for _, q := range f.Questions {
		if !f.Answers.IsAnswered(q.ID) {
			return q
		}
	}
	return nil
}
