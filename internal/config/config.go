package config

import "fmt"

type Config struct {
	HTTPPort int            `envconfig:"http_port" default:"8000"`
	HTTPMode HTTPModeConfig `envconfig:"http_mode"`

	WebsocketPort int `envconfig:"websocket_port" default:"8001"`

	LogToConsole bool   `envconfig:"log_to_console" default:"true"`
	LogFile      string `envconfig:"log_file"`
}

type HTTPModeConfig string

const (
	HTTPModeRelease = "release"
	HTTPModeDebug   = "debug"
)

func (c *HTTPModeConfig) Decode(value string) error {
	if value == "" {
		value = HTTPModeDebug
	}
	if value != HTTPModeRelease && value != HTTPModeDebug {
		return fmt.Errorf("invalid http mode: %s", value)
	}
	*c = HTTPModeConfig(value)
	return nil
}
