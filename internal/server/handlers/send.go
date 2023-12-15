package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/smtp"
	"x-platform-emailer/internal/config"
	"x-platform-emailer/internal/lib/logger/sl"
	resp "x-platform-emailer/internal/lib/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	To  []string `json:"to"`
	Msg string   `json:"msg"`
}

type Response struct {
	resp.Response
}

const (
	PackagePath = "internal.server.handlers.send."
)

func NewSendHandler(log *slog.Logger, cfg *config.Mailbox) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = PackagePath + "NewSendHandler"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty", sl.Err(err))
				resp.SendResponse(w, r, http.StatusBadRequest, resp.Err("empty request body"))
				return
			}
			log.Error("failed to decode request body", sl.Err(err))
			resp.SendResponse(w, r, http.StatusBadRequest, resp.Err("failed to decode request body"))
			return
		}

		auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host)

		err = smtp.SendMail(cfg.Host+":"+cfg.Port, auth, cfg.From, req.To, []byte(req.Msg))
		if err != nil {
			log.Error("failed to send email", sl.Err(err))
			resp.SendResponse(w, r, http.StatusInternalServerError, resp.Err("failed to send email"))
			return
		}

		log.Info("message have been sent")
		resp.SendResponse(w, r, http.StatusOK, resp.OK())
	}
}
