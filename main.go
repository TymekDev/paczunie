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

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	switch {
	case *debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	c, err := packages.NewClientWithSQLiteStorage(*dbName)
	if err != nil {
		log.Fatal().Stack().Err(err).Send()
	}

	log.Info().Int("port", *port).Msg("Started listening")
	if err := http.ListenAndServe(fmt.Sprint(":", *port), c); err != nil {
		log.Fatal().Stack().Err(err).Send()
	}
}
