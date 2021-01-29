package tamework

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-macaron/inject"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	// DefaultWaitTimeout how long we wait for user input
	DefaultWaitTimeout = time.Second * 60
)

// Tamework main instance
type Tamework struct {
	bot *tgbotapi.BotAPI
	*Router
	methods     map[string]string
	waiter      *Waiter
	WaitTimeout time.Duration
	AutoTyping  bool
	State       StateStorage
	FormStore   FormStore

	inject.Injector
	handlers []Handler
	action   Handler

	NotFound Handler

	RejectOldUpdates int //seconds

	Locales []func(translationID string, args ...interface{}) string
}

type initOpFunc func(t *Tamework) error

// New returns Tamework instance
func New(accessToken string, funcs ...initOpFunc) (_ *Tamework, err error) {
	bot, err := tgbotapi.NewBotAPIWithClient(accessToken, &http.Client{Timeout: 10 * time.Second})
	if err != nil {
		return
	}

	tw := &Tamework{
		bot:        bot,
		methods:    make(map[string]string),
		waiter:     NewWaiter(DefaultWaitTimeout),
		AutoTyping: true,
		action:     func() {},
		State:      newMemStateStorage(),
		FormStore:  newMemFormStore(),
	}
	tw.Router = NewRouter(tw)

	for _, fn := range funcs {
		err = fn(tw)
		if err != nil {
			return nil, err
		}
	}

	return tw, nil
}

// WithStateStorage replace default memory state storage. can be used for persistent chat state.
func WithStateStorage(stateStorage StateStorage) func(*Tamework) error {
	return func(t *Tamework) error {
		t.State = stateStorage
		return nil
	}
}

// WithFormStorage replace default memory state storage. can be used for persistent chat state.
func WithFormStorage(formStorage FormStore) func(*Tamework) error {
	return func(t *Tamework) error {
		t.FormStore = formStorage
		return nil
	}
}

// MsgOp args func for flexible change message structure
type MsgOp func(m *tgbotapi.MessageConfig)

// ToChat set target receipent for message
func ToChat(id int) MsgOp {
	return func(m *tgbotapi.MessageConfig) {
		m.BaseChat.ChatID = int64(id)
	}
}

// WithKeyboard add keyboard markup to message
func WithKeyboard(kb *Keyboard) MsgOp {
	return func(m *tgbotapi.MessageConfig) {
		m.ReplyMarkup = kb.Markup()
	}
}

// Markdown tell telegram servers parse message text as markdown
func Markdown() MsgOp {
	return func(m *tgbotapi.MessageConfig) {
		m.ParseMode = tgbotapi.ModeMarkdown
	}
}

// Send sends message to id
func (tw *Tamework) Send(text string, fns ...MsgOp) (int, error) {
	kbmsg := tgbotapi.NewMessage(0, text)
	for _, fn := range fns {
		fn(&kbmsg)
	}
	resp, err := tw.bot.Send(kbmsg)
	if err != nil {
		return 0, err
	}
	return resp.MessageID, nil
}

// Use registre middleware
// This func will be used in each request
func (tw *Tamework) Use(handler Handler) {
	tw.handlers = append(tw.handlers, handler)
}

// Bot returns *tgbotapi.BotAPI
func (tw *Tamework) Bot() *tgbotapi.BotAPI {
	return tw.bot
}

// Run starts the pooler and send new updates
// to router
func (tw *Tamework) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := tw.bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		go tw.handleUpdate(update)
	}
}

// HandleUpdateWebhook implements http.Handler
func (tw *Tamework) HandleUpdateWebhook(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		defer req.Body.Close()
		var update tgbotapi.Update
		err := json.NewDecoder(req.Body).Decode(&update)
		if err != nil {
			log.Println(err)
			return
		}
		go tw.handleUpdate(update)
	}
}

func (tw *Tamework) handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {

		log.Println(update.Message)
	} else if update.InlineQuery != nil {
		log.Println(update.InlineQuery)
	}

	up := NewUpdate(update)
	if !tw.waiter.NeedNext(up.ChatID(), up) {
		return
	}

	if tw.AutoTyping {
		ca := tgbotapi.NewChatAction(up.ChatID(), tgbotapi.ChatTyping)
		_, err := tw.bot.Send(ca)
		if err != nil {
			log.Println(err)
		}
	}

	tw.Handle(update)
}

// RegistreMethod registre an alias for method
func (tw *Tamework) RegistreMethod(method string, buttonCaption string) {
	tw.methods[buttonCaption] = method
}

// Resolve resolve method name for passed alias
func (tw *Tamework) Resolve(text string) (method string, has bool) {
	method, has = tw.methods[text]
	if !has {
		method = text
	}
	return
}
