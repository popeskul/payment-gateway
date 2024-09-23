package ports

type DatabaseConfig interface {
	GetHost() string
	GetPort() string
	GetUser() string
	GetPassword() string
	GetDBName() string
	GetMaxConnections() int
	GetMinConnections() int
	GetMaxConnLifetime() int
	GetMaxConnIdleTime() int
}
