package main

import (
	"github.com/bwmarrin/discordgo"
)

func dcTrains(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🚅 I like trains! 🚅",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
