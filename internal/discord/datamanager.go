package discord

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	fdhandler "github.com/one-scope/discord-time-tracker/internal/filedirectoryhandler"
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
	tStatus := &statuslog{
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

// 未実装：ステータス情報を保存するとき、User情報があるか確認してなかったらとりにいく。
// 未実装：User情報は変化がない時もあるだろうから変更があったか判定した方がよさそう
// 未実装：失敗した場合、リトライする
// 未実装：バックアップも定期的に作りたいなー
func (aManager *dataManager) flushData() func() {
	return func() {
		log.Println("save data")
		// データディレクトリ作成
		if tError := fdhandler.CreateDataDirectory(aManager.DataPathBase); tError != nil {
			log.Printf("failed to create base data directory: %v", tError)
			return
		}
		if tError := aManager.flushUsersData(); tError != nil {
			log.Printf("failed to flush users data: %v", tError)
			return
		}
	}
}

func (aManager *dataManager) flushUsersData() error {
	// ロック
	aManager.UsersMutex.Lock()
	defer aManager.UsersMutex.Unlock()
	tUsers := map[string]*user{}

	// ユーザーファイルがあるなら読み込み
	if fdhandler.IsExistsUsersFile(aManager.DataPathBase) {
		if tError := fdhandler.DecodeUsersFile(aManager.DataPathBase, &tUsers); tError != nil {
			return fmt.Errorf("failed to decode users file: %w", tError)
		}
	}
	if tUsers == nil { //Decodeでnilにされることがあるため
		tUsers = map[string]*user{}
	}

	// メモリのユーザー情報をファイルのユーザー情報に上書き
	for _, tUser := range aManager.Users {
		tUsers[tUser.UserID] = tUser
	}

	// メモリのユーザー情報をファイルに書き込み
	if tError := fdhandler.EncodeUsersFile(aManager.DataPathBase, tUsers); tError != nil {
		return fmt.Errorf("failed to encode users file: %w", tError)
	}

	// 古いファイルを新規ユーザーファイルに置き換え
	if tError := fdhandler.RenameUserFile(aManager.DataPathBase); tError != nil {
		return fmt.Errorf("failed to rename users file: %w", tError)
	}
	return nil
}
