package discord

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/one-scope/discord-time-tracker/internal/db"
)

// DataMangager(メモリ)にUser情報を一時保存
func (aManager *dataManager) updateUser(aMember *discordgo.Member, aIsMember db.IsMember) error {
	aManager.UserByID[aMember.User.ID] = &db.User{
		ID:       aMember.User.ID,
		Name:     aMember.User.Username,
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

// 未実装：失敗した場合、リトライする
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

		log.Println("全てのDBにあるユーザー情報")
		tUsers, tError := aManager.DB.GetAllUsers()
		if tError != nil {
			log.Printf("failed to get all users: %v", tError)
			return
		}
		for _, tUser := range tUsers {
			log.Println(tUser)
			// ロール情報を取得
			tRoles, tError := aManager.DB.GetAllRolesByUserID(tUser.ID)
			if tError != nil {
				log.Printf("failed to get all roles: %v", tError)
				return
			}
			for _, tRole := range tRoles {
				log.Println(tRole)
			}
		}
	}
}

// 未実装：ユーザーを登録して、メモリリセットするまでの間に、更新があるかもしれない。あるなら処理もしくはロックが必要
func (aManager *dataManager) flushUsersData() error {
	// ユーザー情報の更新がないなら何もしない
	if len(aManager.UserByID) == 0 {
		return nil
	}

	// ユーザー情報をDBに保存
	for _, tUser := range aManager.UserByID {
		tIsExists, tError := aManager.DB.IsExistsUserByID(tUser.ID)
		if tError != nil {
			return fmt.Errorf("failed to check user exists: %w", tError)
		}
		if tIsExists {
			if tError := aManager.DB.UpdateUser(tUser); tError != nil {
				return fmt.Errorf("failed to update user: %w", tError)
			}
		} else {
			if tError := aManager.DB.InsertUser(tUser); tError != nil {
				return fmt.Errorf("failed to add user: %w", tError)
			}
		}
	}

	// ロール情報をDBに保存
	for _, tUser := range aManager.UserByID {
		tRolesMap, tError := aManager.DB.GetRolesIDMapByUserID(tUser.ID)
		if tError != nil {
			return fmt.Errorf("failed to get roles: %w", tError)
		}
		for _, tRole := range tUser.Roles {
			tIsExists, tError := aManager.DB.IsExistsRoleByUserID(tUser.ID, tRole)
			if tError != nil {
				return fmt.Errorf("failed to check role exists: %w", tError)
			}
			if !tIsExists {
				if tError := aManager.DB.InsertRoleByUserID(tUser.ID, tRole); tError != nil {
					return fmt.Errorf("failed to add role: %w", tError)
				}
			} else {
				delete(tRolesMap, tRole)
			}
		}
		for tRole := range tRolesMap {
			if tError := aManager.DB.DeleteRoleByUserID(tUser.ID, tRole); tError != nil {
				return fmt.Errorf("failed to delete role: %w", tError)
			}
		}
	}

	// メモリのユーザー情報を初期化
	aManager.UserByID = map[string]*db.User{}

	return nil
}

// 未実装：ステータス情報を保存するとき、User情報があるか確認してなかったらとりにいく。
func (aManager *dataManager) flushStatusesData() error {
	// ステータス情報がないなら何もしない
	if len(aManager.StatusByID) == 0 {
		return nil
	}

	return nil
}
