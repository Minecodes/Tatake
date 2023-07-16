package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

func dcWeather(s *discordgo.Session, i *discordgo.InteractionCreate) {
	loc := i.Interaction.ApplicationCommandData().Options[0].StringValue()
	// make http request
	req, err := http.NewRequest("GET", "https://wttr.in/"+loc+"?0&T", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "curl")
	req.Response, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer req.Response.Body.Close()
	body, err := ioutil.ReadAll(req.Response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if i.Interaction.ApplicationCommandData().Options[1].BoolValue() {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "",
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Weather in " + loc,
						Description: "```" + string(body) + "```",
						Color:       0x0E86D4,
						Timestamp:   time.Now().Format(time.RFC3339),
						Provider: &discordgo.MessageEmbedProvider{
							Name: "wttr.in",
							URL:  "https://wttr.in",
						},
						Author: &discordgo.MessageEmbedAuthor{
							Name: "wttr.in",
							URL:  "https://wttr.in",
						},
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "",
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Weather in " + loc,
						Description: "```" + string(body) + "```",
						Color:       0x0E86D4,
						Timestamp:   time.Now().Format(time.RFC3339),
						Provider: &discordgo.MessageEmbedProvider{
							Name: "wttr.in",
							URL:  "https://wttr.in",
						},
						Author: &discordgo.MessageEmbedAuthor{
							Name: "wttr.in",
							URL:  "https://wttr.in",
						},
					},
				},
			},
		})
	}
}
