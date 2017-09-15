package tamework

import (
	"time"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Context struct {
	Method string

	UserID int64
	ChatID int64
	Text   string

	waiter *Waiter

	update   Update
	tamework *Tamework
	Keyboard *Keyboard
}

func NewContext(update tgbotapi.Update, tamework *Tamework) *Context {
	c := &Context{
		tamework: tamework,
		update:   NewUpdate(update),
		waiter:   tamework.waiter,
		Keyboard: NewKeyboard(nil),
	}

	c.ChatID = c.update.ChatID()
	c.Text = c.update.Text()
	c.Method = c.update.Method()

	c.Text, _ = tamework.Resolve(c.Text)

	return c
}

func (c *Context) Send(text string) error {
	return c.send(text, "")
}

func (c *Context) HTML(html string) error {
	return c.send(html, tgbotapi.ModeHTML)
}

func (c *Context) Markdown(md string) error {
	return c.send(md, tgbotapi.ModeMarkdown)
}

func (c *Context) send(text string, mode string) error {
	return c.sendTo(c.ChatID, text, mode)
}

func (c *Context) SendTo(id int64, text string) error {
	return c.sendTo(id, text, "")
}

func (c *Context) sendTo(id int64, text string, mode string) error {
	kbmsg := tgbotapi.NewMessage(id, text)
	if c.Keyboard != nil {
		kbmsg.ReplyMarkup = c.Keyboard.Markup()
	}
	kbmsg.ParseMode = mode
	_, err := c.tamework.bot.Send(kbmsg)
	if err != nil {
		return err
	}
	// reset keyboard
	c.NewKeyboard(nil)
	return nil
}

func (c *Context) SendInvoice() error {
	color.Green("New invoice for %d", c.ChatID)
	invoice := tgbotapi.NewInvoice(c.ChatID, "New invoice", "description here", "lalka",
		"361519591:TEST:68ca07b04a6cb4f8c7b68b78dbfd5c0a", "12345", "RUB",
		&[]tgbotapi.LabeledPrice{{Label: "RUB", Amount: 10000}})
	resp, err := c.tamework.bot.Send(invoice)
	color.Green("%s", resp)
	return err
}

func (c *Context) SendShippingAnswer() error {
	sc := tgbotapi.ShippingConfig{
		ShippingQueryID: c.update.ShippingQuery.ID,
		OK:              true,
		ShippingOptions: &[]tgbotapi.ShippingOption{{
			ID:     "allo1",
			Title:  "nogami",
			Prices: &[]tgbotapi.LabeledPrice{{Label: "r", Amount: 60}},
		}},
	}
	_, err := c.tamework.bot.AnswerShippingQuery(sc)
	return err
}

func (c *Context) SendPrecheckoutAnswer() error {
	pca := tgbotapi.PreCheckoutConfig{
		OK:                 true,
		PreCheckoutQueryID: c.update.PreCheckoutQuery.ID,
	}
	_, err := c.tamework.bot.AnswerPreCheckoutQuery(pca)
	return err
}

func (c *Context) Wait(stopword string, durations ...time.Duration) (Update, bool) {
	return c.waiter.Wait(c.ChatID, stopword, durations...)
}

func (c *Context) NewKeyboard(values interface{}) {
	c.Keyboard = NewKeyboard(values)
}

func (c *Context) BotAPI() *tgbotapi.BotAPI {
	return c.tamework.bot
}
