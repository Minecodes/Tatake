package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func dcGhUser(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Interaction.ApplicationCommandData().Options[0].Name == "user" {
		user := i.Interaction.ApplicationCommandData().Options[0].Options[0].StringValue()
		req, err := http.NewRequest("GET", "https://api.github.com/users/"+user, nil)
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

		type User struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
			Name              string `json:"name"`
			Company           string `json:"company"`
			Blog              string `json:"blog"`
			Location          string `json:"location"`
			Email             string `json:"email"`
			Hireable          bool   `json:"hireable"`
			Bio               string `json:"bio"`
			TwitterUsername   string `json:"twitter_username"`
			PublicRepos       int    `json:"public_repos"`
			PublicGists       int    `json:"public_gists"`
			Followers         int    `json:"followers"`
			Following         int    `json:"following"`
			CreatedAt         string `json:"created_at"`
			UpdatedAt         string `json:"updated_at"`
		}

		var data User
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatal(err)
		}

		var fields = []*discordgo.MessageEmbedField{
			{
				Name:   "Username",
				Value:  data.Login,
				Inline: true,
			},
			{
				Name:   "ID",
				Value:  strconv.Itoa(data.ID),
				Inline: true,
			},
			{
				Name:   "Public Repos",
				Value:  strconv.Itoa(data.PublicRepos),
				Inline: true,
			},
			{
				Name:   "Public Gists",
				Value:  strconv.Itoa(data.PublicGists),
				Inline: true,
			},
			{
				Name:   "Followers",
				Value:  strconv.Itoa(data.Followers),
				Inline: true,
			},
			{
				Name:   "Following",
				Value:  strconv.Itoa(data.Following),
				Inline: true,
			},
			{
				Name:   "Created At",
				Value:  data.CreatedAt,
				Inline: true,
			},
		}

		if data.Bio != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Bio",
				Value:  data.Bio,
				Inline: true,
			})
		}
		if data.Company != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Company",
				Value:  data.Company,
				Inline: true,
			})
		}
		if data.Location != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Location",
				Value:  data.Location,
				Inline: true,
			})
		}
		if data.Email != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Email",
				Value:  data.Email,
				Inline: true,
			})
		}
		if data.Blog != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Blog",
				Value:  data.Blog,
				Inline: true,
			})
		}
		if data.TwitterUsername != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Twitter",
				Value:  "https://twitter.com/" + data.TwitterUsername,
				Inline: true,
			})
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "",
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Github: " + user,
						Color:       0x24292D,
						Description: "",
						Provider: &discordgo.MessageEmbedProvider{
							Name: "Github",
							URL:  "https://github.com",
						},
						Author: &discordgo.MessageEmbedAuthor{
							Name:    data.Name,
							URL:     data.HTMLURL,
							IconURL: data.AvatarURL,
						},
						Timestamp: time.Now().Format(time.RFC3339),
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Provided by Github",
							IconURL: "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png",
						},
						Fields: fields,
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	}
}
