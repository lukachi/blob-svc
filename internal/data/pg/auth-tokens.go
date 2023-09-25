package pg

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/lukachi/blob-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const sessionsTableName = "sessions"

type SessionsQ struct {
	db  *pgdb.DB
	sql squirrel.SelectBuilder
}

func NewSessionsQ(db *pgdb.DB) data.SessionsQ {
	return &SessionsQ{
		db:  db.Clone(),
		sql: squirrel.Select("s.*").From(fmt.Sprintf("%s as s", sessionsTableName)),
	}
}

func (s SessionsQ) New() data.SessionsQ {
	return NewSessionsQ(s.db)
}

func (s SessionsQ) Get() (*data.Session, error) {
	var session data.Session

	err := s.db.Get(&session, s.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &session, err
}

func (s SessionsQ) Select() ([]data.Session, error) {
	var result []data.Session

	err := s.db.Select(&result, s.sql)

	if len(result) == 0 {
		return nil, nil
	}

	return result, err
}

func (s SessionsQ) FilterById(id string) (*data.Session, error) {
	var result []data.Session

	err := s.db.Select(&result, s.sql.Where(squirrel.Eq{"id": id}))

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], err
}

func (s SessionsQ) Insert(data data.Session) (string, error) {
	clauses := structs.Map(data)
	var id string

	stmt := squirrel.Insert(sessionsTableName).SetMap(clauses).Suffix("returning id")
	err := s.db.Get(&id, stmt)

	return id, err
}

func (s SessionsQ) Delete(id ...string) error {
	stmt := squirrel.Delete(sessionsTableName).Where(squirrel.Eq{"id": id})
	err := s.db.Exec(stmt)

	return err
}
