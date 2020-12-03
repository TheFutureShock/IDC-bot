package nukeprediction

import (
	"fmt"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
)

const defaultPerms = 36785728

type NukePrediction struct {
	GuildID            string
	SuspicionLevel     int
	RestorableChannels []*discordgo.Channel
	RestorableRoles    []*discordgo.Role
	Triggered          bool
	Timer              *time.Timer
	Client             *discordgo.Session
	Strikes            []string
	Cache              *struct {
		Pins map[string][]*discordgo.Message
	}
}

func (N *NukePrediction) IncreaseSuspicionLevel(amount int) {
	if N.Triggered {
		return
	}

	N.SuspicionLevel += amount

	if N.SuspicionLevel >= 100 {
		success := N.Timer.Stop()
		if success != true {
			<-N.Timer.C
		}
		N.Lockdown()
		N.Triggered = true
	} else {
		N.ResetTimer()
	}

}

func (N *NukePrediction) Lockdown() {
	N.SuspicionLevel = 0
	roles, err := N.Client.GuildRoles(N.GuildID)

	if err != nil {
		fmt.Println(err)
	}

	for _, role := range roles {
		_, err := N.Client.GuildRoleEdit(N.GuildID, role.ID, role.Name, role.Color, role.Hoist, defaultPerms, role.Mentionable)

		if err == nil {
			N.RestorableRoles = append(N.RestorableRoles, role)
		}

	}

}

func (N *NukePrediction) Restore() {
	for _, role := range N.RestorableRoles {
		N.Client.GuildRoleEdit(N.GuildID, role.ID, role.Name, role.Color, role.Hoist, role.Permissions, role.Mentionable)
	}

	sort.SliceStable(N.RestorableChannels, func(i, j int) bool {
		if N.RestorableChannels[i].Type == discordgo.ChannelTypeGuildCategory {
			return true
		}
		return false
	})

	for _, channel := range N.RestorableChannels {
		channel1, err := N.Client.GuildChannelCreateComplex(N.GuildID, discordgo.GuildChannelCreateData{
			Name:                 channel.Name,
			Type:                 channel.Type,
			Topic:                channel.Topic,
			Bitrate:              channel.Bitrate,
			UserLimit:            channel.UserLimit,
			RateLimitPerUser:     channel.RateLimitPerUser,
			Position:             channel.Position,
			PermissionOverwrites: channel.PermissionOverwrites,
			ParentID:             channel.ParentID,
			NSFW:                 channel.NSFW,
		})

		if err == nil {
			messages := N.Cache.Pins[channel.ID]

			for _, m := range messages {
				N.Client.ChannelMessageSend(channel1.ID, m.Content)
			}
		}

		time.Sleep(time.Second * 1)
	}

	N.RestorableChannels = []*discordgo.Channel{}
	N.RestorableRoles = []*discordgo.Role{}
	N.Strikes = []string{}
	N.Triggered = false

}

func (N *NukePrediction) ResetTimer() {
	if N.Timer == nil {
		N.Timer = time.AfterFunc(12*time.Second, N.ResetSuspicon)
	} else {
		N.Timer.Stop()
		N.Timer.Reset(12 * time.Second)
	}
}

func (N *NukePrediction) ResetSuspicon() {
	N.SuspicionLevel = 0
	N.RestorableChannels = []*discordgo.Channel{}
	N.RestorableRoles = []*discordgo.Role{}
	N.Strikes = []string{}
}
