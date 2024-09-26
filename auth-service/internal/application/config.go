package application

import (
	"os"

	"github.com/go-playground/validator/v10"
)

type config struct {
	Port        string `validate:"required"`
	PostgresDsn string `validate:"required"`
	JwtSecret   string `validate:"required"`
	DbName      string `validate:"required"`
}

func (c *config) Validate() error {
	validator := validator.New()
	err := validator.Struct(c)
	if err != nil {
		return err
	}
	return nil
}

func loadEnvConfig() (*config, error) {
	port := os.Getenv("port")
	jwtSecret := os.Getenv("jwtSecret")
	postgresDsn := os.Getenv("postgresDsn")
	dbName := os.Getenv("dbName")
	config := config{
		Port:        port,
		JwtSecret:   jwtSecret,
		PostgresDsn: postgresDsn,
		DbName:      dbName,
	}
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	return &config, config.Validate()
}
