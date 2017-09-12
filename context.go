package tamework

import tgbotapi "gopkg.in/telegram-bot-api.v4"

type Context struct {
	Method string

	UserID int64
	ChatID int64
	Text   string

	waiter *Waiter

	update   tgbotapi.Update
	tamework *Tamework
}

func NewContext(update tgbotapi.Update, tamework *Tamework) *Context {
	c := &Context{
		tamework: tamework,
		update:   update,
		waiter:   tamework.waiter,
	}
	if update.Message != nil {
		c.UserID = int64(update.Message.From.ID)
		c.ChatID = update.Message.Chat.ID
		c.Method = "message"
		if update.Message.ReplyToMessage != nil {
			c.Method = "reply"
		}
		c.Text = update.Message.Text
	} else if update.CallbackQuery != nil {
		c.UserID = int64(update.CallbackQuery.Message.From.ID)
		c.ChatID = update.CallbackQuery.Message.Chat.ID
		c.Method = "callback_query"
		c.Text = update.CallbackQuery.Data
	}

	c.Text, _ = tamework.Resolve(c.Text)

	return c
}

func (c *Context) Send(text string) error {
	kbmsg := tgbotapi.NewMessage(c.ChatID, text)
	_, err := c.tamework.bot.Send(kbmsg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) Wait() (string, bool) {
	return c.waiter.Wait(c.ChatID)
}
