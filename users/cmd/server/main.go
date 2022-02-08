package main

import (
	"flag"
	"fmt"
	"go-grpc-demo/users/api/server"
	"go-grpc-demo/users/db"
	"os"

	"github.com/golang/glog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitialMigration for project with db.AutoMigrate
func InitialMigration() *gorm.DB {
	pgHost := "localhost"
	if len(os.Getenv("POSTGRES_USER")) > 0 {
		// see docker-compose.yml
		pgHost = "database"
	}
	pgUser := os.Getenv("POSTGRES_USER")
	if len(pgUser) == 0 {
		pgUser = "vs"
	}
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if len(pgPassword) == 0 {
		pgPassword = "111"
	}
	pgDb := os.Getenv("POSTGRES_DB")
	if len(pgDb) == 0 {
		pgDb = "vs"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", pgHost, pgUser, pgPassword, pgDb)
	con, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	// Migrate the schema
	_ = con.AutoMigrate(&db.User{})
	return con
}
func main() {
	DB := InitialMigration()
	srv, err := server.NewServer(DB)
	if err != nil {
		fmt.Println("Could not create server", err)
	}

	flag.Parse()
	defer glog.Flush()

	if err := srv.Serve(); err != nil {
		glog.Fatal(err)
	}

}
