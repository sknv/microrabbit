package cfg

import (
	"os"

	"github.com/sknv/microrabbit/app/lib/xflags"
)

type Config struct {
	RabbitAddr string `long:"rabbit-addr" env:"RABBIT_ADDR" default:"amqp://guest:guest@localhost:5672" description:"rabbitmq address"`
}

func Parse() *Config {
	cfg := new(Config)
	if _, err := xflags.ParseArgs(os.Args[1:], cfg); err != nil {
		os.Exit(1)
	}
	return cfg
}
