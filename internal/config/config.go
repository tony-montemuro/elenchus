package config

import (
	"flag"
	"fmt"
	"log/slog"
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
	config.parseFlags()
	return config
}

func (c *Config) parseFlags() {
	c.Addr = flag.String("addr", ":4000", "HTTP network address")

	c.MinLogLevel = slog.LevelInfo
	flag.Func("minLogLevel", "Minimum logging level (see slog.Level; default \"INFO\")", func(level string) error {
		levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}

		for _, l := range levels {
			if strings.ToUpper(level) == l.String() {
				c.MinLogLevel = l
				return nil
			}
		}

		return fmt.Errorf("invalid log level: %s", level)
	})

	flag.Parse()
}
