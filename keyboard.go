package tamework

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	// DefaultKeyboardRowLen represend max row len
	// if buttons count > DefaultKeyboardRowLen,
	// new rows will be created automaticaly
	DefaultKeyboardRowLen = 3
)

// KeyboardType represend type of keyboard
type KeyboardType string

const (
	// KeyboardInline inline kb
	KeyboardInline KeyboardType = "inline"
	// KeyboardReply reply kb
	KeyboardReply = "reply"
)

// Keyboard helper for play with keyboards
type Keyboard struct {
	rowLen         int
	values         interface{}
	inlineKeyboard tgbotapi.InlineKeyboardMarkup
	typ            KeyboardType
	remove         bool

	enabled bool
}

// NewKeyboard returns new Keyboard
func NewKeyboard(values interface{}) *Keyboard {
	enabled := false
	if values != nil {
		enabled = true
	}
	return &Keyboard{values: values, rowLen: DefaultKeyboardRowLen, typ: KeyboardReply, enabled: enabled}
}

// SetRowLen set max row length
func (k *Keyboard) SetRowLen(l int) {
	k.enabled = true
	k.rowLen = l
}

// SetType set one of two types
func (k *Keyboard) SetType(typ KeyboardType) {
	k.enabled = true
	k.typ = typ
}

// Remove indicate that you need to delete keyboard
func (k *Keyboard) Remove() *Keyboard {
	k.enabled = true
	k.remove = true
	return k
}

// Reset the keyboard
func (k *Keyboard) Reset() *Keyboard {
	k.values = nil
	k.enabled = false
	return k
}

func (k *Keyboard) AddContactButton(text string) *Keyboard {
	return k.addReplyButton(text, true)
}

// AddURLButton add inline button
func (k *Keyboard) AddURLButton(text, uri string) *Keyboard {
	return k.addInlineButton(text, uri, "url")
}

// AddCallbackButton add inline buton
func (k *Keyboard) AddCallbackButton(text string, datas ...string) *Keyboard {
	data := text
	if len(datas) > 0 {
		data = datas[0]
	}
	return k.addInlineButton(text, data, "cb")
}

// addInlineButton just helper
func (k *Keyboard) addInlineButton(text, data, typ string) *Keyboard {
	k.enabled = true
	k.typ = KeyboardInline

	// for pointer
	datacopy := data

	var button = tgbotapi.InlineKeyboardButton{
		Text: text,
	}
	if typ == "cb" {
		button.CallbackData = &datacopy
	} else if typ == "url" {
		button.URL = &datacopy
	}

	if k.values == nil {
		k.values = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{
					button,
				},
			},
		}
	} else {
		kb, ok := k.values.(tgbotapi.InlineKeyboardMarkup)
		if !ok {
			return k
		}
		if text == "" {
			kb.InlineKeyboard = append(kb.InlineKeyboard, []tgbotapi.InlineKeyboardButton{})
		} else {
			kb.InlineKeyboard[len(kb.InlineKeyboard)-1] = append(kb.InlineKeyboard[len(kb.InlineKeyboard)-1],
				button)
		}
		k.values = kb
	}
	return k
}

func (k *Keyboard) AddReplyButton(text string) *Keyboard {
	return k.addReplyButton(text, false)
}

// addReplyButton add an reply button
func (k *Keyboard) addReplyButton(text string, isContact bool) *Keyboard {

	k.enabled = true

	k.typ = KeyboardReply
	if k.values == nil {
		k.values = tgbotapi.ReplyKeyboardMarkup{
			Keyboard: [][]tgbotapi.KeyboardButton{
				{
					{
						Text:           text,
						RequestContact: isContact,
					},
				},
			},
			ResizeKeyboard: true,
		}
	} else {
		if kb, ok := k.values.(tgbotapi.ReplyKeyboardMarkup); ok {
			rows := len(kb.Keyboard)
			if rows == 0 {
				// it not possible
				return k
			}
			if text == "" {
				kb.Keyboard = append(kb.Keyboard, []tgbotapi.KeyboardButton{})
			} else {
				kb.Keyboard[rows-1] = append(kb.Keyboard[rows-1], tgbotapi.KeyboardButton{
					Text: text,
				})
			}
			k.values = kb
		}
	}
	return k
}

// Markup return interface for which can be used in tgbotapi.Message.Markup
func (k *Keyboard) Markup() interface{} {

	if k.remove {
		return tgbotapi.NewRemoveKeyboard(false)
	}

	if !k.enabled {
		return nil
	}

	switch v := k.values.(type) {
	case tgbotapi.ReplyKeyboardMarkup:
		return v
	case tgbotapi.InlineKeyboardMarkup:
		return v
	case []string:
		if k.typ == KeyboardReply {
			keyboard := [][]tgbotapi.KeyboardButton{{}}
			curRow := 0
			for _, button := range v {
				if len(keyboard[curRow]) == k.rowLen || button == "" {
					keyboard = append(keyboard, []tgbotapi.KeyboardButton{})
					curRow++
				}
				if button != "" {
					keyboard[curRow] = append(keyboard[curRow], tgbotapi.KeyboardButton{Text: button})
				}
			}
			return tgbotapi.ReplyKeyboardMarkup{
				Keyboard:       keyboard,
				ResizeKeyboard: true,
			}
		}
		keyboard := tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{}},
		}
		curRow := 0
		for _, button := range v {
			if len(keyboard.InlineKeyboard[curRow]) == k.rowLen || button == "" {
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{})
				curRow++
			}
			if button != "" {
				btn := button
				keyboard.InlineKeyboard[curRow] = append(keyboard.InlineKeyboard[curRow],
					tgbotapi.InlineKeyboardButton{Text: button, CallbackData: &btn})
			}
		}
		return keyboard
	case string:
		return tgbotapi.ReplyKeyboardMarkup{
			Keyboard:       [][]tgbotapi.KeyboardButton{{{Text: v}}},
			ResizeKeyboard: true,
		}

	}

	return nil
}

// InlineKeyboardMarkup returns keyboard typed tgbotapi.InlineKeyboardMarkup
func (k *Keyboard) InlineKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	if k.Markup() == nil {
		return tgbotapi.InlineKeyboardMarkup{}
	}
	if markup, ok := k.Markup().(tgbotapi.InlineKeyboardMarkup); ok {
		return markup
	}
	return tgbotapi.InlineKeyboardMarkup{}
}
