package retrieve

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"start1/internal/storage"
	resp "start1/lib/api/response"
	"start1/lib/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

type Response struct {
	URL string `json:"url"`
	resp.Response
}

//go:generate mockery --name=URLRetriever
type URLRetriever interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlRetriever URLRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "http-server.handlers.url.retrieve.New"
		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded")

		if err = validator.New().Struct(req); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErrs))

			return
		}

		url, err := urlRetriever.GetURL(req.Alias)
		log.Info(fmt.Sprintf("ERROR in retrieve: %s", url))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("url is not found", slog.String("alias", req.Alias))
				render.JSON(w, r, resp.Error("url is not found"))

				return
			}

			log.Error("failed to get URL", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get URL"))

			return
		}

		log.Info("url retrieved", slog.String("url", url))
		render.JSON(w, r, Response{
			URL:      url,
			Response: resp.OK(),
		})
	}
}
