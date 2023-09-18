package pg

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/lukachi/blob-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const blobsTableName = "blobs"

type BlobsQ struct {
	db  *pgdb.DB
	sql squirrel.SelectBuilder
}

func NewBlobsQ(db *pgdb.DB) data.BlobsQ {
	return &BlobsQ{
		db:  db.Clone(),
		sql: squirrel.Select("b.*").From(fmt.Sprintf("%s as b", blobsTableName)),
	}
}

func (q *BlobsQ) New() data.BlobsQ {
	return NewBlobsQ(q.db)
}

func (q *BlobsQ) Get() (*data.Blob, error) {
	var blob data.Blob

	err := q.db.Get(&blob, q.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &blob, err
}

func (q *BlobsQ) Select() ([]data.Blob, error) {
	var result []data.Blob

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *BlobsQ) Insert(data data.Blob) (string, error) {
	clauses := structs.Map(data)
	var id string // FIXME: use uuidv4

	stmt := squirrel.Insert(blobsTableName).SetMap(clauses).Suffix("returning id")
	err := q.db.Get(&id, stmt)

	return id, err
}

func (q *BlobsQ) Delete(id ...string) error {
	s := squirrel.Delete(blobsTableName).Where(squirrel.Eq{"id": id})
	err := q.db.Exec(s)

	return err
}
