package data

import "time"

type SessionsQ interface {
	New() SessionsQ

	Get() (*Session, error)
	Select() ([]Session, error)
	FilterById(id string) (*Session, error)

	Insert(data Session) (string, error)
	Delete(id ...string) error
}

type Session struct {
	ID        string    `db:"id" structs:"id"`
	UserID    string    `db:"user_id" structs:"user_id"`
	CreatedAt time.Time `db:"created_at" structs:"created_at"`
	ExpiresAt time.Time `db:"expires_at" structs:"expires_at"`
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}
