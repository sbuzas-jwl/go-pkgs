package main

import "log/slog"

type LogErrorHandler struct {
	log *slog.Logger
}

func (eh LogErrorHandler) Handle(err error) {
	eh.log.With("error", err).Error("unexpected error encountered.")
}
