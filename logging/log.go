package logging

import (
	"log/slog"
	"os"
)

func Setup(leveler slog.Leveler) {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: leveler},
			),
		),
	)
}

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}
