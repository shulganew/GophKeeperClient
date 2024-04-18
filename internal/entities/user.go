package entities

import "database/sql"

type User struct {
	JWT      sql.NullString `json:"-"`
	Login    string         `json:"login"`
	Password string         `json:"password"`
}
