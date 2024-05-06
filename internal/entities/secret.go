package entities

// Type of secter data.
type SecretType string

const (
	SITE SecretType = "SITE"
	CARD SecretType = "CARD"
	TEXT SecretType = "TEXT"
	FILE SecretType = "BIN"
	ANY  SecretType = "ANY" // State no metter.
)

func (s *SecretType) String() string {
	return string(*s)
}
