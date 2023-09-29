package handlers

import (
	"context"
	"github.com/lukachi/blob-svc/internal/data"
	"github.com/lukachi/blob-svc/internal/service/handlers/helpers"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	blobsQCtxKey
	usersQCtxKey
	jwtQCtxKey
	JWTUsersClaimCtxKey
	VerifiedBlobCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxBlobsQ(entry data.BlobsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, blobsQCtxKey, entry)
	}
}

func BlobsQ(r *http.Request) data.BlobsQ {
	return r.Context().Value(blobsQCtxKey).(data.BlobsQ).New()
}

func CtxUsersQ(entry data.UsersQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, usersQCtxKey, entry)
	}
}

func UsersQ(r *http.Request) data.UsersQ {
	return r.Context().Value(usersQCtxKey).(data.UsersQ).New()
}

func CtxJWT(entry helpers.JWTManager) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, jwtQCtxKey, entry)
	}
}

func JWT(r *http.Request) helpers.JWTManager {
	return r.Context().Value(jwtQCtxKey).(helpers.JWTManager).New()
}

func UserClaim(r *http.Request) *helpers.UserClaims {
	return r.Context().Value(JWTUsersClaimCtxKey).(*helpers.UserClaims)
}

func VerifiedBlob(r *http.Request) *data.Blob {
	return r.Context().Value(VerifiedBlobCtxKey).(*data.Blob)
}
