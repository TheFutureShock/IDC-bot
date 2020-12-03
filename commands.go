package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const BrandedColor = 16747590
const reportChannel = "783212613915770891"

func OnMsg(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if antispam.Increase(msg.Author.ID, msg.ID) {
		s.ChannelMessagesBulkDelete(msg.ChannelID, antispam.Users[msg.Author.ID].M)
	}

	if !strings.HasPrefix(msg.Content, prefix) {
		return
	}

	args := strings.Split(msg.Content, " ")

	command := strings.TrimPrefix(args[0], prefix)

	if command == "status" {
		p := getNukePredictor(msg.GuildID)
		var embed *discordgo.MessageEmbed

		if p.Triggered {
			strikes := strings.Join(p.Strikes, "\n")
			embed = &discordgo.MessageEmbed{
				Title: "Server Status : Lockdown",
				Color: BrandedColor,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Strikes: ",
						Value: fmt.Sprintf("```diff\n\n%s```", strikes),
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Written by [REDACTED]#4242",
					IconURL: "https://cdn.discordapp.com/attachments/781803550598627341/782155510429909012/pfp2.png",
				},
			}
		} else {
			embed = &discordgo.MessageEmbed{
				Title: "Server Status : Normal",
				Color: BrandedColor,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:   "Suspicion Level: ",
						Value:  strconv.Itoa(p.SuspicionLevel),
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Written by [REDACTED]#4242",
					IconURL: "https://cdn.discordapp.com/attachments/781803550598627341/782155510429909012/pfp2.png",
				},
			}
		}
		s.ChannelMessageSendEmbed(msg.ChannelID, embed)
	}

	if command == "restore" {
		p := getNukePredictor(msg.GuildID)
		if !p.Triggered {
			s.ChannelMessageSend(msg.ChannelID, "Server not in lockdown.")
			return
		}
		strikes := strings.Join(p.Strikes, "\n")

		embed := &discordgo.MessageEmbed{
			Title: "Order Restored",
			Color: BrandedColor,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Strikes from previous raid: ",
					Value: fmt.Sprintf("```diff\n\n%s```", strikes),
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Written by [REDACTED]#4242",
				IconURL: "https://cdn.discordapp.com/attachments/781803550598627341/782155510429909012/pfp2.png",
			},
		}

		go p.Restore()

		s.ChannelMessageSendEmbed(msg.ChannelID, embed)
	}

	if command == "help" {
		embed := &discordgo.MessageEmbed{
			Title: "Help",
			Color: BrandedColor,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Status: ",
					Value: "shows information about the server's current situation",
				},
				&discordgo.MessageEmbedField{
					Name:  "Restore: ",
					Value: "restores channels and roles deleted in the raid",
				},
				&discordgo.MessageEmbedField{
					Name:  "Report: ",
					Value: "reports a situation to the official IDC discord server\n Usage:  report (situation)",
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Written by [REDACTED]#4242",
				IconURL: "https://cdn.discordapp.com/attachments/781803550598627341/782155510429909012/pfp2.png",
			},
		}

		s.ChannelMessageSendEmbed(msg.ChannelID, embed)
	}

	if command == "report" {
		client.ChannelMessageSend(msg.ChannelID, "report sent")
		guild, _ := client.Guild(msg.GuildID)
		invite, _ := client.ChannelInviteCreate(msg.ChannelID, discordgo.Invite{})
		embed := &discordgo.MessageEmbed{
			Title: "Report Recieved",
			Color: BrandedColor,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Server: ",
					Value: guild.Name,
				},
				&discordgo.MessageEmbedField{
					Name:  "Invite: ",
					Value: fmt.Sprintf("https://discord.gg/%s", invite.Code),
				},
				&discordgo.MessageEmbedField{
					Name:  "Situation: ",
					Value: strings.Join(args[1:], " "),
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Written by [REDACTED]#4242",
				IconURL: "https://cdn.discordapp.com/attachments/781803550598627341/782155510429909012/pfp2.png",
			},
		}

		s.ChannelMessageSendEmbed(reportChannel, embed)
	}

	if command == "sponsors" || command == "credits" {
		embed := &discordgo.MessageEmbed{
			Title: "Sponsors",
			Description: `
			The following people helped to support the creation and maintenance of the bot:
			* ...
			* ...
			* ...
			`,
			Color: BrandedColor,
			Image: &discordgo.MessageEmbedImage{
				URL: "https://cdn.discordapp.com/attachments/781803550598627341/783238686145904660/thanks.png",
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Written by [REDACTED]#4242",
				IconURL: "https://cdn.discordapp.com/attachments/781803550598627341/782155510429909012/pfp2.png",
			},
		}

		s.ChannelMessageSendEmbed(msg.ChannelID, embed)
	}

}
