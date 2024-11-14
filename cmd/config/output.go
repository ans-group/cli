package config

type Context struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type ContextCollection []Context

func (m ContextCollection) DefaultColumns() []string {
	return []string{"name", "active"}
}
