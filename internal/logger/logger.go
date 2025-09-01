package logger

import (
	"log/slog"
	"os"

	"github.com/vin-rmdn/http3-gui-client-go/internal/config"
)

// Setup initializes and returns a structured logger based on the provided configuration.
//
// The function reads the log level from the given configuration and applies it to the logger.
// If the log level is invalid or cannot be unmarshaled, the function logs an error and
// terminates the application with an appropriate exit code.
//
// Parameters:
//   - config: A Configuration object containing the logging settings.
//
// Returns:
//   - *slog.Logger: A configured logger instance ready for use.
//
// Note:
//   - The function exits the application with code 65 if the log level is invalid.
func Setup(config config.Configuration) (*slog.Logger, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(config.Log.Level)); err != nil {
		return nil, err
	}

	log := slog.New(slog.NewTextHandler(
		os.Stdout, &slog.HandlerOptions{Level: level}),
	)

	return log, nil
}
