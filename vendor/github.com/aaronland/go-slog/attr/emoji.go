package attr

import (
	"fmt"
	"log/slog"
)

type EmojiLevel struct {
	Emoji string
	Label string
}

type EmojiLevelMap map[slog.Level]*EmojiLevel

func DefaultEmojiLevelMap() EmojiLevelMap {

	emoji_map := map[slog.Level]*EmojiLevel{
		LevelTrace: &EmojiLevel{
			Emoji: "⁉️ ",
			Label: "TRACE",
		},
		LevelDebug: &EmojiLevel{
			Emoji: "🔍",
			Label: "DEBUG",
		},
		LevelInfo: &EmojiLevel{
			Emoji: "💬",
			Label: "INFO",
		},
		LevelWarning: &EmojiLevel{
			Emoji: "⚠️ ",
			Label: "WARNING",
		},
		LevelError: &EmojiLevel{
			Emoji: "🔥",
			Label: "ERROR",
		},
		LevelEmergency: &EmojiLevel{
			Emoji: "💥",
			Label: "EMERGENCY",
		},
	}

	return emoji_map
}

func EmojiLevelFunc() func(groups []string, a slog.Attr) slog.Attr {
	
	emoji_map := DefaultEmojiLevelMap()
	return EmojiLevelFuncWithMap(emoji_map)
}

func EmojiLevelFuncWithMap(emoji_map EmojiLevelMap) func(groups []string, a slog.Attr) slog.Attr {

	fn := func(groups []string, a slog.Attr) slog.Attr {

		if a.Key == slog.LevelKey {

			// Handle custom level values.
			level := a.Value.Any().(slog.Level)

			var emoji_level *EmojiLevel
			var match bool

			switch {
			case level < LevelDebug:
				emoji_level, match = emoji_map[LevelTrace]
			case level < LevelInfo:
				emoji_level, match = emoji_map[LevelDebug]
			case level < LevelNotice:
				emoji_level, match = emoji_map[LevelInfo]
			case level < LevelWarning:
				emoji_level, match = emoji_map[LevelWarning]
			case level < LevelError:
				emoji_level, match = emoji_map[LevelError]
			case level < LevelEmergency:
				emoji_level, match = emoji_map[LevelEmergency]
			default:
				emoji_level, match = emoji_map[LevelEmergency]
			}

			if match {
				a.Value = slog.StringValue(fmt.Sprintf("%s %s", emoji_level.Label, emoji_level.Emoji))
			}
		}

		return a

	}

	return fn
}
