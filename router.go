package tamework

import (
	"context"
	"strings"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// Message normal message update
	Message = "message"
	// CallbackQuery cbquery update
	CallbackQuery = "callback_query"
	// Reply reply to message update
	Reply = "reply"
	// ShippingQuery update
	ShippingQuery = "shipping_query"
	// PreCheckoutQuery update
	PreCheckoutQuery = "pre_checkout_query"
	// Photo update
	Photo = "photo"
	// InlineQuery update
	InlineQuery = "inline_query"

	// Prefix update
	Prefix = "prefix"
)

// Router route messages by text contents
type Router struct {
	routeTable map[string]map[string]Handler
	stateTable map[State]Handler
	formTable  map[string]FormHandler
	tamework   *Tamework

	aliases map[string]string
}

// NewRouter returns new *Router instance
func NewRouter(tamework *Tamework) *Router {
	return &Router{
		tamework:   tamework,
		routeTable: make(map[string]map[string]Handler),
		stateTable: make(map[State]Handler),
		formTable:  make(map[string]FormHandler),
		aliases:    make(map[string]string),
	}
}

// HandleFunc type for handlers
type HandleFunc func(c *Context)

// CallbackQuery registre handler for CallbackQuery which have pattern as text
func (r *Router) CallbackQuery(pattern string, fn Handler) {
	r.registre(CallbackQuery, pattern, fn)
}

// Cb is alias for CallbackQuery
func (r *Router) Cb(pattern string, fn Handler) {
	r.registre(CallbackQuery, pattern, fn)
}

// InlineQuery registre handler for InlineQuery which have pattern as text
func (r *Router) InlineQuery(fn Handler) {
	r.registre(InlineQuery, "", fn)
}

// ShippingQuery registre handler for ShippingQuery which have pattern as text
func (r *Router) ShippingQuery(pattern string, fn Handler) {
	r.registre(ShippingQuery, pattern, fn)
}

// PreCheckoutQuery registre handler for PreCheckoutQuery which have pattern as text
func (r *Router) PreCheckoutQuery(pattern string, fn Handler) {
	r.registre(PreCheckoutQuery, pattern, fn)
}

// Text registre handler for Text which have pattern as text
func (r *Router) Text(pattern string, fn Handler) {
	r.registre(Message, pattern, fn)
}

// Reply registre handler for Reply which have pattern as text
func (r *Router) Reply(pattern string, fn Handler) {
	r.registre(Reply, pattern, fn)
}

// Prefix registre handler in router. Router will be check all messages - if
// any message will contain a prefix, prefix will be deleted from message text
// and handler will process this update
func (r *Router) Prefix(pattern string, handler Handler) {
	r.registre(Prefix, pattern, handler)
}

// State registre handler in router. Router use tamework.StateStorage for define current
// user status. If current user status == state, handler will be called.
//
// State handlers has low priority compared to Reply or Text handlers.
func (r *Router) State(state State, handler Handler) {
	r.stateTable[state] = handler
}

func (r *Router) Form(text string, form *Form, handler FormHandler) {
	r.formTable[text] = handler
	r.registre(Message, text, func(ctx *Context) {
		form := form.Copy()
		question := form.GetNextQuestion()
		if question == nil {
			return
		}
		err := r.tamework.FormStore.SaveForm(ctx.Context, int(ctx.ChatID), int(ctx.UserID), form)
		if err != nil {
			ctx.Send("db error")
			return
		}
		// build buttons
		kb := ctx.NewKeyboard(nil)
		for _, v := range question.Answers {
			kb.AddCallbackButton(v)
		}
		_, err = ctx.Markdown(question.Text)
		if err != nil {
			ctx.Send("network error")
			return
		}
	})
}

func (r *Router) registre(method string, pattern string, fn Handler) {
	if r.routeTable[method] == nil {
		r.routeTable[method] = make(map[string]Handler)
	}
	r.routeTable[method][pattern] = fn
	for _, fn := range r.tamework.Locales {
		if fn(pattern) != pattern {
			r.aliases[fn(pattern)] = pattern
			color.Cyan("Registre alias: %s", fn(pattern))
		}
	}
}

// Handle is main router func which handle all updates
func (r *Router) Handle(update tgbotapi.Update) {
	var (
		ctx         = r.tamework.createContext(update)
		currHandler Handler
	)

	if aliase, has := r.aliases[ctx.Text]; has {
		ctx.Text = aliase
	}

	for key, handler := range r.routeTable[Prefix] {
		if strings.HasPrefix(ctx.Text, key) {
			ctx.Text = strings.TrimSpace(strings.TrimPrefix(ctx.Text, key))
			currHandler = handler
		}
	}

	if currHandler == nil {
		if m, has := r.routeTable[ctx.Method]; has {
			var txt = ctx.Text
			if ctx.Method == InlineQuery {
				txt = ""
			}
			if handler, has := m[txt]; has {
				currHandler = handler
			}
		}
	}

	if state, err := r.tamework.State.GetState(context.Background(), int(ctx.ChatID), int(ctx.UserID)); err == nil && state != "" {
		currHandler = r.stateTable[state]
	} else if err != nil {
		// TODO: handle storage error here
	}

	if form, err := r.tamework.FormStore.GetActiveForm(context.Background(), int(ctx.ChatID), int(ctx.UserID)); err == nil && form != nil {
		currHandler = form.MakeHandler(r.tamework.FormStore, r.formTable[form.Keyword])
	}

	if currHandler == nil && r.tamework.NotFound != nil {
		currHandler = r.tamework.NotFound
	}

	if currHandler != nil {
		ctx.handlers = append(ctx.tamework.handlers, currHandler)
	}
	ctx.run()
}
