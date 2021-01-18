package main

import (
	"fmt"

	"time"

	"github.com/bwmarrin/discordgo"

	"./nukeprediction"
)

func fillPredictors() {
	for _, guild := range client.State.Guilds {
		predictor := &nukeprediction.NukePrediction{
			GuildID: guild.ID,
			Client:  client,
			Cache: &struct {
				Pins map[string][]*discordgo.Message
			}{
				Pins: map[string][]*discordgo.Message{},
			},
		}

		channels, _ := client.GuildChannels(guild.ID)

		for _, channel := range channels {
			messages, err := client.ChannelMessagesPinned(channel.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			predictor.Cache.Pins[channel.ID] = messages
		}

		nukePredictors[guild.ID] = predictor

		time.Sleep(time.Second * 1)
	}
}



