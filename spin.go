package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type SpinCycleInfo struct {
	spin *discordgo.Channel
	cycle *discordgo.Channel
	hasSpinCycle bool
}

func getSpinCycleInfo(s *discordgo.Session, guildID string) (SpinCycleInfo, error) {
	var info SpinCycleInfo

	channels, err := s.GuildChannels(guildID)
	if err != nil { return info, err }

	for _, channel := range channels {
		if channel.Name == "spin" {
			info.spin = channel
		}

		if channel.Name == "cycle" {
			info.cycle = channel
		}
	}

	info.hasSpinCycle = true
	if info.spin == nil || info.cycle == nil {
		info.hasSpinCycle = false
	}

	return info, nil
}

func updateSpinCycle(s *discordgo.Session, m *discordgo.VoiceStateUpdate) error {
	info, err := getSpinCycleInfo(s, m.GuildID)
	if err != nil { return err }

	if !info.hasSpinCycle { return nil }

	time.Sleep(time.Millisecond * 750)

	if m.VoiceState.ChannelID == info.spin.ID {
		s.GuildMemberMove(info.cycle.GuildID, m.UserID, &info.cycle.ID)
	}

	if m.VoiceState.ChannelID == info.cycle.ID {
		s.GuildMemberMove(info.spin.GuildID, m.UserID, &info.spin.ID)
	}

	return nil
}

func voiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	err := updateSpinCycle(s, m)
	if err != nil {
		fmt.Println(err.Error())
	}
}
