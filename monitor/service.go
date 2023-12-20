package monitor

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//encore:service
type Service struct {
	db *gorm.DB
}

var secrets struct {
	user     string
	password string
}

// initService initializes the site service.
// It is automatically called by Encore on service startup.
func initService() (*Service, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", "localhost", string(secrets.user), string(secrets.password), "checks", "5432", "disable", "Asia/Shanghai")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Somethings Wrong due to, ", err.Error())
		return nil, err
	}

	log.Printf("Database Connected to %v\n", db.Migrator().CurrentDatabase())
	return &Service{db: db}, nil
}
