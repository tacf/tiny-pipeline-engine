package types

type Plugin interface {
	Exec()
	GetName() string
}
