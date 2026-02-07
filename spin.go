package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
)

type SpinCycleInfo struct {
	spin *discordgo.Channel
	cycle *discordgo.Channel
	guildID string
	hasSpinCycle bool
}

var cache []SpinCycleInfo

const cacheSize = 10
func initSpinCycleCache() {
	cache = make([]SpinCycleInfo, cacheSize)
}

func checkSpinCycleCache(guildID string) (SpinCycleInfo, bool) {
	var info SpinCycleInfo
	ok := false
	for _, item := range cache {
		if item.guildID == guildID { 
			info = item 
			ok = true
			break
		}
	}

	return info, ok
}

// not running a check here because this function should only be called in case of
// a cache invalidation anyway
func updateSpinCycleCache(info SpinCycleInfo) {
	slices.Reverse(cache[:len(cache)-1])
	slices.Reverse(cache)
	cache[0] = info
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

	info.guildID = guildID
	info.hasSpinCycle = true
	if info.spin == nil || info.cycle == nil {
		info.hasSpinCycle = false
	}

	return info, nil
}

func updateSpinCycle(s *discordgo.Session, m *discordgo.VoiceStateUpdate) error {
	var info SpinCycleInfo
	info, ok := checkSpinCycleCache(m.GuildID)
	if !ok {
		info, err := getSpinCycleInfo(s, m.GuildID)
		if err != nil { return err }
		updateSpinCycleCache(info)
	}

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
