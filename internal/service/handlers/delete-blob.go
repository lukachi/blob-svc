package handlers

import (
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

type DeleteBlobByIdRequest struct {
	ID string
}

func DeleteBlobById(w http.ResponseWriter, r *http.Request) {
	req, err := NewDeleteBlobByIdRequest(r)

	err = BlobsQ(r).Delete(req.ID)

	if err != nil {
		Log(r).WithError(err).Error("Failed to delete blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusNoContent)
}

func NewDeleteBlobByIdRequest(r *http.Request) (DeleteBlobByIdRequest, error) {
	request := DeleteBlobByIdRequest{}

	id := chi.URLParam(r, "id")

	if _, err := uuid.Parse(id); err != nil {
		Log(r).WithError(err).Error("Failed to parse id")
		return request, err
	}

	request.ID = id

	return request, nil
}
