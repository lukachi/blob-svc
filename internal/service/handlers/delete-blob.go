package handlers

import (
	"database/sql"
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

	if err != nil {
		Log(r).WithError(err).Error("Failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userClaims := UserClaim(r)

	blob, err := BlobsQ(r).FilterById(req.ID).Get()

	if err == sql.ErrNoRows || blob == nil {
		Log(r).WithField("id", req.ID).Error("Blob not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if err != nil {
		Log(r).WithError(err).Error("Failed to get blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if blob.OwnerId != userClaims.ID {
		Log(r).WithError(err).Error("Unauthorized")
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

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
