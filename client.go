package main

import (
	"html/template"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Storage interface {
	Store(Pkg)
	Values() []Pkg
}

type Client struct {
	s Storage
	t *template.Template
}

var _ http.Handler = (*Client)(nil)

func NewClient(s Storage) (*Client, error) {
	const fName = "NewClient"
	t, err := template.ParseFiles("index.html")
	if err != nil {
		return nil, errors.Wrap(err, fName)
	}
	return &Client{s: s, t: t}, nil
}

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.handleGET(w, r)
	case "POST":
		c.handlePOST(w, r)
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

func (c *Client) handleGET(w http.ResponseWriter, r *http.Request) {
	const fName = "Client.handleGET"
	log.Trace().Str("method", "GET").Msg("Serve template")
	if err := c.t.Execute(w, c.s.Values()); err != nil {
		err = errors.Wrap(err, fName)
		log.Error().Stack().Err(err).Str("method", "GET").Send()
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
func (c *Client) handlePOST(w http.ResponseWriter, r *http.Request) {
	const fName = "Client.handlePOST"
	log.Trace().Str("method", "POST").Msg("Receive form")
	if err := r.ParseForm(); err != nil {
		err = errors.Wrap(err, fName)
		log.Error().Stack().Err(err).Str("method", "POST").Send()
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	status := Ordered
	if r.Form.Has("shipped") {
		status = Shipped
	}
	p := NewPkg(
		r.Form.Get("name"),
		WithInpost(r.Form.Has("inpost")),
		WithStatus(status),
	)
	log.Trace().Interface("Pkg", p).Msg("Store")
	c.s.Store(p)

	http.Redirect(w, r, r.URL.Path, 302)
}
