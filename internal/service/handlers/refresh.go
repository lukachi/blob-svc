package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lukachi/blob-svc/internal/data"
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

	hashedRefreshToken := sha256.Sum256([]byte(req.RefreshToken))

	session, err := SessionsQ(r).FilterById(fmt.Sprintf("%x", hashedRefreshToken[:]))

	if err != nil {
		Log(r).WithError(err).Error("failed to get session")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if session == nil {
		Log(r).WithError(err).Error("session not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if err := validateSession(session); err != nil {
		Log(r).WithError(err).Error("session expired")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	user, err := UsersQ(r).FilterById(session.UserID)

	if err != nil {
		Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user == nil {
		Log(r).WithError(err).Error("user not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	authTokens, err := JWT(r).Gen(helpers.UserClaims{
		ID:       user.ID,
		Username: user.Username,
	})

	if err != nil {
		Log(r).WithError(err).Error("failed to generate auth tokens")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	hashedNewRefreshToken := sha256.Sum256([]byte(authTokens.RefreshToken))

	_, err = SessionsQ(r).Insert(data.Session{
		ID:        fmt.Sprintf("%x", hashedNewRefreshToken[:]),
		UserID:    user.ID,
		ExpiresAt: authTokens.ExpiresAt,
		CreatedAt: authTokens.CreatedAt,
	})

	if err != nil {
		Log(r).WithError(err).Error("failed to insert session")
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

func validateSession(s *data.Session) error {
	if time.Now().After(s.ExpiresAt) {
		return errors.New("session expired")
	}

	return nil
}
