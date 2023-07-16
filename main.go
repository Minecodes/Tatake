package main

import (
	//"errors"

	"flag"
	"fmt"
	"math/rand"

	//"fmt"

	"log"
	"os"
	"os/signal"
	"sync"

	//"strings"
	//"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
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
		{
			Name:        "password",
			Description: "Generate a random password",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "length",
					Description: "Length of the password",
					Required:    true,
					MinValue:    &minPasswordLength,
					MaxValue:    maxPasswordLength,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"trains":   dcTrains,
		"ping":     dcPing,
		"weather":  dcWeather,
		"gh":       dcGhUser,
		"httpdog":  dcHTTPDog,
		"httpcat":  dcHTTPCat,
		"fox":      dcFox,
		"qrcode":   dcQRCode,
		"password": dcPassword,
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

func handleMessage(s xmpp.Sender, p stanza.Packet) {
	msg, ok := p.(stanza.Message)
	if !ok {
		_, _ = fmt.Fprintf(os.Stdout, "Ignoring packet: %T\n", p)
		return
	}

	_, _ = fmt.Fprintf(os.Stdout, "Body = %s - from = %s\n", msg.Body, msg.From)
	reply := stanza.Message{Attrs: stanza.Attrs{To: msg.From}, Body: msg.Body}
	_ = s.Send(reply)
}

func errorHandler(err error) {
	fmt.Println(err.Error())
}

func dc(wg *sync.WaitGroup) {
	defer wg.Done()
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
	os.Exit(0)
}

func xmppBot(wg *sync.WaitGroup) {
	defer wg.Done()

	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: "etc.minecodes.de:5222",
		},
		Jid:          "",
		Credential:   xmpp.Password(""),
		StreamLogger: os.Stdout,
		Insecure:     false,
		// TLSConfig: tls.Config{InsecureSkipVerify: true},
	}

	router := xmpp.NewRouter()
	router.HandleFunc("message", handleMessage)

	client, err := xmpp.NewClient(&config, router, errorHandler)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// If you pass the client to a connection manager, it will handle the reconnect policy
	// for you automatically.
	cm := xmpp.NewStreamManager(client, nil)
	log.Fatal(cm.Run())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	cm.Stop()
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	go dc(&wg)
	go xmppBot(&wg)

	wg.Wait()

	/**config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: "etc.minecodes.de:5222",
		},
		Jid:          "tatake@etc.minecodes.de",
		Credential:   xmpp.Password("T*AM5%843#H&w!6krd&!"),
		StreamLogger: os.Stdout,
		Insecure:     false,
		// TLSConfig: tls.Config{InsecureSkipVerify: true},
	}

	router := xmpp.NewRouter()
	router.HandleFunc("message", handleMessage)

	client, err := xmpp.NewClient(&config, router, errorHandler)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// If you pass the client to a connection manager, it will handle the reconnect policy
	// for you automatically.
	cm := xmpp.NewStreamManager(client, nil)
	log.Fatal(cm.Run())**/
}
