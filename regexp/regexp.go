package regexp

import (
	"regexp"
)

var (
	// ContainsNumber matches on any numbers
	ContainsNumber = regexp.MustCompile(`[0-9]+`)

	// DiscordInviteCodeRegexp matches a Discord Invite Code
	DiscordInviteCodeRegexp = regexp.MustCompile(
		`^[A-Za-z0-9-]+$`,
	)

	// DiscordInviteRegexp matches a linked Discord Invite
	DiscordInviteRegexp = regexp.MustCompile(
		`^(http(s)?:\/\/)?(discord\.gg(\/invite)?|discordapp\.com\/invite)\/([A-Za-z0-9-]+)(\/?)$`,
	)

	// EmojiRegex matches a built in Emoji or discord emoji
	EmojiRegex = regexp.MustCompile(`[\x{00A0}-\x{1F9EF}]|<(a)?:[^<>:]+:[0-9]+>`)

	UnicodeEmojiRegex = regexp.MustCompile(`[\x{00A0}-\x{1F9EF}]`)

	// DiscordEmojiRegexp matches an Emoji, or an Emoji ID
	DiscordEmojiRegexp = regexp.MustCompile(`(<(a)?:([^<>:]+):)?([0-9]+)>?`)

	DiscordEmojiStrictRegexp = regexp.MustCompile(`<a?:(\w+):(\d+)>`)

	SnowflakeRegexp = regexp.MustCompile(`^[0-9]+$`)

	DiscordMessageLinkRegexp = regexp.MustCompile(`https:\/\/discord(app)?.com\/channels\/(\d+)\/(\d+)\/(\d+)`)
)
