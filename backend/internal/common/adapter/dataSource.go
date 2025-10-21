package adapter

type DataSourceProperties interface {
	GetURL() string
	GetUsername() string
	GetPassword() string
	GetCatalog() string
	GetSchema() string
}
