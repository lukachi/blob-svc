package service

import (
  "github.com/lukachi/blob-svc/internal/service/handlers"
  "github.com/go-chi/chi"
  "gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
  r := chi.NewRouter()

  r.Use(
    ape.RecoverMiddleware(s.log),
    ape.LoganMiddleware(s.log),
    ape.CtxMiddleware(
      handlers.CtxLog(s.log),
    ),
  )
  r.Route("/integrations/blob-svc", func(r chi.Router) {
    // configure endpoints here
  })

  return r
}
