package tamework

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestRouter_EditedMessage(t *testing.T) {
	t.Skip("skipping (change token to run)")
	// Create a new Tamework instance
	tw, err := New("test-token")
	if err != nil {
		t.Fatalf("Failed to create Tamework instance: %v", err)
	}
	router := NewRouter(tw)

	// Track if the handler was called
	handlerCalled := false
	router.EditedMessage(func(c *Context) {
		handlerCalled = true
		if c.Text != "edited text" {
			t.Errorf("Expected text 'edited text', got '%s'", c.Text)
		}
	})

	// Create an edited message update
	update := tgbotapi.Update{
		EditedMessage: &tgbotapi.Message{
			MessageID: 123,
			From: &tgbotapi.User{
				ID:        456,
				FirstName: "TestUser",
			},
			Chat: &tgbotapi.Chat{
				ID: 789,
			},
			Text: "edited text",
		},
	}

	// Handle the update
	router.Handle(update)

	// Assert that the handler was called
	if !handlerCalled {
		t.Error("Edited message handler should have been called")
	}
}
