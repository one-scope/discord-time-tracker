package discord

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/one-scope/discord-time-tracker/internal/dbhandler"
)

// DataMangager(メモリ)にUser情報を一時保存
func (aManager *dataManager) updateUser(aMember *discordgo.Member, aIsMember dbhandler.IsMember) error {
	aManager.UserByID[aMember.User.ID] = &dbhandler.User{
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
	tStatus := &statuslog{
		UserID:       aVoiceState.UserID,
		ChannelID:    aVoiceState.ChannelID,
		Timestamp:    tNow,
		VoiceState:   statusMap(aVoiceState),
		OnlineStatus: aOnline,
	}
	aManager.StatusByID[aVoiceState.UserID] = append(aManager.StatusByID[aVoiceState.UserID], tStatus)

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
		if tError := aManager.flushUsersData(); tError != nil {
			log.Printf("failed to flush users data: %v", tError)
			return
		}
		if tError := aManager.flushStatusesData(); tError != nil {
			log.Printf("failed to flush statuses data: %v", tError)
			return
		}
	}
}

func (aManager *dataManager) flushUsersData() error {
	// ユーザー情報がないなら何もしない
	if len(aManager.UserByID) == 0 {
		return nil
	}

	// ユーザーディレクトリがないなら作成
	if tError := dbhandler.CreateUsersDataDirectory(aManager.DataPathBase); tError != nil {
		return fmt.Errorf("failed to create users directory: %w", tError)
	}

	tUsers := map[string]*dbhandler.User{}

	// ユーザーファイルがあるなら読み込み
	if dbhandler.IsExistsUsersJsonFile(aManager.DataPathBase) {
		if tError := dbhandler.DecodeUsersJsonFile(aManager.DataPathBase, &tUsers); tError != nil {
			return fmt.Errorf("failed to decode users file: %w", tError)
		}
	}

	// メモリのユーザー情報をファイルのユーザー情報に上書き
	for _, tUser := range aManager.UserByID {
		tUsers[tUser.UserID] = tUser
	}

	// メモリのユーザー情報をファイルに書き込み
	if tError := dbhandler.EncodeUsersJsonFile(aManager.DataPathBase, &tUsers); tError != nil {
		return fmt.Errorf("failed to encode users file: %w", tError)
	}

	// 古いファイルを新規ユーザーファイルに置き換え
	if tError := dbhandler.RenameUsersJsonFile(aManager.DataPathBase); tError != nil {
		return fmt.Errorf("failed to rename users file: %w", tError)
	}

	// メモリのユーザー情報を初期化
	aManager.UserByID = map[string]*dbhandler.User{}

	return nil
}

func (aManager *dataManager) flushStatusesData() error {
	// ステータス情報がないなら何もしない
	if len(aManager.StatusByID) == 0 {
		return nil
	}

	return nil
}
