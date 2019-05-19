package discord

import (
	"github.com/bwmarrin/discordgo"
)

func SetAPIBase(apiBase string) {
	if apiBase == "" {
		return
	}

	discordgo.EndpointDiscord = apiBase
	discordgo.EndpointAPI = discordgo.EndpointDiscord + "api/v" + discordgo.APIVersion + "/"
	discordgo.EndpointGuilds = discordgo.EndpointAPI + "guilds/"
	discordgo.EndpointChannels = discordgo.EndpointAPI + "channels/"
	discordgo.EndpointUsers = discordgo.EndpointAPI + "users/"
	discordgo.EndpointGateway = discordgo.EndpointAPI + "gateway"
	discordgo.EndpointGatewayBot = discordgo.EndpointGateway + "/bot"
	discordgo.EndpointWebhooks = discordgo.EndpointAPI + "webhooks/"
}
