package tamework

import (
	"context"
	"time"

	"github.com/go-macaron/inject"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Context will be created for all requests
type Context struct {
	context.Context
	inject.Injector

	index  int
	router map[string]map[string]Handler

	handlers []Handler
	action   Handler

	Method string

	UserID int64
	ChatID int64
	Text   string
	State  State

	waiter *Waiter

	update   Update
	tamework *Tamework
	Keyboard *Keyboard

	Data map[string]interface{}

	// if before middleware exit
	exited bool

	T func(translationID string, args ...interface{}) string
}

// Update reuturns update
func (c Context) Update() Update {
	return c.update
}

// Exit shows that the request is processed and next handlers will not be called
func (c *Context) Exit() {
	c.exited = true
}

// createContextjust create context
func (tw *Tamework) createContext(update tgbotapi.Update) *Context {
	c := &Context{
		Context:  context.Background(),
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

// Send send text
func (c *Context) Send(text string) (msg tgbotapi.Message, err error) {
	return c.send(text, "")
}

// HTML send html
func (c *Context) HTML(html string) (msg tgbotapi.Message, err error) {
	return c.send(html, tgbotapi.ModeHTML)
}

// Markdown send markdown
func (c *Context) Markdown(md string) (msg tgbotapi.Message, err error) {
	return c.send(md, tgbotapi.ModeMarkdown)
}

// send is just helper
func (c *Context) send(text string, mode string) (msg tgbotapi.Message, err error) {
	return c.sendTo(c.ChatID, text, mode)
}

// SendTo send message to specific chat
func (c *Context) SendTo(id int64, text string) (msg tgbotapi.Message, err error) {
	return c.sendTo(id, text, "")
}

// MarkdownTo send message to specific chat
func (c *Context) MarkdownTo(id int64, text string) (msg tgbotapi.Message, err error) {
	return c.sendTo(id, text, tgbotapi.ModeMarkdown)
}

// sendTo just helper
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

// TODO: implement this
// SendInvoice send new invoice
func (c *Context) SendInvoice() error {
	return nil
}

// TODO:implement this
// SendShippingAnswer send shipping query answer
func (c *Context) SendShippingAnswer() (err error) {
	// sc := tgbotapi.ShippingConfig{
	// 	ShippingQueryID: c.update.ShippingQuery.ID,
	// 	OK:              true,
	// 	ShippingOptions: &[]tgbotapi.ShippingOption{{
	// 		ID:     "allo1",
	// 		Title:  "nogami",
	// 		Prices: &[]tgbotapi.LabeledPrice{{Label: "r", Amount: 60}},
	// 	}},
	// }
	// _, err := c.tamework.bot.AnswerShippingQuery(sc)
	// return err
	return
}

// AnswerTo answer for callback query
func (c *Context) AnswerTo(to string, text string, alerts ...bool) error {
	alert := false
	if len(alerts) > 0 {
		alert = alerts[0]
	}

	cfg := tgbotapi.NewCallback(to, text)
	cfg.ShowAlert = alert
	_, err := c.tamework.bot.Send(cfg)
	return err
}

// Answer answer for context.Update()
func (c *Context) Answer(text string, alerts ...bool) error {
	return c.AnswerTo(c.update.CallbackQuery.ID, text, alerts...)
}

// AnswerInline answer for inline query
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
	_, err := c.tamework.bot.Send(cfg)
	return err
}

// SendPrecheckoutAnswer send PreCheckoutQueryanswer
// it always sends true
func (c *Context) SendPrecheckoutAnswer() error {
	pca := tgbotapi.PreCheckoutConfig{
		OK:                 true,
		PreCheckoutQueryID: c.update.PreCheckoutQuery.ID,
	}
	_, err := c.tamework.bot.Send(pca)
	return err
}

// EditText edit text for message from context.Update
func (c *Context) EditText(newMessage string) error {
	return c.newEdit(newMessage, "")
}

// EditMarkdown edit text for message from context.Update
func (c *Context) EditMarkdown(newMessage string) error {
	return c.newEdit(newMessage, tgbotapi.ModeMarkdown)
}

// newEdit just helper
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

// EditCaption for file
func (c *Context) EditCaption(newCaption string) error {
	cnf := tgbotapi.NewEditMessageCaption(c.ChatID, c.update.CallbackQuery.Message.MessageID, newCaption)
	_, err := c.tamework.bot.Send(cnf)
	return err
}

// EditReplyMarkup edit keyboard
func (c *Context) EditReplyMarkup(kb *Keyboard) error {
	var iface interface{}
	if kb != nil {
		iface = kb.Markup()
	}
	if rkb, ok := iface.(tgbotapi.InlineKeyboardMarkup); ok {
		cnf := tgbotapi.NewEditMessageReplyMarkup(c.ChatID,
			c.update.CallbackQuery.Message.MessageID, rkb)
		_, err := c.tamework.bot.Send(cnf)
		return err
	}
	cnf := tgbotapi.NewEditMessageReplyMarkup(c.ChatID,
		c.update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{})
	_, err := c.tamework.bot.Send(cnf)
	return err
}

// Wait input from user
// if messge text will be contains stopword, waiter will be canceled
func (c *Context) Wait(stopword string, durations ...time.Duration) (Update, bool) {
	return c.waiter.Wait(c.ChatID, stopword, durations...)
}

// NewKeyboard helper for creating keyboard
func (c *Context) NewKeyboard(values interface{}) *Keyboard {
	c.Keyboard = NewKeyboard(values)
	return c.Keyboard
}

// BotAPI returns tgbotapi.BotAPI instance
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

// Next call next handler
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
