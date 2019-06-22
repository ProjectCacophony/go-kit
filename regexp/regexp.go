package regexp

import (
	"regexp"
)

var (
	// DiscordInviteCodeRegexp matches a Discord Invite Code
	DiscordInviteCodeRegexp = regexp.MustCompile(
		`^[A-Za-z0-9-]+$`,
	)

	// DiscordInviteRegexp matches a linked Discord Invite
	DiscordInviteRegexp = regexp.MustCompile(
		`^(http(s)?:\/\/)?(discord\.gg(\/invite)?|discordapp\.com\/invite)\/([A-Za-z0-9-]+)(\/?)$`,
	)

	// DiscordEmojiRegexp matches an Emoji, or an Emoji ID
	DiscordEmojiRegexp = regexp.MustCompile(`(<(a)?:([^<>:]+):)?([0-9]+)>?`)
)
