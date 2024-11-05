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

	// test database
	TestHost     string `mapstructure:"test_host"`
	TestPort     int    `mapstructure:"test_port"`
	TestUsername string `mapstructure:"test_username"`
	TestPassword string `mapstructure:"test_password"`
	TestName     string `mapstructure:"test_name"`

	// parameters
	MaxOpenConns int `mapstructure:"max_open_connections"`
	MaxIdleConns int `mapstructure:"max_idle_connections"`
}

var App AppConfig
