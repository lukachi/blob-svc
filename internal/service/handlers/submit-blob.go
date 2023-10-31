package handlers

import (
	"github.com/google/uuid"
	"github.com/lukachi/blob-svc/internal/data"
	"github.com/lukachi/blob-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func SubmitBlob(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateBlobRequest(r)

	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userClaims := UserClaim(r)

	blob := data.Blob{
		ID:      uuid.NewString(),
		Value:   string(request.Value),
		OwnerId: userClaims.ID,
	}

	horizonBlobsQ := HorizonBlobsQ(r)

	id, err := horizonBlobsQ.WriteBlob(&blob)

	blob.ID = strconv.FormatUint(uint64(id), 10)

	if err != nil {
		Log(r).WithError(err).Error("failed to write blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if err != nil {
		Log(r).WithError(err).Error("failed to insert blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.BlobResponse{
		Data: newBlobModel(blob),
	}

	ape.Render(w, result)
}
