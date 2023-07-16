package cmds

import (
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	minPasswordLength = float64(12)
	maxPasswordLength = float64(100)
	passwordChars     = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+{}[]:;?/.,<>")
)

func gen(length int64) string {
	password := make([]rune, length)
	for i := range password {
		password[i] = passwordChars[rand.Intn(len(passwordChars))]
	}
	// check if password has at least one number, one uppercase letter, one lowercase letter and one special character
	// if not, generate a new password
	if !strings.ContainsAny(string(password), "0123456789") || !strings.ContainsAny(string(password), "abcdefghijklmnopqrstuvwxyz") || !strings.ContainsAny(string(password), "ABCDEFGHIJKLMNOPQRSTUVWXYZ") || !strings.ContainsAny(string(password), "!@#$%^&*()_+{}[]:;?/.,<>") {
		return gen(length)
	}
	return string(password)
}

func dcPassword(s *discordgo.Session, i *discordgo.InteractionCreate) {
	length := i.Interaction.ApplicationCommandData().Options[0].IntValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			Flags:   discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Password",
					Color:       0x00ff00,
					Description: "Your password is: `" + gen(length) + "`",
				},
			},
		},
	})
}
