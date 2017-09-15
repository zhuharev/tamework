package tamework

import (
	"fmt"
	"pure/api/socs/telegram/trinity/pkg/middleware"
	"strings"
)

var (
	Message          = "message"
	CallbackQuery    = "callback_query"
	Reply            = "reply"
	ShippingQuery    = "shipping_query"
	PreCheckoutQuery = "pre_checkout_query"
	Photo            = "photo"

	Prefix = "prefix"
)

type Router struct {
	routeTable map[string]map[string]HandleFunc
	tamework   *Tamework
}

func NewRouter(tamework *Tamework) *Router {
	return &Router{
		tamework:   tamework,
		routeTable: make(map[string]map[string]HandleFunc),
	}
}

type HandleFunc func(c *Context)

func (r *Router) CallbackQuery(pattern string, fn HandleFunc) {
	r.registre(CallbackQuery, pattern, fn)
}

func (r *Router) ShippingQuery(pattern string, fn HandleFunc) {
	r.registre(ShippingQuery, pattern, fn)
}

func (r *Router) PreCheckoutQuery(pattern string, fn HandleFunc) {
	r.registre(PreCheckoutQuery, pattern, fn)
}

func (r *Router) Text(pattern string, fn HandleFunc) {
	r.registre(Message, pattern, fn)
}

func (r *Router) Reply(pattern string, fn HandleFunc) {
	r.registre(Reply, pattern, fn)
}

func (r *Router) Prefix(pattern string, fn HandleFunc) {
	r.registre(Prefix, pattern, fn)
}

func (r *Router) registre(method string, pattern string, fn HandleFunc) {
	if r.routeTable[method] == nil {
		r.routeTable[method] = make(map[string]HandleFunc)
	}
	r.routeTable[method][pattern] = fn
	for _, fn := range r.tamework.locals {
		middleware.RegistreMethod(pattern, fn(pattern))
	}
}

func (r *Router) Handle(ctx *Context) {

	for key, handler := range r.routeTable[Prefix] {
		if strings.HasPrefix(ctx.Text, key) {
			ctx.Text = strings.TrimPrefix(ctx.Text, key)
			handler(ctx)
			return
		}
	}

	if m, has := r.routeTable[ctx.Method]; !has {
		ctx.Send(fmt.Sprintf("Команда %s (%s) не найдена", ctx.Text, ctx.Method))
		return
	} else {
		if handler, has := m[ctx.Text]; !has {
			ctx.Send(fmt.Sprintf("Команда %s (%s) не найдена", ctx.Text, ctx.Method))
		} else {
			handler(ctx)
		}
	}
}
