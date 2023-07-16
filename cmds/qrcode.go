package cmds

import (
	"github.com/bwmarrin/discordgo"
)

func dcQRCode(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
}
