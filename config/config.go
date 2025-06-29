package config

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type TokenConfig struct {
	AppName				string
	JWTSignatureKey		[]byte
	JWTSigningMethod	*jwt.SigningMethodHMAC
	AccessTokenLifetime	time.Duration
}

type Config struct {
	DBConfig
	ApiConfig
	TokenConfig
}


func (c *Config) ReadConfig() error {
	_ = godotenv.Load()

	c.DBConfig = DBConfig{
		Host:	os.Getenv("DB_HOST"),
		Port:	os.Getenv("DB_PORT"),
		Dbname:	os.Getenv("DB_NAME"),
		Username:	os.Getenv("DB_USER"),
		Password:	os.Getenv("DB_PASSWORD"),
		Driver: os.Getenv("DB_DRIVER"),
	}

	c.ApiConfig = ApiConfig{
		ApiPort:	os.Getenv("API_PORT"),
	}

	accessTokenLifetime := time.Duration(15) * time.Minute

	c.TokenConfig = TokenConfig{
		AppName: "Pay-Share",
		JWTSignatureKey: []byte(os.Getenv("ACCESS_TOKEN")),
		JWTSigningMethod: jwt.SigningMethodHS256,
		AccessTokenLifetime: accessTokenLifetime,
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