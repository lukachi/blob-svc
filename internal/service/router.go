package service

import (
	"github.com/go-chi/chi"
	"github.com/lukachi/blob-svc/internal/config"
	"github.com/lukachi/blob-svc/internal/data/horizon"
	"github.com/lukachi/blob-svc/internal/data/pg"
	"github.com/lukachi/blob-svc/internal/service/handlers"
	"github.com/lukachi/blob-svc/internal/service/handlers/helpers"
	"github.com/lukachi/blob-svc/internal/service/handlers/middlewares"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/tokend/keypair"
	"os"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	masterSeed := os.Getenv("MASTER_SEED")
	if len(masterSeed) == 0 {
		panic("the MASTER_SEED enviroment variable does not exist")
	}

	masterKp, err := keypair.ParseSeed(masterSeed)
	if err != nil {
		panic(err)
	}

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxBlobsQ(pg.NewBlobsQ(cfg.DB())),
			handlers.CtxHorizonBLobsQ(horizon.NewBlobsQ(cfg.Log(), cfg.Horizon(), &masterKp)),
			handlers.CtxUsersQ(pg.NewUsersQ(cfg.DB())),
			handlers.CtxJWT(helpers.NewJwtManager([]byte(cfg.Secret()))),
		),
	)
	r.Route("/blob-svc", func(r chi.Router) {
		// configure endpoints here

		r.Group(func(r chi.Router) {
			r.Use(middlewares.VerifyAccessToken())

			r.Post("/", handlers.CreateBlob)
			r.Post("/submit", handlers.SubmitBlob)
			r.Get("/submitted/{id}", handlers.GetSubmittedBlobById)

			r.Group(func(r chi.Router) {
				r.Use(middlewares.VerifyBlobOwner())

				r.Get("/{id}", handlers.GetBlob)
				r.Delete("/{id}", handlers.DeleteBlobById)
			})
		})

		r.Post("/sign-up", handlers.SignUp)
		r.Post("/sign-in", handlers.SignIn)
		r.Post("/refresh", handlers.Refresh)
	})

	return r
}
