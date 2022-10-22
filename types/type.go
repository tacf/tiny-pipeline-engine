package types

type Plugin interface {
        Exec(map[string]string)
        GetName() string
}

