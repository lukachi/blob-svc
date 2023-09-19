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
	req, err := NewGetBlobRequest(r)

	if err != nil {
		Log(r).WithError(err).Error("Failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
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

	result := resources.BlobResponse{
		Data: newBlobModel(*blob),
	}

	ape.Render(w, result)
}

func NewGetBlobRequest(r *http.Request) (GetBlobByIDRequest, error) {
	request := GetBlobByIDRequest{}

	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		Log(r).WithError(err).Error("Failed to parse id")
		return request, err
	}

	request.ID = id

	return request, nil
}
