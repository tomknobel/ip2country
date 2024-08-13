package db

import "github.com/tomknobel/ip2country/internal/models"

type Db interface {
	Connect() error
	Close() error
	Find(ip string) (*models.Country, error)
}
type DbConfig struct {
	ConnectionString string
}

func DbFactory(dbType string, config DbConfig) Db {
	switch dbType {
	case "csv":
		return NewCsvDb(config.ConnectionString)
	default:
		return NewCsvDb(config.ConnectionString)
	}

}
