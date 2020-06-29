package preferences

import (
	"github.com/kelseyhightower/envconfig"
)

type Preferences struct {
	LogLevel          int    `envconfig:"LOG_LEVEL"`
	LogAsJSON         bool   `envconfig:"LOG_AS_JSON"`
	MaxTimeoutInMsec  uint   `envconfig:"MAX_TIMEOUT_IN_MSEC" required:"true"`
	WriteTimeoutInSec int    `envconfig:"WRITE_TIMEOUT_IN_SEC" required:"true"`
	ReadTimeoutInSec  int    `envconfig:"READ_TIMEOUT_IN_SEC" required:"true"`
	HttpAddress       string `envconfig:"HTTP_ADDRESS" required:"true"`
}

func Get() (*Preferences, error) {
	var p Preferences
	err := envconfig.Process("", &p)
	return &p, err
}
