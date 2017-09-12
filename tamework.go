package tamework

import (
	"time"

	"github.com/nicksnyder/go-i18n/i18n"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var (
	DefaultWaitTimeout = time.Second * 5
)

// Tamework main instance
type Tamework struct {
	bot *tgbotapi.BotAPI
	*Router
	locals      map[string]i18n.TranslateFunc
	methods     map[string]string
	waiter      *Waiter
	WaitTimeout time.Duration
}

// New return
func New(accessToken string) (_ *Tamework, err error) {
	bot, err := tgbotapi.NewBotAPI(accessToken)
	if err != nil {
		return
	}
	tw := &Tamework{
		bot:     bot,
		locals:  make(map[string]i18n.TranslateFunc),
		methods: make(map[string]string),
		waiter:  NewWaiter(DefaultWaitTimeout),
	}
	tw.Router = NewRouter(tw)
	return tw, nil
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

func (tw *Tamework) handleUpdate(update tgbotapi.Update) {
	ctx := NewContext(update, tw)

	if !tw.waiter.NeedNext(ctx.ChatID,
		ctx.Text) {
		return
	}
	tw.Handle(ctx)
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
