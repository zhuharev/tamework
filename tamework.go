package tamework

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	DefaultWaitTimeout = time.Second * 60
)

// Tamework main instance
type Tamework struct {
	bot *tgbotapi.BotAPI
	*Router
	locals      map[string]i18n.TranslateFunc
	methods     map[string]string
	waiter      *Waiter
	WaitTimeout time.Duration
	AutoTyping  bool
}

// New return
func New(accessToken string) (_ *Tamework, err error) {
	bot, err := tgbotapi.NewBotAPI(accessToken)
	if err != nil {
		return
	}
	bot.Debug = true
	tw := &Tamework{
		bot:        bot,
		locals:     make(map[string]i18n.TranslateFunc),
		methods:    make(map[string]string),
		waiter:     NewWaiter(DefaultWaitTimeout),
		AutoTyping: true,
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
	} else {
		log.Printf("%v", update.InlineQuery)
	}

	ctx := NewContext(update, tw)
	if tw.AutoTyping {
		ca := tgbotapi.NewChatAction(ctx.ChatID, tgbotapi.ChatTyping)
		_, err := tw.bot.Send(ca)
		if err != nil {
			log.Println(err)
		}
	}

	if !tw.waiter.NeedNext(ctx.ChatID,
		ctx.update) {
		return
	}
	color.Cyan("%s (%s) %d", ctx.Method, ctx.Text, ctx.ChatID)

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
