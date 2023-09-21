package handlers

import (
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/lukachi/blob-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

type GetBlobByIDRequest struct {
	ID string
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

	blob, err := BlobsQ(r).FilterById(req.ID)

	if err != nil {
		Log(r).WithError(err).Error("Failed to get blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if blob == nil {
		Log(r).WithField("id", req.ID).Error("Blob not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if blob.Owner != userClaims.ID {
		Log(r).WithError(err).Error("Unauthorized")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	result := resources.BlobResponse{
		Data: newBlobModel(*blob),
	}

	ape.Render(w, result)
}

func NewGetBlobRequest(r *http.Request) (GetBlobByIDRequest, *http.Header, error) {
	request := GetBlobByIDRequest{}

	headers := r.Header

	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		Log(r).WithError(err).Error("Failed to parse id")
		return request, &headers, err
	}

	request.ID = id

	return request, &headers, nil
}
