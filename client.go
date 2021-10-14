package main

import (
	"html/template"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Client is responsible for handling HTTP requests. It fulfills http.Handler
// interface.
type Client struct {
	s Storage
	t *template.Template
}

var _ http.Handler = (*Client)(nil)

// NewClient creates a Client object associated with provided Storage.
func NewClient(s Storage) (*Client, error) {
	const fName = "NewClient"
	t, err := template.ParseFiles("index.html")
	if err != nil {
		return nil, errors.Wrap(err, fName)
	}
	return &Client{s: s, t: t}, nil
}

// ServeHTTP handles the request with a dedicated handler function based on
// request's method.
func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = c.handleGET(w, r)
	case "POST":
		err = c.handlePOST(w, r)
	default:
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error().Stack().Err(errors.WithStack(err)).Send()
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (c *Client) handleGET(w http.ResponseWriter, r *http.Request) error {
	pkgs, err := c.s.LoadPkgs()
	if err != nil {
		return errors.WithStack(err)
	}
	if err := c.t.Execute(w, pkgs); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Client) handlePOST(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return errors.WithStack(err)
	}
	log.Debug().Interface("form", r.Form).Msg("Parsed form")

	status := Ordered
	if r.Form.Has("shipped") {
		status = Shipped
	}
	p := NewPkg(
		r.Form.Get("name"),
		WithInpost(r.Form.Has("inpost")),
		WithStatus(status),
	)
	if err := c.s.StorePkg(p); err != nil {
		return errors.WithStack(err)
	}
	log.Debug().Interface("Pkg", p).Msg("Stored package")

	http.Redirect(w, r, r.URL.Path, 302)

	return nil
}
