package logging

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// DiscordgoLogger returns a logger to use with github.com/bwmarrin/discordgo
func DiscordgoLogger(logger *zap.Logger) func(msgL, caller int, format string, a ...interface{}) {

	return func(msgL, caller int, format string, a ...interface{}) {
		pc, file, line, _ := runtime.Caller(caller)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		l := logger.With(
			zap.String("file", fmt.Sprintf("%s:%d", file, line)),
			zap.String("method", name),
		)

		switch msgL {
		case discordgo.LogError:
			l.Error(
				fmt.Sprintf(format, a...),
			)

			return
		case discordgo.LogWarning:
			l.Warn(
				fmt.Sprintf(format, a...),
			)

			return
		case discordgo.LogInformational:
			l.Info(
				fmt.Sprintf(format, a...),
			)

			return
		case discordgo.LogDebug:
			l.Debug(
				fmt.Sprintf(format, a...),
			)

			return

		}

		l.Info(
			fmt.Sprintf(format, a...),
			zap.Int("log_level", msgL),
		)
	}
}
