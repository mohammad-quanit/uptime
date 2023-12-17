package monitor

import (
	"fmt"

	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	db *gorm.DB
}

var ChecksDB = sqldb.NewDatabase("checks", sqldb.DatabaseConfig{
	Migrations: "./migrations",
}).Stdlib()

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: ChecksDB,
	}))
	if err != nil {
		fmt.Println("This is connection error.............", err.Error())
		return nil, err
	}
	return &Service{db: db}, nil
}
