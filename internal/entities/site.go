package entities

import (
	"github.com/gofrs/uuid"
)

type Site struct {
	UUID    uuid.UUID `json:"-" db:"credential_id"`
	SiteURL string    `json:"site" db:"site_url"`
	SLogin  string    `json:"login" db:"login"` // SLogin mean site login (saving creadentials, ie not user login.
	SPw     string    `json:"pw" db:"pw"`
}
