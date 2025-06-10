package config

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/tony-montemuro/elenchus/internal/assert"
)

func TestLoadConfig(t *testing.T) {
	defaultAddr := ":4000"
	customAddr := ":9000"

	defaultMinLogLevel := slog.LevelInfo
	customMinLogLevel := slog.LevelDebug

	defaultDsn := "web:pass@/elenchus?parseTime=true"
	customDsn := "migration:foobar@/elenchus?parseTime=true"

	tests := []struct {
		name        string
		args        []string
		want        Config
		expectError bool
	}{
		{
			name:        "No args",
			args:        []string{},
			want:        Config{Addr: &defaultAddr, Dsn: &defaultDsn, MinLogLevel: defaultMinLogLevel},
			expectError: false,
		},
		{
			name:        "All args specificed",
			args:        []string{fmt.Sprintf("-addr=%s", customAddr), fmt.Sprintf("-dsn=%s", customDsn), fmt.Sprintf("-minLogLevel=%s", customMinLogLevel.String())},
			want:        Config{Addr: &customAddr, Dsn: &customDsn, MinLogLevel: customMinLogLevel},
			expectError: false,
		},
		{
			name:        "Specified addr arg",
			args:        []string{fmt.Sprintf("-addr=%s", customAddr)},
			want:        Config{Addr: &customAddr, Dsn: &defaultDsn, MinLogLevel: defaultMinLogLevel},
			expectError: false,
		},
		{
			name:        "Specified minLogLevel arg",
			args:        []string{fmt.Sprintf("-minLogLevel=%s", customMinLogLevel.String())},
			want:        Config{Addr: &defaultAddr, MinLogLevel: customMinLogLevel},
			expectError: false,
		},
		{
			name:        "Specified dsn arg",
			args:        []string{fmt.Sprintf("-dsn=%s", customDsn)},
			want:        Config{Addr: &defaultAddr, Dsn: &customDsn, MinLogLevel: defaultMinLogLevel},
			expectError: false,
		},
		{
			name:        "Invalid minLogLevel arg",
			args:        []string{"-minLogLevel=INVALID"},
			want:        Config{},
			expectError: true,
		},
		{
			name:        "Case insensitive minLogLevel arg",
			args:        []string{fmt.Sprintf("-minLogLevel=%s", strings.ToLower(customMinLogLevel.String()))},
			want:        Config{Addr: &defaultAddr, MinLogLevel: customMinLogLevel},
			expectError: false,
		},
		{
			name:        "Multiple invalid args",
			args:        []string{"-minLogLevel=INVALID", "-addr=INVALID"},
			want:        Config{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := flag.NewFlagSet(tt.name, flag.ContinueOnError)
			flags.SetOutput(io.Discard)

			config := &Config{}
			config.parseFlags(flags)
			err := flags.Parse(tt.args)

			if err != nil {
				if !tt.expectError {
					t.Errorf("problem parsing flags: %s", err.Error())
				}
			} else {
				if tt.expectError {
					t.Error("expected error but got none")
				}

				assert.Equal(t, *config.Addr, *tt.want.Addr)
				assert.Equal(t, config.MinLogLevel, tt.want.MinLogLevel)
			}

		})
	}
}
