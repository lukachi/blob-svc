package data

type UsersQ interface {
	New() UsersQ

	Get() (*User, error)
	Select() ([]User, error)

	FilterById(id string) UsersQ
	FilterByLogin(login string) UsersQ

	Insert(data User) (string, error)
	Delete(id ...string) error
}

type User struct {
	ID       string `db:"id" structs:"id"`
	Login    string `db:"login" structs:"login"`
	Password string `db:"password" structs:"password"`
	Username string `db:"username" structs:"username"`
}
