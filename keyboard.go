package tamework

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	DefaultKeyboardRowLen = 3
)

type KeyboardType string

const (
	KeyboardInline KeyboardType = "inline"
	KeyboardReply               = "reply"
)

type Keyboard struct {
	rowLen         int
	values         interface{}
	inlineKeyboard tgbotapi.InlineKeyboardMarkup
	typ            KeyboardType
	remove         bool

	enabled bool
}

func NewKeyboard(values interface{}) *Keyboard {
	enabled := false
	if values != nil {
		enabled = true
	}
	return &Keyboard{values: values, rowLen: DefaultKeyboardRowLen, typ: KeyboardReply, enabled: enabled}
}

func (k *Keyboard) SetRowLen(l int) {
	k.enabled = true
	k.rowLen = l
}

func (k *Keyboard) SetType(typ KeyboardType) {
	k.enabled = true
	k.typ = typ
}

func (k *Keyboard) Remove() {
	k.enabled = true
	k.remove = true
}

func (k *Keyboard) Reset() *Keyboard {
	k.values = nil
	k.enabled = false
	return k
}

func (k *Keyboard) AddCallbackButton(text string, datas ...string) *Keyboard {

	k.enabled = true

	data := text
	if len(datas) == 1 {
		data = datas[0]
	}
	k.typ = KeyboardInline
	if k.values == nil {
		k.values = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.InlineKeyboardButton{
						Text:         text,
						CallbackData: &data,
					},
				},
			},
		}
	}
	return k
}

func (k *Keyboard) AddReplyButton(text string) *Keyboard {

	k.enabled = true

	k.typ = KeyboardReply
	if k.values == nil {
		k.values = tgbotapi.ReplyKeyboardMarkup{
			Keyboard: [][]tgbotapi.KeyboardButton{
				[]tgbotapi.KeyboardButton{
					tgbotapi.KeyboardButton{
						Text: text,
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
			kb.Keyboard[rows-1] = append(kb.Keyboard[rows-1], tgbotapi.KeyboardButton{
				Text: text,
			})
			k.values = kb
		}
	}
	return k
}

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
			keyboard := [][]tgbotapi.KeyboardButton{[]tgbotapi.KeyboardButton{}}
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
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{[]tgbotapi.InlineKeyboardButton{}},
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
			Keyboard:       [][]tgbotapi.KeyboardButton{[]tgbotapi.KeyboardButton{tgbotapi.KeyboardButton{Text: v}}},
			ResizeKeyboard: true,
		}

	}

	return nil
}
