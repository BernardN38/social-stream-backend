package application

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/handler"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/router"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/service"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type App struct {
	Router *router.Router
}

func NewApp() *App {
	config, err := loadEnvConfig()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	// Connect to the database
	db, err := sql.Open("postgres", config.PostgresDsn)
	if err != nil {
		log.Fatalln("unable to connect to the database:", err)
		return nil
	}
	defer db.Close()
	// Check if the database exists
	if err := createDatabaseIfNotExists(db, config.DbName); err != nil {
		log.Fatalln("unable to create or check the database:", err)
		return nil
	}

	// Connect to the specific database
	db, err = sql.Open("postgres", config.PostgresDsn+" dbname="+config.DbName)
	if err != nil {
		log.Fatalln("unable to connect to the specific database:", err)
		return nil
	}
	// Run database migrations
	if err := RunDatabaseMigrations(db); err != nil {
		log.Fatalln("unable to run database migrations:", err)
		return nil
	}

	//start service layer
	service := service.NewService(db)

	//create request handler
	hanlder := handler.NewHandler(service)

	//create request router
	router := router.NewRouter(hanlder)
	return &App{
		Router: router,
	}
}

func (a *App) Run() error {
	return http.ListenAndServe(":8080", a.Router.R)
}

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB(connStr string) (*sql.DB, error) {
	// Open the connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure connection is established
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to PostgreSQL database!")
	return db, nil
}

func createDatabaseIfNotExists(db *sql.DB, dbName string) error {
	result, err := db.Exec(fmt.Sprintf("select 1 from pg_database where datname = '%s'", dbName))
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return err
		}
	}

	return err
}
