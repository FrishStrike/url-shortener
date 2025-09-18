package redirect

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"start1/internal/http-server/handlers/url/retrieve"
	"start1/internal/storage"
	resp "start1/lib/api/response"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func New(log *slog.Logger, urlRetriever retrieve.URLRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "http-server.handlers.redirect.New"
		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias is empty"))
			return
		}

		url, err := urlRetriever.GetURL(alias)
		log.Info(fmt.Sprintf("Url: %s, alias: %s", url, alias))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("alias is not found")
				render.JSON(w, r, resp.Error("alias is not found"))
				return
			}

			log.Error("invalid request")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		log.Info("got url", slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
