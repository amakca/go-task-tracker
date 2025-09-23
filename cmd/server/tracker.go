package main

import (
	"context"
	"net/http"

	"log/slog"
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
		slog.Error("http server failed", "err", err)
		panic(err)
	}
}

func (t *Tracker) Shutdown(ctxShutdown context.Context) {
	if err := t.srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("server shutdown error", "err", err)
	}
}
