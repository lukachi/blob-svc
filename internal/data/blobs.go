package data

import (
	"gitlab.com/tokend/go/xdr"
)

type BlobsQ interface {
	New() BlobsQ

	Get() (*Blob, error)
	Select() ([]Blob, error)
	FilterById(id string) BlobsQ

	Insert(data Blob) (string, error)
	Delete(id ...string) error
}

type HorizonBlobsQ interface {
	New() HorizonBlobsQ

	WriteBlob(blob *Blob) (xdr.Uint64, error)
	GetBlob(string) (*Blob, error)
}

type Blob struct {
	ID      string `db:"id" structs:"id"`
	Value   string `db:"value" structs:"value"`
	OwnerId string `db:"owner_id" structs:"owner_id"`
}
