package handlers

import (
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/lukachi/blob-svc/internal/data"
	"github.com/lukachi/blob-svc/resources"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

func NewGetBlobModel(blob data.Blob) resources.GetBlob {
	result := resources.GetBlob{
		Key: resources.Key{
			ID:   blob.ID,
			Type: resources.BLOB,
		},
		Attributes: resources.GetBlobAttributes{
			Value: string(blob.Value),
		},
		Relationships: resources.GetBlobRelationships{
			Owner: resources.Relation{
				Data: &resources.Key{
					ID:   blob.OwnerId,
					Type: resources.USER,
				},
			},
		},
	}

	return result
}

func GetBlob(w http.ResponseWriter, r *http.Request) {
	req, headers, err := NewGetBlobRequest(r)

	if err != nil {
		Log(r).WithError(err).Error("Failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userClaims, err := JWT(r).ParseAccessToken(headers.Get("Authorization"))

	if err != nil {
		Log(r).WithError(err).Error("Unauthorized")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	blob, err := BlobsQ(r).FilterById(req.Id)

	if err != nil {
		Log(r).WithError(err).Error("Failed to get blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if blob == nil {
		Log(r).WithField("id", req.Id).Error("Blob not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if blob.OwnerId != userClaims.ID {
		Log(r).WithError(err).Error("Unauthorized")
		ape.RenderErr(w, problems.Unauthorized())
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

	user, err := UsersQ(r).FilterById(blob.OwnerId)

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

type GetBlobRequest struct {
	Id          string `url:"-"`
	IncludeUser bool   `include:"user"`
}

func NewGetBlobRequest(r *http.Request) (GetBlobRequest, *http.Header, error) {
	/* because in generated resources response type == request type */
	request := GetBlobRequest{}

	headers := r.Header

	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		Log(r).WithError(err).Error("Failed to parse id")
		return request, &headers, err
	}

	err := urlval.Decode(r.URL.Query(), &request)

	if err != nil {
		return request, &headers, errors.Wrap(err, "failed to unmarshal")
	}

	request.Id = id

	return request, &headers, nil
}
