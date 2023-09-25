package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/lukachi/blob-svc/internal/data"
	"github.com/lukachi/blob-svc/internal/service/handlers/helpers"
	"github.com/lukachi/blob-svc/resources"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func newAuthTokensModel(authTokens *data.AuthTokens) resources.AuthTokens {
	result := resources.AuthTokens{
		Key: resources.Key{
			ID:   uuid.NewString(),
			Type: resources.AUTH_TOKENS,
		},
		Attributes: resources.AuthTokensAttributes{
			AccessToken:  authTokens.AccessToken,
			RefreshToken: authTokens.RefreshToken,
		},
	}

	return result
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	request, err := newSignUpRequest(r)

	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	salt := uuid.NewString()

	hashedPassword := sha256.Sum256([]byte(request.Password))

	password := fmt.Sprintf("%x", hashedPassword[:]) + salt

	user := data.User{
		ID:       uuid.NewString(),
		Login:    request.Login,
		Password: password,
		Salt:     salt,
		Username: request.Username,
	}

	user2, err := UsersQ(r).FilterByLogin(request.Login)

	if err != nil {
		Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user2 != nil && user2.Login == user.Login {
		Log(r).WithError(err).Error("user already exists")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	user.ID, err = UsersQ(r).Insert(user)

	if err != nil {
		Log(r).WithError(err).Error("failed to insert user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	userClaims := helpers.UserClaims{
		ID:       user.ID,
		Username: user.Username,
	}

	authTokens, err := JWT(r).Gen(userClaims)

	if err != nil {
		Log(r).WithError(err).Error("failed to generate auth tokens")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	hashedRefreshToken := sha256.Sum256([]byte(authTokens.RefreshToken))

	_, err = SessionsQ(r).Insert(data.Session{
		ID: fmt.Sprintf("%x", hashedRefreshToken[:]),
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

func newSignUpRequest(r *http.Request) (resources.SignUpRequestAttributes, error) {
	request := struct {
		Data resources.SignUpRequest `json:"data"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data.Attributes, errors.Wrap(err, "failed to unmarshal")
	}

	return request.Data.Attributes, validateSignUpRequest(request.Data)
}

func validateSignUpRequest(r resources.SignUpRequest) error {
	return validation.Errors{
		"/data/attributes/login":    validation.Validate(&r.Attributes.Login, validation.Required),
		"/data/attributes/password": validation.Validate(&r.Attributes.Password, validation.Required),
		"/data/attributes/username": validation.Validate(&r.Attributes.Username, validation.Required),
	}.Filter()
}
