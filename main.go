package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	port := flag.Int("p", 8080, "port to listen on")
	debug := flag.Bool("debug", false, "sets log level to debug")
	trace := flag.Bool("trace", false, "sets log level to trace")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	switch {
	case *trace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case *debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	r := mux.NewRouter()
	if err := Register(r); err != nil {
		log.Fatal().Stack().Err(err).Send()
	}

	log.Info().Int("port", *port).Msg("Started listening")
	if err := http.ListenAndServe(fmt.Sprint(":", *port), r); err != nil {
		log.Fatal().Stack().Err(err).Send()
	}
}
