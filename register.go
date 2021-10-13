package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Storage interface {
	Store(Pkg)
	Values() []Pkg
}

// Register adds handlers to provided mux.Router.
func Register(r *mux.Router, s Storage) error {
	const fName = "Register"

	t, err := template.ParseFiles("index.html")
	if err != nil {
		return errors.Wrap(err, fName)
	}

	r.
		Methods("GET").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Trace().Str("method", "GET").Msg("Serve template")
			if err := t.Execute(w, s.Values()); err != nil {
				err = errors.Wrap(err, fName)
				log.Error().Stack().Err(err).Str("method", "GET").Send()
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		})

	r.
		Methods("POST").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			s.Store(p)

			http.Redirect(w, r, r.URL.Path, 302)
		})

	return nil
}
