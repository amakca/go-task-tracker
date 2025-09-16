package main

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Tracker struct {
	srv *http.Server
}

func NewTracker(srv *http.Server) *Tracker {
	tracker := &Tracker{
		srv: srv,
	}

	return tracker
}

func (t *Tracker) Start() {
	if err := t.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("http server failed")
	}
}

func (t *Tracker) Shutdown(ctxShutdown context.Context) {
	if err := t.srv.Shutdown(ctxShutdown); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}
}
