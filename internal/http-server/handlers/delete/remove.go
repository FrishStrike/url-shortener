package remove

import (
	"log/slog"
	"net/http"
	resp "start1/lib/api/response"
	"start1/lib/sl"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type URLRemover interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "http-server.handlers.delete.New"
		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Error("alias is empty")
			render.JSON(w, r, resp.Error("alias is empty"))
			return
		}

		err := urlRemover.DeleteURL(alias)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		log.Info("alias was deleted", slog.String("alias", alias))
		render.JSON(w, r, resp.Response{Status: resp.StatusOK})
	}
}
