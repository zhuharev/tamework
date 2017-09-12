package tamework

import (
	"fmt"
	"pure/api/socs/telegram/trinity/pkg/middleware"
)

var (
	Message       = "message"
	CallbackQuery = "callback_query"
	Reply         = "reply"
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

func (r *Router) Text(pattern string, fn HandleFunc) {
	r.registre(Message, pattern, fn)
}

func (r *Router) Reply(pattern string, fn HandleFunc) {
	r.registre(Reply, pattern, fn)
}

func (r *Router) registre(method string, pattern string, fn HandleFunc) {
	if r.routeTable[pattern] == nil {
		r.routeTable[pattern] = make(map[string]HandleFunc)
	}
	r.routeTable[pattern][method] = fn
	for _, fn := range r.tamework.locals {
		middleware.RegistreMethod(pattern, fn(pattern))
	}
}

func (r *Router) Handle(ctx *Context) {
	if m, has := r.routeTable[ctx.Text]; !has {
		ctx.Send(fmt.Sprintf("Команда %s (%s) не найдена", ctx.Text, ctx.Method))
		return
	} else {
		if handler, has := m[ctx.Method]; !has {
			ctx.Send(fmt.Sprintf("Команда %s (%s) не найдена", ctx.Text, ctx.Method))
		} else {
			handler(ctx)
		}
	}
}
