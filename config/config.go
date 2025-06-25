package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     	string
	Port     	string
	Username 	string
	Password 	string
	Dbname		string
	Driver		string
}

type ApiConfig struct {
	ApiPort		string
}

type Config struct {
	DBConfig
	ApiConfig
}

func (c *Config) ReadConfig() error {
	_ = godotenv.Load()

	c.DBConfig = DBConfig{
		Host:	os.Getenv("DB_HOST"),
		Port:	os.Getenv("DB_PORT"),
		Username:	os.Getenv("DB_USER"),
		Password:	os.Getenv("DB_PASSWORD"),
		Dbname:	os.Getenv("DB_NAME"),
	}

	c.ApiConfig = ApiConfig{
		ApiPort:	os.Getenv("API_PORT"),
	}

	if c.Host == "" || c.Port == "" || c.Username == "" || c.Password == "" || c.Dbname == "" {
		return errors.New("Some config is empty")
	}

	return  nil
}

func NewConfig() (*Config, error) {
	config := &Config{}

	if err := config.ReadConfig(); err != nil {
		return nil, err
	}

	return  config, nil
}