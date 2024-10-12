package application

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/BernardN38/ebuy-server/authentication-service/messaging"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/handler"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/router"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/service"

	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
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
	conn, err := amqp.Dial(config.RabbitmqURL)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	err = initExchangesAndQueues(conn)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	rabbitmqEmitter := messaging.New(channel)
	//start service layer
	service := service.NewService(db, rabbitmqEmitter)

	//create request handler
	jwtAuth := jwtauth.New("HS512", []byte(config.JwtSecret), nil)
	hanlder := handler.NewHandler(service, jwtAuth)

	//create request router
	router := router.NewRouter(hanlder)
	return &App{
		Router: router,
	}
}

func (a *App) Run() error {
	log.Println("listening on port 8080")
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

type ExchangeQueueDeclaration struct {
	exchangeName string
	exchangeType string
	queueName    string
	routingKey   string
}

func initExchangesAndQueues(conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	declarations := []ExchangeQueueDeclaration{
		{
			exchangeName: "user_events",
			exchangeType: "topic",
			queueName:    "user_updates",
			routingKey:   "user.#",
		},
	}
	for _, v := range declarations {
		err := messaging.DeclareExchangeAndQueue(channel, v.exchangeName, v.exchangeType, v.queueName, v.routingKey)
		if err != nil {
			return err
		}
	}
	return nil
}
