package update

import (
	"errors"
	"net/http"

	"github.com/Svoevolin/url-shortener/internal/database"
	resp "github.com/Svoevolin/url-shortener/internal/lib/api/response"
	"github.com/Svoevolin/url-shortener/internal/lib/logger/sl"
	"github.com/Svoevolin/url-shortener/internal/lib/random"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLUpdater
type URLUpdater interface {
	UpdateURL(url string, newAlias string) (int64, error)
}

func New(log *slog.Logger, urlUpdater URLUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.update.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err = validator.New().Struct(req)
		if err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomStrings(aliasLength)
		}

		id, err := urlUpdater.UpdateURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, database.ErrURLNotFound) {
				log.Info("url not found", slog.String("alias", alias))

				render.JSON(w, r, resp.Error("url not found"))

				return
			}

			log.Error("failed to update url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to update url"))

			return
		}

		log.Info("url updated", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
