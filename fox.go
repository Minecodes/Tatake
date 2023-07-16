package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func dcFox(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			Flags:   discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "Random Fox",
						URL:     "https://randomfox.ca",
						IconURL: "https://randomfox.ca/favicon.ico",
					},
					Color: 0xf48b00,
					Image: &discordgo.MessageEmbedImage{
						URL: "https://randomfox.ca/images/" + fmt.Sprint(rand.Intn(123)) + ".jpg",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		},
	})
}
