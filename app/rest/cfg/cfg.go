package cfg

import (
	"os"

	"github.com/sknv/microrabbit/app/lib/xflags"
)

type Config struct {
	Addr string `long:"rest-addr" env:"REST_ADDR" default:"localhost:8080" description:"rest api address"`
}

func Parse() *Config {
	var cfg Config
	if _, err := xflags.ParseArgs(os.Args[1:], &cfg); err != nil {
		os.Exit(1)
	}
	return &cfg
}
