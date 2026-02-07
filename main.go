package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	initSpinCycleCache()
	err := runBot()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func runBot() error {
	contents, err := os.ReadFile("./token")
	if err != nil { return err }

	token := strings.ReplaceAll(string(contents), "\n", "")
	discord, err := discordgo.New("Bot " + token)
	if err != nil { return err }

	discord.AddHandler(voiceStateUpdate)

	discord.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuildMembers

	err = discord.Open()
	if err != nil { return err }

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()

	return nil;
}

