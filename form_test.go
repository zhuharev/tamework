package tamework

import "testing"

func TestFromBuilderExample(t *testing.T) {
	_ = NewForm().
		AddQuestion(NewQuestion(
			"How are you?",
			[]string{"fine", "sad"}),
		).
		AddQuestion(NewQuestion(
			"Your name is:",
			[]string{"John", "Doe"},
		))
}
