package discord

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// DataMangager(メモリ)にUser情報を一時保存
func (aManager *dataManager) updateUser(aMember *discordgo.Member, aIsMember IsMember) error {
	aManager.UsersMutex.Lock()
	defer aManager.UsersMutex.Unlock()
	aManager.Users[aMember.User.ID] = &user{
		UserID:   aMember.User.ID,
		UserName: aMember.User.Username,
		IsMember: aIsMember,
		Roles:    aMember.Roles,
	}
	return nil
}

// DataMangager(メモリ)にステータス情報を一時保存
func (aManager *dataManager) updateStatus(aVoiceState *discordgo.VoiceState, aOnline OnlineStatus) error {
	tNow := time.Now()
	aManager.StatusesMutex.Lock()
	defer aManager.StatusesMutex.Unlock()
	tStatus := &status{
		UserID:       aVoiceState.UserID,
		ChannelID:    aVoiceState.ChannelID,
		Timestamp:    tNow,
		VoiceState:   statusMap(aVoiceState),
		OnlineStatus: aOnline,
	}
	aManager.Statuses[aVoiceState.UserID] = append(aManager.Statuses[aVoiceState.UserID], tStatus)

	return nil
}
func statusMap(aVoiceState *discordgo.VoiceState) string {
	if aVoiceState.ChannelID == "" {
		return "offline"
	}
	if aVoiceState.SelfDeaf || aVoiceState.Deaf {
		return "speaker-mute"
	}
	if aVoiceState.SelfMute || aVoiceState.Mute {
		return "mic-mute"
	}
	return "mic-on"
}

// ステータス情報を保存するとき、User情報があるか確認する。
// User情報は変化がない時もあるだろうから変更があったか判定した方がよさそう
func (aManager *dataManager) flushStatuses() func() {
	return func() {
		log.Println("save data")
	}
}
