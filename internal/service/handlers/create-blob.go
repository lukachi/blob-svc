package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/lukachi/blob-svc/internal/data"
	"github.com/lukachi/blob-svc/resources"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

type CreateBlobRequest struct {
	Value json.RawMessage `json:"value"`
}

func newBlobModel(blob data.Blob) resources.Blob {
	result := resources.Blob{
		Key: resources.Key{
			ID:           blob.ID,
			ResourceType: "Blobs",
		},
		Attributes: resources.BlobAttributes{
			Value: string(blob.Value),
		},
	}

	return result
}

func CreateBlob(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateBlobRequest(r)

	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	blob := data.Blob{
		ID:    uuid.NewString(),
		Value: string(request.Value),
	}

	blob.ID, err = BlobsQ(r).Insert(blob)

	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.BlobResponse{
		Data: newBlobModel(blob),
	}

	ape.Render(w, result)
}

func NewCreateBlobRequest(r *http.Request) (CreateBlobRequest, error) {
	request := struct {
		Data CreateBlobRequest `json:"data"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data, errors.Wrap(err, "failed to unmarshal")
	}

	return request.Data, request.Data.validate()
}

func (r CreateBlobRequest) validate() error {
	return validation.Errors{
		"/data/attributes/value": validation.Validate(&r.Value, validation.Required),
	}.Filter()
}
