package config

type AppConfig struct {
	Name        string
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

var App AppConfig
