package cmds

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func dcHTTPDog(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			Flags:   discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "HTTP Dog",
					Color: 0xfa7c91,
					Image: &discordgo.MessageEmbedImage{
						URL:    "https://http.dog/" + fmt.Sprint(i.Interaction.ApplicationCommandData().Options[0].FloatValue()) + ".jpg",
						Width:  1400,
						Height: 1600,
					},
					Provider: &discordgo.MessageEmbedProvider{
						Name: "HTTP Dog",
						URL:  "https://http.dog",
					},
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "HTTP Dog",
						URL:     "https://http.dog",
						IconURL: "https://http.dog/favicon.ico",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		},
	})
}
