package config

type Config struct {
	HTTPPort     int    `envconfig:"http_port" default:"8000"`
	LogToConsole bool   `envconfig:"log_to_console" default:"true"`
	LogFile      string `envconfig:"log_file"`
}
