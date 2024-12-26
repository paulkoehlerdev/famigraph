package repository

type OTC interface {
	Generate() (string, error)
}
