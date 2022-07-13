package tamework

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Update warap tgbotapi.Update and provide
// several helpers
type Update struct {
	tgbotapi.Update
}

type UpdateType int

const (
	UpdateTypeMessage UpdateType = iota + 1
	UpdateTypeEditedMessage
	UpdateTypeChannelPost
	UpdateTypeEditedChannelPost
	UpdateTypeInlineQuery
	UpdateTypeChosenInlineResult
	UpdateTypeCallbackQuery
	UpdateShippingQuery
	UpdatePreCheckoutQuery
	UpdatePoll
)

func (u Update) Type() UpdateType {
	if u.Update.Message != nil {
		return UpdateTypeMessage
	}
	if u.Update.EditedMessage != nil {
		return UpdateTypeMessage
	}
	if u.Update.ChannelPost != nil {
		return UpdateTypeChannelPost
	}
	if u.Update.EditedChannelPost != nil {
		return UpdateTypeEditedChannelPost
	}
	if u.Update.InlineQuery != nil {
		return UpdateTypeInlineQuery
	}
	if u.Update.ChosenInlineResult != nil {
		return UpdateTypeChosenInlineResult
	}
	if u.Update.CallbackQuery != nil {
		return UpdateTypeCallbackQuery
	}
	if u.Update.ShippingQuery != nil {
		return UpdateShippingQuery
	}
	if u.Update.PreCheckoutQuery != nil {
		return UpdatePreCheckoutQuery
	}
	//TOOD: add poll
	return 0
}

// NewUpdate returns Update by tgbotapi.Update
func NewUpdate(tupdate tgbotapi.Update) Update {
	update := Update{
		Update: tupdate,
	}

	return update
}

// ChatID detect update type and return ChatID
func (u Update) ChatID() int64 {
	if u.Message != nil {
		return u.Message.Chat.ID
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.Message.Chat.ID
	} else if u.ShippingQuery != nil {
		return int64(u.ShippingQuery.From.ID)
	} else if u.PreCheckoutQuery != nil {
		return int64(u.PreCheckoutQuery.From.ID)
	}
	return 0

}

// UserID detect update type and return UserID
func (u Update) UserID() int64 {
	if u.Message != nil {
		return int64(u.Message.From.ID)
	} else if u.CallbackQuery != nil {
		return int64(u.CallbackQuery.From.ID)
	} else if u.ShippingQuery != nil {
		return int64(u.ShippingQuery.From.ID)
	} else if u.PreCheckoutQuery != nil {
		return int64(u.PreCheckoutQuery.From.ID)
	} else if u.InlineQuery != nil {
		return int64(u.InlineQuery.From.ID)
	}
	return 0

}

// Method detects update type and return it
func (u Update) Method() string {
	if u.Message != nil {
		if u.Message.ReplyToMessage != nil {
			return Reply
		} else if u.Message.Photo != nil {
			return Photo
		}
		return Message
	} else if u.CallbackQuery != nil {
		return CallbackQuery
	} else if u.ShippingQuery != nil {

		return ShippingQuery
		//c.Text == ?
	} else if u.PreCheckoutQuery != nil {
		return PreCheckoutQuery
	} else if u.InlineQuery != nil {
		return InlineQuery
	}

	return ""
}

// Text returns text of message
func (u Update) Text() string {
	if u.Message != nil {
		return u.Message.Text
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.Data
	} else if u.ShippingQuery != nil {
		return ""
		//c.Text == ?
	} else if u.PreCheckoutQuery != nil {
		return ""
	} else if u.InlineQuery != nil {
		return u.InlineQuery.Query
	}
	return ""
}

// Username returns sender username
func (u Update) Username() (res string) {
	if u.Message != nil {
		res = u.Message.From.UserName
	} else if u.CallbackQuery != nil {
		res = u.CallbackQuery.From.UserName
	} else if u.ShippingQuery != nil {
		res = u.ShippingQuery.From.UserName
	} else if u.PreCheckoutQuery != nil {
		res = u.PreCheckoutQuery.From.UserName
	}

	if res == "" {
		return "_" + strconv.Itoa(int(u.UserID()))
	}

	return
}

// FirstName returns sender first name
func (u Update) FirstName() string {
	if u.Message != nil {
		return u.Message.Chat.FirstName
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.From.FirstName
	} else if u.ShippingQuery != nil {
		return u.ShippingQuery.From.FirstName
	} else if u.PreCheckoutQuery != nil {
		return u.PreCheckoutQuery.From.FirstName
	}
	return ""
}

// LastName returns sender last name
func (u Update) LastName() string {
	if u.Message != nil {
		return u.Message.Chat.LastName
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.From.LastName
	} else if u.ShippingQuery != nil {
		return u.ShippingQuery.From.LastName
	} else if u.PreCheckoutQuery != nil {
		return u.PreCheckoutQuery.From.LastName
	}
	return ""
}
