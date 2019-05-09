package regexp

import (
	"regexp"
)

var (
	// DiscordInviteRegexp matches a linked Discord Invite
	DiscordInviteRegexp = regexp.MustCompile(
		`^(http(s)?:\/\/)?(discord\.gg(\/invite)?|discordapp\.com\/invite)\/([A-Za-z0-9-]+)$`,
	)
)
