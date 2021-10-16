package packages

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

//go:embed index.html static
var _fs embed.FS

// Client is responsible for handling HTTP requests. It fulfills http.Handler
// interface.
type Client struct {
	r *mux.Router
	s Storage
	t *template.Template
}

var _ http.Handler = (*Client)(nil)

// NewClient creates a Client object associated with provided Storage.
func NewClient(s Storage) (*Client, error) {
	const fName = "NewClient"
	t, err := template.ParseFS(_fs, "index.html")
	if err != nil {
		return nil, errors.Wrap(err, fName)
	}
	c := &Client{r: mux.NewRouter(), s: s, t: t}
	// TODO: prevent directory listing
	c.r.Methods("GET").PathPrefix("/static").
		Handler(http.FileServer(http.FS(_fs)))
	c.r.Methods("GET").
		HandlerFunc(c.handleError(c.handleGET))
	c.r.Methods("POST").
		HandlerFunc(c.handleError(c.handlePOST))
	return c, nil
}

// ServeHTTP calls ServeHTTP on underlying *mux.Router.
func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.r.ServeHTTP(w, r)
}

func (c *Client) handleError(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Error().Stack().Err(errors.WithStack(err)).Send()
			http.Error(w, "", http.StatusInternalServerError)
		}
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
