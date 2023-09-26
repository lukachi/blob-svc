package data

type BlobsQ interface {
	New() BlobsQ

	Get() (*Blob, error)
	Select() ([]Blob, error)
	FilterById(id string) (*Blob, error)

	Insert(data Blob) (string, error)
	Delete(id ...string) error
}

type Blob struct {
	ID      string `db:"id" structs:"id"`
	Value   string `db:"value" structs:"value"`
	OwnerId string `db:"owner_id" structs:"owner_id"`
}
