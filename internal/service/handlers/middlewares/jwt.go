package middlewares

import (
	"context"
	"github.com/lukachi/blob-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func VerifyAccessToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				handlers.Log(r).Error("no access token provided")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			isAccessTokenValid, userClaims, err := handlers.JWT(r).ParseAccessToken(r.Header.Get("Authorization"))

			if !isAccessTokenValid {
				handlers.Log(r).WithError(err).Error("access token is not valid")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			if err != nil {
				handlers.Log(r).WithError(err).Error("failed to parse access token")
				ape.RenderErr(w, problems.InternalError())
				return
			}

			ctx := context.WithValue(r.Context(), handlers.JWTUsersClaimCtxKey, userClaims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
