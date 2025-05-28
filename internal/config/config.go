package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Addr        *string
	MinLogLevel slog.Level
}

// This function should only be called once: at the start of the application
// See `flag.Parse()` for details on why
func LoadConfig() *Config {
	config := &Config{}

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config.parseFlags(flags)
	flags.Parse(os.Args[1:])

	return config
}

func (c *Config) parseFlags(flags *flag.FlagSet) {
	c.Addr = flags.String("addr", ":4000", "HTTP network address")

	c.MinLogLevel = slog.LevelInfo
	flags.Func("minLogLevel", "Minimum logging level (see slog.Level; default \"INFO\")", func(level string) error {
		levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}

		for _, l := range levels {
			if strings.ToUpper(level) == l.String() {
				c.MinLogLevel = l
				return nil
			}
		}

		return fmt.Errorf("invalid log level: %s", level)
	})

}
