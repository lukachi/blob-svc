package pg

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/lukachi/blob-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const usersTableName = "users"

func NewUsersQ(db *pgdb.DB) data.UsersQ {
	return &UsersQ{
		db:  db.Clone(),
		sql: squirrel.Select("b.*").From(fmt.Sprintf("%s as b", usersTableName)),
	}
}

type UsersQ struct {
	db  *pgdb.DB
	sql squirrel.SelectBuilder
}

func (q *UsersQ) New() data.UsersQ {
	return NewUsersQ(q.db)
}

func (q *UsersQ) Get() (*data.User, error) {
	var user data.User

	err := q.db.Get(&user, q.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (q *UsersQ) Select() ([]data.User, error) {
	var result []data.User

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *UsersQ) FilterById(id string) data.UsersQ {
	q.sql = q.sql.Where(squirrel.Eq{"id": id})

	return q
}

func (q *UsersQ) FilterByLogin(login string) data.UsersQ {
	q.sql = q.sql.Where(squirrel.Eq{"login": login})

	return q
}

func (q *UsersQ) Insert(data data.User) (string, error) {
	clauses := structs.Map(data)
	var id string

	stmt := squirrel.Insert(usersTableName).SetMap(clauses).Suffix("returning id")

	err := q.db.Get(&id, stmt)

	return id, err
}

func (q *UsersQ) Delete(id ...string) error {
	s := squirrel.Delete(usersTableName).Where(squirrel.Eq{"id": id})
	err := q.db.Exec(s)

	return err
}
