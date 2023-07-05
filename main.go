package main

import (
	//"errors"

	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	//"strings"
	//"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	InfoMessages   = []string{
		`Hello, I'm a bot made by <@!556848982433857537>!`,
		`Hello SlimeDiamond`,
		`"Never trust a tech guy with a rat tail—too easy to carve secrets out of him." - Lone Star (Mr. Robot)`,
		`me gusta los trenes`,
		`Now open source on GitHub, git.minecodes.de and CodeBerg`,
		`Not Windows created or running, just Debian and Arch`,
		`Schnitzel mit Spätzle ist auch lercker`,
		`Does anyone also hear the Doom music getting louder?`,
		`HAM radio is cool :D`,
		`No, I'm not proprietary, I'm open source`,
		`I was coded in NodeJS, but now I'm coded in Go`,
		`If you need a job: created Linux VM, then "sudo rm -rf /*" and bye bye VM`,
		`Encryption should be a human right`,
	}
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	err = godotenv.Load()
	if BotToken == nil || *BotToken == "" {
		*BotToken = os.Getenv("TOKEN")
	}
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "trains",
			Description: "I like trains",
		},
		{
			Name:        "ping",
			Description: "Ping the bot",
		},
		{
			Name:        "weather",
			Description: "Get the weather",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "city",
					Description: "Name of the city",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "private",
					Description: "Send the weather as a message that only you can see",
					Required:    true,
				},
			},
		},
		{
			Name:        "gh",
			Description: "Github commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "user",
					Description: "Get the user info",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "user",
							Description: "Name of the user",
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "httpdog",
			Description: "HTTP status codes to dog pictures",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "code",
					Description: "HTTP status code",
					Required:    true,
					MinValue:    &integerOptionMinValue,
					MaxValue:    599,
				},
			},
		},
		{
			Name:        "httpcat",
			Description: "HTTP status codes to dog pictures",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "code",
					Description: "HTTP status code",
					Required:    true,
					MinValue:    &integerOptionMinValue,
					MaxValue:    599,
				},
			},
		},
		{
			Name:        "fox",
			Description: "Get a random fox picture",
		},
		{
			Name:        "qrcode",
			Description: "Generate a QR code",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "Text to encode",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"trains": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "🚅 I like trains! 🚅",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
		"weather": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"gh": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"httpdog": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"httpcat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"fox": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"qrcode": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			value := i.Interaction.ApplicationCommandData().Options[0].StringValue()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "",
					Flags:   discordgo.MessageFlagsEphemeral,
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "QR Code",
							Color: 0xffff66,
							Image: &discordgo.MessageEmbedImage{
								URL:    "https://api.qrserver.com/v1/create-qr-code/?size=1000x1000&data=" + value,
								Width:  1000,
								Height: 1000,
							},
						},
					},
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.UpdateGameStatus(0, InfoMessages[rand.Intn(len(InfoMessages))])
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	if err != nil {
		log.Fatalf("Cannot get guilds: %v", err)
	}
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	s.UpdateGameStatus(0, InfoMessages[rand.Intn(len(InfoMessages))])

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
