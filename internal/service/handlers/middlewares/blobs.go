package middlewares

import (
	"context"
	"database/sql"
	"github.com/lukachi/blob-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func VerifyBlobOwner() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, err := handlers.NewGetBlobRequest(r)

			if err != nil {
				handlers.Log(r).WithError(err).Error("Failed to parse request")
				ape.RenderErr(w, problems.BadRequest(err)...)
				return
			}

			userClaims := handlers.UserClaim(r)

			blob, err := handlers.BlobsQ(r).FilterById(req.Id).Get()

			if err == sql.ErrNoRows || blob == nil {
				handlers.Log(r).WithField("id", req.Id).Error("Blob not found")
				ape.RenderErr(w, problems.NotFound())
				return
			}

			if err != nil {
				handlers.Log(r).WithError(err).Error("Failed to get blob")
				ape.RenderErr(w, problems.InternalError())
				return
			}

			if blob.OwnerId != userClaims.ID {
				handlers.Log(r).WithError(err).Error("User is not the owner of the blob")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			ctx := context.WithValue(r.Context(), handlers.VerifiedBlobCtxKey, blob)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
