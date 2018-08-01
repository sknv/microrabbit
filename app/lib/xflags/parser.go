package xflags

import (
	"os"

	"github.com/jessevdk/go-flags"
)

func ParseArgs(args []string, dest interface{}) ([]string, error) {
	flagParser := flags.NewParser(dest, flags.Default)
	return flagParser.ParseArgs(os.Args[1:])
}
