package tamework

import (
	"strings"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	Message          = "message"
	CallbackQuery    = "callback_query"
	Reply            = "reply"
	ShippingQuery    = "shipping_query"
	PreCheckoutQuery = "pre_checkout_query"
	Photo            = "photo"
	InlineQuery      = "inline_query"

	Prefix = "prefix"
)

type Router struct {
	routeTable map[string]map[string]Handler
	tamework   *Tamework

	aliases map[string]string
}

func NewRouter(tamework *Tamework) *Router {
	return &Router{
		tamework:   tamework,
		routeTable: make(map[string]map[string]Handler),
		aliases:    make(map[string]string),
	}
}

type HandleFunc func(c *Context)

func (r *Router) CallbackQuery(pattern string, fn Handler) {
	r.registre(CallbackQuery, pattern, fn)
}

func (r *Router) InlineQuery(fn Handler) {
	r.registre(InlineQuery, "", fn)
}

func (r *Router) ShippingQuery(pattern string, fn Handler) {
	r.registre(ShippingQuery, pattern, fn)
}

func (r *Router) PreCheckoutQuery(pattern string, fn Handler) {
	r.registre(PreCheckoutQuery, pattern, fn)
}

func (r *Router) Text(pattern string, fn Handler) {
	r.registre(Message, pattern, fn)
}

func (r *Router) Reply(pattern string, fn Handler) {
	r.registre(Reply, pattern, fn)
}

func (r *Router) Prefix(pattern string, fn Handler) {
	r.registre(Prefix, pattern, fn)
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
		if m, has := r.routeTable[ctx.Method]; !has {
			if r.tamework.NotFound != nil {
				currHandler = r.tamework.NotFound
			}
			//ctx.Send(fmt.Sprintf("Команда %s (%s) не найдена", ctx.Text, ctx.Method))
			//log.Println(ctx.update.Update)
			//	return
		} else {
			var txt = ctx.Text
			if ctx.Method == InlineQuery {
				txt = ""
			}
			if handler, has := m[txt]; !has {
				//ctx.Send(fmt.Sprintf("Команда %s (%s) не найдена", ctx.Text, ctx.Method))
				//log.Println(ctx.update.Update)
				if r.tamework.NotFound != nil {
					currHandler = r.tamework.NotFound
				}
			} else {
				currHandler = handler
			}
		}
	}
	if currHandler != nil {
		ctx.handlers = append(ctx.tamework.handlers, currHandler)
	}
	ctx.run()
}
