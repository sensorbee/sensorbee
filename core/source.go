package core

type Source interface {
	GenerateStream(w Writer) error
	Schema() *Schema
}
