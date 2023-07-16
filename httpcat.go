package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func dcHTTPCat(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			Flags:   discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "HTTP Cat",
					Color: 0xfa7c91,
					Image: &discordgo.MessageEmbedImage{
						URL:    "https://http.cat/" + fmt.Sprint(i.Interaction.ApplicationCommandData().Options[0].FloatValue()) + ".jpg",
						Width:  1400,
						Height: 1600,
					},
					Provider: &discordgo.MessageEmbedProvider{
						Name: "HTTP Cat",
						URL:  "https://http.cat",
					},
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "HTTP Cat",
						URL:     "https://http.cat",
						IconURL: "https://http.cat/favicon.ico",
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		},
	})
}
