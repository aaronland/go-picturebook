package attr

import (
	"log/slog"
)

func ReplaceAttrChain(chain ...func(groups []string, a slog.Attr) slog.Attr) func(groups []string, a slog.Attr) slog.Attr {

	fn := func(groups []string, a slog.Attr) slog.Attr {

		for _, fn := range chain {
			a = fn(groups, a)
		}

		return a
	}

	return fn
}
