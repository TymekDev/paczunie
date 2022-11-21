package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"packages/packages"
)

func main() {
	port := flag.Int("p", 8080, "port to listen on")
	debug := flag.Bool("debug", false, "sets log level to debug")
	dbName := flag.String("db", "packages.db", "path to SQLite3 database")
	flag.Parse()

	const timeFormat = "2006-01-02 15:04 -0700"
	zerolog.TimeFieldFormat = timeFormat
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{TimeFormat: timeFormat, Out: os.Stderr})
	switch {
	case *debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	c, err := packages.NewClientWithSQLiteStorage(*dbName)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Int("port", *port).Msg("Started listening")
	if err := http.ListenAndServe(fmt.Sprint(":", *port), c); err != nil {
		log.Fatal().Err(err).Send()
	}
}
