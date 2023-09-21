package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lukachi/blob-svc/internal/service/handlers/helpers"
	"github.com/lukachi/blob-svc/resources"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	request, err := NewSignInRequest(r)

	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := UsersQ(r).FilterByLogin(request.Login)

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

	hashedPassword := sha256.Sum256([]byte(request.Password))
	password := fmt.Sprintf("%x", hashedPassword[:]) + user.Salt

	if password != user.Password {
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

	result := resources.AuthTokensResponse{
		Data: newAuthTokensModel(&authTokens),
	}

	ape.Render(w, result)
}

func NewSignInRequest(r *http.Request) (resources.SignInRequestAttributes, error) {
	request := struct {
		Data resources.SignInRequest `json:"data"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data.Attributes, errors.Wrap(err, "failed to unmarshal")
	}

	return request.Data.Attributes, validateSignInRequest(request.Data)
}

func validateSignInRequest(r resources.SignInRequest) error {
	return validation.Errors{
		"/data/attributes/login":    validation.Validate(&r.Attributes.Login, validation.Required),
		"/data/attributes/password": validation.Validate(&r.Attributes.Password, validation.Required),
	}.Filter()
}
