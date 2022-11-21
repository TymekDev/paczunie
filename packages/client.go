package packages

import (
	"database/sql"
	"embed"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
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
	t, err := template.ParseFS(_fs, "index.html")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c := &Client{r: mux.NewRouter(), s: s, t: t}
	// TODO: prevent directory listing
	c.r.Methods("GET").PathPrefix("/static").
		Handler(http.FileServer(http.FS(_fs)))
	c.r.Methods("GET").
		HandlerFunc(c.handleError(c.handleGET))
	c.r.Methods("POST").
		HandlerFunc(c.handleError(c.handlePOST))
	c.r.Methods("PATCH").
		HandlerFunc(c.handleError(c.handlePATCH))
	c.r.Methods("DELETE").
		HandlerFunc(c.handleError(c.handleDELETE))
	return c, nil
}

// NewClientWithSQLiteStorage is a wrapper on opening connection to SQLite3
// database, creating a Storage with it, and creating Client with the Storage.
func NewClientWithSQLiteStorage(dbName string) (*Client, error) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dbs, err := NewDBStorage(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewClient(dbs)
}

// ServeHTTP calls ServeHTTP on underlying *mux.Router.
func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.r.ServeHTTP(w, r)
}

func (c *Client) handleError(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Error().Stack().Err(err).Send()
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func (c *Client) handleGET(w http.ResponseWriter, r *http.Request) error {
	pkgs, err := c.s.LoadPkgs()
	if err != nil {
		return err
	}
	if err := c.t.Execute(w, pkgs); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Client) handlePOST(w http.ResponseWriter, r *http.Request) error {
	if err := parseForm(r, "POST"); err != nil {
		return errors.WithStack(err)
	}

	name := r.Form.Get("name")
	if name == "" {
		const msg = "empty name provided"
		return errors.New(msg)
	}
	status := Ordered
	if r.Form.Has("shipped") {
		status = Shipped
	}
	p := NewPkg(
		name,
		WithInpost(r.Form.Has("inpost")),
		WithStatus(status),
	)
	if err := c.s.StorePkg(p); err != nil {
		return errors.WithStack(err)
	}
	log.Debug().Interface("pkg", p).Msg("Stored package")

	// r.URL.Path is needed in case Client listend on a different handle than "/"
	http.Redirect(w, r, r.URL.Path, http.StatusMovedPermanently)

	return nil
}

func (c *Client) handlePATCH(w http.ResponseWriter, r *http.Request) error {
	if err := parseForm(r, "PATCH"); err != nil {
		return errors.WithStack(err)
	}

	id, err := uuid.Parse(r.Form.Get("id"))
	if err != nil {
		return errors.WithStack(err)
	}

	status, err := StatusFromString(r.Form.Get("status"))
	if err != nil {
		return errors.WithStack(err)
	}

	if err := c.s.UpdatePkgStatus(id, status); err != nil {
		return errors.WithStack(err)
	}
	log.Debug().Interface("id", id).Interface("status", status).Msg("Updated status")

	w.Write([]byte(status.String()))

	return nil
}

func (c *Client) handleDELETE(w http.ResponseWriter, r *http.Request) error {
	// For some reason it is not possible to use URL encoded form with DELETE
	// using XMLHttpRequest in JS.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	s := string(b)
	log.Debug().Str("body", s).Str("method", "DELETE").Msg("Read body")

	id, err := uuid.Parse(s)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := c.s.DeletePkg(id); err != nil {
		return errors.WithStack(err)
	}
	log.Debug().Interface("id", id).Msg("Deleted package")

	return nil
}

func parseForm(r *http.Request, m string) error {
	if err := r.ParseForm(); err != nil {
		return errors.WithStack(err)
	}
	log.Debug().Interface("form", r.Form).Str("method", m).Msg("Parsed form")
	return nil
}
