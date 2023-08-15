package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/one-scope/discord-time-tracker/internal/db"
)

// DataMangager(メモリ)にUser情報を一時保存
func (aManager *dataManager) updateUser(aMember *discordgo.Member, aIsMember db.IsMember) error {
	aManager.UsersByID[aMember.User.ID] = &db.User{
		ID:       aMember.User.ID,
		Name:     aMember.User.Username,
		IsMember: aIsMember,
		Roles:    aMember.Roles,
	}
	return nil
}

// DataMangager(メモリ)にステータス情報を一時保存
func (aManager *dataManager) updateStatus(aVoiceState *discordgo.VoiceState, aOnline db.OnlineStatus) error {
	tNow := time.Now()
	tUUID := fmt.Sprintf("%v", uuid.New())
	tOnline := aManager.determineOnlineStatus(aVoiceState.UserID, aOnline)
	tChannelID := aManager.determineChannelID(aVoiceState.UserID, aVoiceState.ChannelID)

	tStatus := &db.Statuslog{
		ID:           tUUID,
		UserID:       aVoiceState.UserID,
		ChannelID:    tChannelID,
		Timestamp:    tNow,
		VoiceState:   statusMap(aVoiceState),
		OnlineStatus: tOnline,
	}
	aManager.PreViusStatusLogByUserID[aVoiceState.UserID] = tStatus
	aManager.StatusesByID[aVoiceState.UserID] = append(aManager.StatusesByID[aVoiceState.UserID], tStatus)

	return nil
}

// unknownOnlineの場合に前回のログを参照してChannelIDを決定する
func (aManager *dataManager) determineChannelID(aUserID string, aChannelID string) string {
	tChannelID := aChannelID
	//ステータスが分からない場合、前回のログを参照する
	if tChannelID == db.UnknownChannelID {
		if tPreviusLog, tOk := aManager.PreViusStatusLogByUserID[aUserID]; tOk {
			tChannelID = tPreviusLog.ChannelID
		} else {
			tChannelID = ""
		}
	}

	return tChannelID
}

// unknownOnlineの場合に前回のログを参照してOnlineStatusを決定する
func (aManager *dataManager) determineOnlineStatus(aUserID string, aOnlineStatus db.OnlineStatus) db.OnlineStatus {
	tOnlineStatus := aOnlineStatus
	//ステータスが分からない場合、前回のログを参照する
	if tOnlineStatus == db.UnknownOnline {
		if tPreviusLog, tOk := aManager.PreViusStatusLogByUserID[aUserID]; tOk {
			tOnlineStatus = tPreviusLog.OnlineStatus
		} else {
			tOnlineStatus = db.Offline
		}
	}
	return tOnlineStatus
}
func statusMap(aVoiceState *discordgo.VoiceState) db.VoiceState {
	if aVoiceState.ChannelID == "" {
		return db.VoiceOffline
	}
	if aVoiceState.SelfDeaf || aVoiceState.Deaf {
		return db.VoiceDeaf
	}
	if aVoiceState.SelfMute || aVoiceState.Mute {
		return db.VoiceMute
	}
	return db.VoiceOn
}

func (aManager *dataManager) flushData() error {
	if tError := aManager.flushUsersData(); tError != nil {
		return fmt.Errorf("failed to flush users data: %v", tError)
	}
	if tError := aManager.flushStatusesData(); tError != nil {
		return fmt.Errorf("failed to flush statuses data: %v", tError)
	}
	return nil
}

func (aManager *dataManager) flushUsersData() error {
	// ユーザー情報の更新がないなら何もしない
	if len(aManager.UsersByID) == 0 {
		return nil
	}

	// ユーザー情報をDBに保存
	for _, tUser := range aManager.UsersByID {
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
	for _, tUser := range aManager.UsersByID {
		tRolesMap, tError := aManager.DB.GetAllRolesIDMapByUserID(tUser.ID)
		if tError != nil {
			return fmt.Errorf("failed to get roles: %w", tError)
		}
		for _, tRole := range tUser.Roles {
			tIsExists, tError := aManager.DB.IsExistsRoleByUserID(tUser.ID, tRole)
			if tError != nil {
				return fmt.Errorf("failed to check role exists: %w", tError)
			}
			if !tIsExists { // ロールがないなら追加
				if tError := aManager.DB.InsertRoleByUserID(tUser.ID, tRole); tError != nil {
					return fmt.Errorf("failed to add role: %w", tError)
				}
			} else { // ロールがあるなら削除対象から外す
				delete(tRolesMap, tRole)
			}
		}
		for tRole := range tRolesMap { // 削除対象のロールを削除
			if tError := aManager.DB.DeleteRoleByUserID(tUser.ID, tRole); tError != nil {
				return fmt.Errorf("failed to delete role: %w", tError)
			}
		}
	}

	// メモリのユーザー情報を初期化
	aManager.UsersByID = map[string]*db.User{}

	return nil
}

func (aManager *dataManager) flushStatusesData() error {
	// ステータス情報がないなら何もしない
	if len(aManager.StatusesByID) == 0 {
		return nil
	}

	// ステータス情報をDBに保存
	for _, tStatuses := range aManager.StatusesByID {
		for _, tStatus := range tStatuses {
			if tError := aManager.DB.InsertStatus(tStatus); tError != nil {
				return fmt.Errorf("failed to add status: %w", tError)
			}
		}
	}

	// メモリのステータス情報を初期化
	aManager.StatusesByID = map[string][]*db.Statuslog{}

	return nil
}
