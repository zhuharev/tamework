package tamework

import (
	"time"

	"github.com/fatih/color"
	"github.com/go-macaron/inject"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Context struct {
	inject.Injector

	index  int
	router map[string]map[string]Handler

	handlers []Handler
	action   Handler

	Method string

	UserID int64
	ChatID int64
	Text   string

	waiter *Waiter

	update   Update
	tamework *Tamework
	Keyboard *Keyboard

	Data map[string]interface{}

	// if before middleware exit
	exited bool

	T func(translationID string, args ...interface{}) string
}

func (c Context) Update() Update {
	return c.update
}

func (c *Context) Exit() {
	c.exited = true
}

func (tw *Tamework) createContext(update tgbotapi.Update) *Context {
	c := &Context{
		Injector: inject.New(),
		tamework: tw,
		index:    0,
		//	router:   tw.routeTable,
		action:   tw.action,
		update:   NewUpdate(update),
		waiter:   tw.waiter,
		Keyboard: NewKeyboard(nil),
		Data:     map[string]interface{}{},
	}
	c.SetParent(tw)
	c.Map(c)

	c.ChatID = c.update.ChatID()
	c.UserID = c.update.UserID()
	c.Text = c.update.Text()
	c.Method = c.update.Method()

	c.Text, _ = tw.Resolve(c.Text)

	return c
}

func (c *Context) Send(text string) (msg tgbotapi.Message, err error) {
	return c.send(text, "")
}

func (c *Context) HTML(html string) (msg tgbotapi.Message, err error) {
	return c.send(html, tgbotapi.ModeHTML)
}

func (c *Context) Markdown(md string) (msg tgbotapi.Message, err error) {
	return c.send(md, tgbotapi.ModeMarkdown)
}

func (c *Context) send(text string, mode string) (msg tgbotapi.Message, err error) {
	return c.sendTo(c.ChatID, text, mode)
}

func (c *Context) SendTo(id int64, text string) (msg tgbotapi.Message, err error) {
	return c.sendTo(id, text, "")
}

func (c *Context) sendTo(id int64, text string, mode string) (msg tgbotapi.Message, err error) {
	kbmsg := tgbotapi.NewMessage(id, text)
	if c.Keyboard != nil {
		kbmsg.ReplyMarkup = c.Keyboard.Markup()
	}
	kbmsg.ParseMode = mode
	msg, err = c.tamework.bot.Send(kbmsg)
	if err != nil {
		return
	}
	// reset keyboard
	c.NewKeyboard(nil)
	return
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

func (c *Context) AnswerTo(to string, text string, alerts ...bool) error {
	alert := false
	if len(alerts) > 0 {
		alert = alerts[0]
	}

	cfg := tgbotapi.NewCallback(to, text)
	cfg.ShowAlert = alert
	_, err := c.tamework.bot.AnswerCallbackQuery(cfg)
	return err
}

func (c *Context) Answer(text string, alerts ...bool) error {
	return c.AnswerTo(c.update.CallbackQuery.ID, text, alerts...)
}

func (c *Context) AnswerInline() error {
	if c.update.Update.InlineQuery == nil {
		return nil
	}
	queryID := c.update.Update.InlineQuery.ID
	art := tgbotapi.NewInlineQueryResultArticle("da", "hello", "text")
	cfg := tgbotapi.InlineConfig{
		InlineQueryID: queryID,
		Results:       []interface{}{art},
	}
	_, err := c.tamework.bot.AnswerInlineQuery(cfg)
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

func (c *Context) EditText(newMessage string) error {
	return c.newEdit(newMessage, "")
}

func (c *Context) EditMarkdown(newMessage string) error {
	return c.newEdit(newMessage, tgbotapi.ModeMarkdown)
}

func (c *Context) newEdit(message string, parseMode string) error {
	cnf := tgbotapi.NewEditMessageText(c.ChatID, c.update.CallbackQuery.Message.MessageID, message)
	if parseMode != "" {
		cnf.ParseMode = parseMode
	}
	if kb, ok := c.Keyboard.Markup().(tgbotapi.InlineKeyboardMarkup); ok {
		cnf.ReplyMarkup = &kb
	}
	_, err := c.tamework.bot.Send(cnf)
	return err
}

func (c *Context) EditCaption(newCaption string) error {
	cnf := tgbotapi.NewEditMessageCaption(c.ChatID, c.update.CallbackQuery.Message.MessageID, newCaption)
	_, err := c.tamework.bot.Send(cnf)
	return err
}

func (c *Context) EditReplyMurkup(kb *Keyboard) error {
	var iface interface{}
	if kb != nil {
		iface = kb.Markup()
	}
	if rkb, ok := iface.(tgbotapi.InlineKeyboardMarkup); ok {
		cnf := tgbotapi.NewEditMessageReplyMarkup(c.ChatID,
			c.update.CallbackQuery.Message.MessageID, rkb)
		_, err := c.tamework.bot.Send(cnf)
		return err
	} else {
		cnf := tgbotapi.NewEditMessageReplyMarkup(c.ChatID,
			c.update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{})
		_, err := c.tamework.bot.Send(cnf)
		return err
	}
	return nil
}

func (c *Context) Wait(stopword string, durations ...time.Duration) (Update, bool) {
	return c.waiter.Wait(c.ChatID, stopword, durations...)
}

func (c *Context) NewKeyboard(values interface{}) *Keyboard {
	c.Keyboard = NewKeyboard(values)
	return c.Keyboard
}

func (c *Context) BotAPI() *tgbotapi.BotAPI {
	return c.tamework.bot
}

func (c *Context) handler() Handler {
	if c.index < len(c.handlers) {
		return c.handlers[c.index]
	}
	if c.index == len(c.handlers) {
		return c.action
	}
	panic("invalid index for context handler")
}

func (c *Context) Next() {
	c.index++
	c.run()
}

func (c *Context) run() {
	for c.index <= len(c.handlers) {
		if c.exited {
			return
		}
		_, err := c.Invoke(c.handler())
		if err != nil {
			panic(err)
		}
		c.index++
	}
}
