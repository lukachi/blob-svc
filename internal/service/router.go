package service

import (
	"github.com/go-chi/chi"
	"github.com/lukachi/blob-svc/internal/config"
	"github.com/lukachi/blob-svc/internal/data/pg"
	"github.com/lukachi/blob-svc/internal/service/handlers"
	"github.com/lukachi/blob-svc/internal/service/handlers/helpers"
	"github.com/lukachi/blob-svc/internal/service/handlers/middlewares"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxBlobsQ(pg.NewBlobsQ(cfg.DB())),
			handlers.CtxUsersQ(pg.NewUsersQ(cfg.DB())),
			handlers.CtxJWT(helpers.NewJwtManager([]byte(cfg.Secret()))),
		),
	)
	r.Route("/blob-svc", func(r chi.Router) {
		// configure endpoints here
		r.With(middlewares.VerifyAccessToken()).Post("/", handlers.CreateBlob)
		r.With(middlewares.VerifyAccessToken()).Get("/{id}", handlers.GetBlob)
		r.With(middlewares.VerifyAccessToken()).Delete("/{id}", handlers.DeleteBlobById)

		r.Post("/sign-up", handlers.SignUp)
		r.Post("/sign-in", handlers.SignIn)
		r.Post("/refresh", handlers.Refresh)
	})

	return r
}
