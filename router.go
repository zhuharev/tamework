package tamework

import (
	"strings"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	tamework   *Tamework

	aliases map[string]string
}

// NewRouter returns new *Router instance
func NewRouter(tamework *Tamework) *Router {
	return &Router{
		tamework:   tamework,
		routeTable: make(map[string]map[string]Handler),
		aliases:    make(map[string]string),
	}
}

// HandleFunc type for handlers
type HandleFunc func(c *Context)

// CallbackQuery registre handler for CallbackQuery which have pattern as text
func (r *Router) CallbackQuery(pattern string, fn Handler) {
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
