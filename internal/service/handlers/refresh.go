package handlers

import (
	"database/sql"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lukachi/blob-svc/internal/service/handlers/helpers"
	"github.com/lukachi/blob-svc/resources"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"time"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	req, err := NewRefreshRequest(r)

	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	accessToken := r.Header.Get("Authorization")

	// TODO fix this
	if accessToken == "" {
		Log(r).WithError(err).Error("access token is not provided")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	_, accessTokenClaims, err := JWT(r).ParseAccessToken(accessToken)

	isRefreshTokenValid, refreshTokenClaims, err := JWT(r).ParseRefreshToken(req.RefreshToken)

	if err != nil {
		Log(r).WithError(err).Error("failed to parse refresh token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if !isRefreshTokenValid {
		Log(r).WithError(err).Error("refresh token is not valid")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	if time.Until(refreshTokenClaims.ExpiresAt.Time) > 30*time.Second {
		Log(r).WithError(err).Error("refresh token is not refreshable now")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := UsersQ(r).FilterById(accessTokenClaims.ID).Get()

	if err == sql.ErrNoRows || user == nil {
		Log(r).WithError(err).Error("user not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if err != nil {
		Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	authTokens, err := JWT(r).Gen(&helpers.UserClaims{
		ID:       user.ID,
		Username: user.Username,
	})

	if err != nil {
		Log(r).WithError(err).Error("failed to generate auth tokens")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.AuthTokensResponse{
		Data: newAuthTokensModel(&authTokens),
	}

	ape.Render(w, result)
}

func NewRefreshRequest(r *http.Request) (resources.RefreshRequestAttributes, error) {
	request := struct {
		Data resources.RefreshRequest `json:"data"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data.Attributes, errors.Wrap(err, "failed to unmarshal")
	}

	return request.Data.Attributes, validateRefreshRequest(request.Data)
}

func validateRefreshRequest(r resources.RefreshRequest) error {
	return validation.Errors{
		"/data/attributes/refresh_request": validation.Validate(&r.Attributes.RefreshToken, validation.Required),
	}.Filter()
}
