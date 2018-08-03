package cfg

import (
	"os"

	"github.com/sknv/microrabbit/app/lib/xflags"
)

type Config struct {
	RabbitURL string `long:"rabbit-url" env:"RABBIT_URL" default:"amqp://guest:guest@localhost:5672" description:"rabbitmq url"`
}

func Parse() *Config {
	cfg := new(Config)
	if _, err := xflags.ParseArgs(os.Args[1:], cfg); err != nil {
		os.Exit(1)
	}
	return cfg
}
