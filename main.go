package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	port := flag.Int("p", 8080, "port to listen on")
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	switch {
	case *debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	db, err := sql.Open("sqlite3", "bazka.db")
	if err != nil {
		log.Fatal().Stack().Err(err).Send()
	}

	c, err := NewClient(newDBStorage(db))
	if err != nil {
		log.Fatal().Stack().Err(err).Send()
	}

	log.Info().Int("port", *port).Msg("Started listening")
	if err := http.ListenAndServe(fmt.Sprint(":", *port), c); err != nil {
		log.Fatal().Stack().Err(err).Send()
	}
}
