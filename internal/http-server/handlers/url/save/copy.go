package save

import (
	"crypto/rand"
	"errors"
	"log/slog"
	"net/http"
	"start1/internal/storage"
	resp "start1/lib/api/response"
	"start1/lib/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Requestt struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Responsee struct {
	resp.Response `json:"response"`
	Alias         string `json:"alias,omitempty"`
}

type URLSaverr interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func Neww(log slog.Logger, urlSaver URLSaverr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "http-server.handlers.url.save.Neww"
		log = *slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Requestt
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request was decoded")

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = rand.Text()
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Error("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, resp.Error("url already exists"))

				return
			}
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))
		render.JSON(w, r, ResponseOK(alias))

	}
}

func ResponseOK(alias string) *Response {
	return &Response{
		Response: resp.OK(),
		Alias:    alias,
	}
}
