package entities

import "database/sql"

type User struct {
	JWT      sql.NullString `json:"jwt"`
	Login    string         `json:"login"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
}
