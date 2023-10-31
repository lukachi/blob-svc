package handlers

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/lukachi/blob-svc/resources"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

func GetSubmittedBlobById(w http.ResponseWriter, r *http.Request) {
	req, err := NewGetSubmittedBlobRequest(r)

	if err != nil {
		Log(r).WithError(err).Error("Failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	blob, err := HorizonBlobsQ(r).GetBlob(req.Id)

	if err != nil {
		Log(r).WithError(err).Error("Failed to get blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if blob == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var response resources.GetBlobResponse

	if !req.IncludeUser {
		response = resources.GetBlobResponse{
			Data: NewGetBlobModel(*blob),
		}

		ape.Render(w, response)

		return
	}

	user, err := UsersQ(r).FilterById(blob.OwnerId).Get()

	if err == sql.ErrNoRows || user == nil {
		Log(r).WithField("id", blob.OwnerId).Error("User not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if err != nil {
		Log(r).WithError(err).Error("Failed to get user for includes")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	userResponse := resources.User{
		Key: resources.Key{
			ID:   user.ID,
			Type: resources.USER,
		},
		Attributes: resources.UserAttributes{
			Username: user.Username,
		},
	}

	includes := resources.Included{}

	includes.Add(&userResponse)

	response = resources.GetBlobResponse{
		Data:     NewGetBlobModel(*blob),
		Included: includes,
	}

	ape.Render(w, response)
}

func NewGetSubmittedBlobRequest(r *http.Request) (GetBlobRequest, error) {
	/* because in generated resources response type == request type */
	request := GetBlobRequest{}

	id := chi.URLParam(r, "id")

	if id == "" {
		return request, errors.New("id is required")
	}

	err := urlval.Decode(r.URL.Query(), &request)

	if err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	request.Id = id

	return request, nil
}
