package add

import (
	resp "crypto_ex_rate/pkg/api/response"
	"crypto_ex_rate/pkg/logger/sl"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Coin string `json:"coin" validate:"required"`
}

type Response struct {
	resp.Response
}

type CryptoExRate interface {
	AddCurrency(coin string) error
}

func New(log *slog.Logger, cryptoExRate CryptoExRate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http.handlers.add.New"
		log := log.With(slog.String("op", op))

		var req Request

		// Decode data from request JSON
		err := render.DecodeJSON(r.Body, &req)
		// Exception: empty request
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			w.WriteHeader(400)
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(400)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Validate required request fields
		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			w.WriteHeader(400)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		err = cryptoExRate.AddCurrency(req.Coin)
		if err != nil {
			log.Error("failed to add cryptocurrency for monitoring", sl.Err(err))
			w.WriteHeader(500)
			render.JSON(w, r, resp.Error("failed to add cryptocurrency for monitoring"))
			return
		}

		render.JSON(w, r, resp.OK())
	}
}
