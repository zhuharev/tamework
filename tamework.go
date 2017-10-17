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

	inject.Injector
	handlers []Handler
	action   Handler

	NotFound Handler

	RejectOldUpdates int //seconds

	Locales []func(translationID string, args ...interface{}) string
}

func (tw *Tamework) Use(handler Handler) {
	tw.handlers = append(tw.handlers, handler)
}

// New return
func New(accessToken string) (_ *Tamework, err error) {
	bot, err := tgbotapi.NewBotAPI(accessToken)
	if err != nil {
		return
	}
	tw := &Tamework{
		bot:        bot,
		methods:    make(map[string]string),
		waiter:     NewWaiter(DefaultWaitTimeout),
		AutoTyping: true,
		action:     func() {},
	}
	tw.Router = NewRouter(tw)
	return tw, nil
}

func (tw *Tamework) Bot() *tgbotapi.BotAPI {
	return tw.bot
}

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

	//ctx := tw.createContext(update, tw)
	if tw.AutoTyping {
		ca := tgbotapi.NewChatAction(up.ChatID(), tgbotapi.ChatTyping)
		_, err := tw.bot.Send(ca)
		if err != nil {
			log.Println(err)
		}
	}

	// if !tw.waiter.NeedNext(ctx.ChatID,
	// 	ctx.update) {
	// 	return
	// }
	//color.Cyan("%s (%s) %d", ctx.Method, ctx.Text, ctx.ChatID)

	tw.Handle(update)
}

func (tw *Tamework) RegistreMethod(method string, buttonCaption string) {
	tw.methods[buttonCaption] = method
}

func (tw *Tamework) Resolve(text string) (method string, has bool) {
	method, has = tw.methods[text]
	if !has {
		method = text
	}
	return
}
