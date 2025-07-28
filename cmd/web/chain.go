package main

import "net/http"

type Constructor func(http.Handler) http.Handler

type Chain struct {
	constructor []Constructor
}

func New(c ...Constructor) Chain {
	return Chain{append(([]Constructor)(nil), c...)}
}

func (c Chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}
	for i := range c.constructor {
		h = c.constructor[len(c.constructor)-i-1](h)
	}
	return h
}

func (c Chain) ThenFunc(h http.HandlerFunc) http.Handler {
	if h == nil {
		return c.Then(nil)
	}
	return c.Then(h)
}

func (c Chain) Append(con ...Constructor) Chain {
	newCons := make([]Constructor, 0, len(c.constructor)+len(con))
	newCons = append(newCons, c.constructor...)
	newCons = append(newCons, con...)
	return Chain{newCons}
}

func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.constructor...)
}
