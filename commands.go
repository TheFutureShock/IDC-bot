package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const BrandedColor = 16747590
const reportChannel = "770899887260303390"

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
		p := nukePredictors[msg.GuildID]
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
		guild, _ := client.Guild(msg.GuildID)
		if msg.Author.ID != guild.OwnerID {
			return
		}

		p := nukePredictors[msg.GuildID]

		if !p.Triggered {
			s.ChannelMessageSend(msg.ChannelID, "Server not in lockdown.")
			return
		}

		if p.Restoring {
			s.ChannelMessageSend(msg.ChannelID, "Server already attempting to restore .")
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

	// permissions

	if command == "permissionset" {

		var permissions PermissionEntree
		PermissionsRetrieve.Get(&permissions, msg.Author.ID)

		if permissions.EditPermissions != 1 && msg.Author.ID != "628298193922424857" { // yes I hardcoded my ID in I'll update the env eventually
			client.ChannelMessageSend(msg.ChannelID, "Missing perms")
			return
		}

		permWatchlist, err := strconv.Atoi(args[2])
		if err != nil {
			client.ChannelMessageSend(msg.ChannelID, "Invalid Parameters")
		}
		permPerm, err := strconv.Atoi(args[3])
		if err != nil {
			client.ChannelMessageSend(msg.ChannelID, "Invalid Parameters")
		}

		PermissionsSet.Exec(args[1], permWatchlist, permPerm)
	}

	// watchlist commands

	if command == "watchlist" {
		if args[1] == "add" {

			var permissions PermissionEntree
			err := PermissionsRetrieve.Get(&permissions, msg.Author.ID)
			fmt.Println(err)

			if permissions.WatchlistAdmin != 1 {
				client.ChannelMessageSend(msg.ChannelID, "Missing perms")
				return
			}

			userID := args[2]
			reason := strings.Join(args[3:], " ")
			user, err := client.User(userID)
			if err != nil {
				client.ChannelMessageSend(msg.ChannelID, "invalid USER ID")
				return
			}

			if reason == "" {
				client.ChannelMessageSend(msg.ChannelID, "missing reason")
				return
			}

			g, _ := client.Guild(msg.GuildID)

			WatchlistADDREPORT.Exec(userID, fmt.Sprintf("%s#%s", user.Username, user.Discriminator), user.AvatarURL(""), reason, msg.GuildID, g.Name)

			row := WatchlistCOUNT.QueryRow()
			var count int
			row.Scan(&count)

			embed := &discordgo.MessageEmbed{
				Title:       "Suspect Added",
				Color:       BrandedColor,
				Description: fmt.Sprintf("%s#%s", user.Username, user.Discriminator),
				Image: &discordgo.MessageEmbedImage{
					URL: user.AvatarURL(""),
				},
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Users in watchlist: ",
						Value: strconv.Itoa(count),
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Written by [REDACTED]#4242",
					IconURL: "https://cdn.discordapp.com/attachments/781803550598627341/782155510429909012/pfp2.png",
				},
			}

			client.ChannelMessageSendEmbed(msg.ChannelID, embed)
		}

		if args[1] == "dump" {
			var list []WatchlistEntree
			WatchlistALL.Select(&list)
			var file string
			for _, entree := range list {
				file += fmt.Sprintf("%s - %s\n", entree.UserTag, entree.Reason)
			}
			client.ChannelFileSend(msg.ChannelID, "dump.txt", strings.NewReader(file))
		}

		if args[1] == "lookup" {
			if len(args) != 3 {
				client.ChannelMessageSend(msg.ChannelID, "invalid USER ID")
				return
			}
			userID := args[2]

			var list []WatchlistEntree
			WatchlistUSERREPORTS.Select(&list, userID)

			if len(list) < 1 {
				client.ChannelMessageSend(msg.ChannelID, "user is not on the list")
				return
			}

			r1 := list[0]

			var fields []*discordgo.MessageEmbedField

			for i, report := range list {
				fields = append(fields, &discordgo.MessageEmbedField{Name: fmt.Sprintf("report %d", i+1), Value: fmt.Sprintf("Reported in: %s\n Reason: %s", report.OriginGuildName, report.Reason)})
			}

			embed := &discordgo.MessageEmbed{
				Title:       "Suspect Found",
				Color:       BrandedColor,
				Description: fmt.Sprintf("%s, reported %d time(s)", r1.UserTag, len(list)),
				Image: &discordgo.MessageEmbedImage{
					URL: r1.UserPFP,
				},
				Fields: fields,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Written by [REDACTED]#4242",
					IconURL: "https://cdn.discordapp.com/attachments/781803550598627341/782155510429909012/pfp2.png",
				},
			}

			client.ChannelMessageSendEmbed(msg.ChannelID, embed)
		}

		if args[1] == "bansync" {
			guild, _ := client.Guild(msg.GuildID)
			if msg.Author.ID != guild.OwnerID {
				client.ChannelMessageSend(msg.ChannelID, "You must be the server owner to do this")
			}

			
			var list []WatchlistEntree
			WatchlistALL.Select(&list)

			go (func() {
				for _, entree := range list {
					client.GuildBanCreateWithReason(msg.GuildID, entree.UserID, entree.Reason, 0)
				}
			})()

			row := WatchlistCOUNT.QueryRow()
			var count int
			row.Scan(&count)

			client.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Banning %d users.", count))
		}

	}

}
