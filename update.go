package tamework

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Update struct {
	tgbotapi.Update
}

func NewUpdate(tupdate tgbotapi.Update) Update {
	update := Update{
		Update: tupdate,
	}

	return update
}

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

// if update.Message != nil {
//   c.UserID = int64(update.Message.From.ID)
//   c.ChatID = update.Message.Chat.ID
//   c.Method = Message
//   if update.Message.ReplyToMessage != nil {
//     c.Method = Reply
//   }
//   c.Text = update.Message.Text
// } else if update.CallbackQuery != nil {
//   c.UserID = int64(update.CallbackQuery.Message.From.ID)
//   c.ChatID = update.CallbackQuery.Message.Chat.ID
//   c.Method = CallbackQuery
//   c.Text = update.CallbackQuery.Data
// } else if update.ShippingQuery != nil {
//   c.UserID = int64(update.ShippingQuery.From.ID)
//   c.ChatID = c.UserID
//   c.Method = ShippingQuery
//   //c.Text == ?
// } else if update.PreCheckoutQuery != nil {
//   c.UserID = int64(update.PreCheckoutQuery.From.ID)
//   c.ChatID = c.UserID
//   c.Method = PreCheckoutQuery
// }

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

func (u Update) Username() string {
	if u.Message != nil {
		return u.Message.Chat.UserName
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.From.UserName
	} else if u.ShippingQuery != nil {
		return u.ShippingQuery.From.UserName
	} else if u.PreCheckoutQuery != nil {
		return u.PreCheckoutQuery.From.UserName
	}
	return ""
}

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
