package main

import (
	AntiSpam "./antispam"
	"./nukeprediction"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var client *discordgo.Session
var nukePredictors []*nukeprediction.NukePrediction
var antispam *AntiSpam.AntiSpam

var prefix string

func main() {
	godotenv.Load()

	client, _ = discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	err := client.Open()
	if err != nil {
		fmt.Println(err)
	}

	prefix = os.Getenv("PREFIX")

	go fillPredictors()

	antispam = &AntiSpam.AntiSpam{
		Users: make(map[string]struct {
			N int
			M []string
		}),
	}

	antispam.Init()

	client.AddHandler(OnMsg)

	// sus events
	client.AddHandler(channelDeleted)
	client.AddHandler(memberBanned)
	client.AddHandler(memberKicked)

	fmt.Printf("Now running | Logged in as %s\n", client.State.User.Username)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
