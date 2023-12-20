package monitor

import (
	"log"

	"encore.dev/storage/sqldb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	db *gorm.DB
}

// var secrets struct {
// 	user     string
// 	password string
// }

var checksDB = sqldb.NewDatabase("checks", sqldb.DatabaseConfig{
	Migrations: "./migrations",
}).Stdlib()

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: checksDB,
	}))
	if err != nil {
		log.Println("Somethings Wrong due to, ", err.Error())
		return nil, err
	}
	log.Printf("Database Connected to %v\n", db.Migrator().CurrentDatabase())
	return &Service{db: db}, nil
}
