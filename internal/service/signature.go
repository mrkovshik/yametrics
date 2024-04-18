package service

type Signature interface {
	Generate() (string, error)
}
