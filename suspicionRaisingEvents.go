package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func channelDeleted(s *discordgo.Session, c *discordgo.ChannelDelete) {
	p := getNukePredictor(c.GuildID)

	p.IncreaseSuspicionLevel(25)

	p.RestorableChannels = append(p.RestorableChannels, c.Channel)

	p.Strikes = append(p.Strikes, fmt.Sprintf("- Deleted Channel: %s", c.Name))
}

func memberBanned(s *discordgo.Session, m *discordgo.GuildBanAdd) {
	p := getNukePredictor(m.GuildID)

	p.IncreaseSuspicionLevel(20)

	p.Strikes = append(p.Strikes, fmt.Sprintf("- Banned User: %s", m.User.Username))
}

func memberKicked(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	p := getNukePredictor(m.GuildID)

	p.IncreaseSuspicionLevel(20)

	p.Strikes = append(p.Strikes, fmt.Sprintf("- Kicked User: %s", m.User.Username))
}
