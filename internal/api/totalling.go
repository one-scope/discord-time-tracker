package api

import (
	"time"

	"github.com/one-scope/discord-time-tracker/internal/db"
)

func GetAllUsers(aDB *db.PostgresDB) ([]*db.User, error) {
	tUsers, tError := aDB.GetAllUsers()
	if tError != nil {
		return nil, tError
	}
	for _, tUser := range tUsers {
		tRoles, tError := aDB.GetAllRolesIDByUserID(tUser.ID)
		if tError != nil {
			return nil, tError
		}
		tUser.Roles = tRoles
	}

	return tUsers, nil
}

type AllUsersTotalStatuses struct {
	Start            time.Time
	End              time.Time
	Period           time.Duration
	StatusesByUserID map[string][]*TotalStatus
}
type TotalStatus struct {
	Start          time.Time
	End            time.Time
	ChannelByID    map[string]TotalChannel
	VoiceByState   map[string]TotalVoiceState
	OnlineByStatus map[string]TotalOnlineStatus
}
type TotalChannel struct {
	TotalTime time.Duration
}
type TotalVoiceState struct {
	TotalTime time.Duration
}
type TotalOnlineStatus struct {
	TotalTime time.Duration
}

func AggregateStatusWithinRangeByUserIDs(aDB *db.PostgresDB, aStart time.Time, aEnd time.Time, aPeriod time.Duration, tUserIDs []string) (*AllUsersTotalStatuses, error) {
	//返り値
	tApiStatuses := &AllUsersTotalStatuses{
		Start:            aStart,
		End:              aEnd,
		Period:           aPeriod,
		StatusesByUserID: map[string][]*TotalStatus{},
	}

	//ユーザーごとに集計
	for _, tUserID := range tUserIDs {
		tStatuses := []*TotalStatus{}

		//初期ステータス設定
		tInitLogStatuses, tError := aDB.GetRecentStatusByUserIDAndTimestamp(tUserID, aStart)
		if tError != nil {
			return nil, tError
		}
		tNowChannel := tInitLogStatuses.ChannelID
		tNowOnlineStatus := string(tInitLogStatuses.OnlineStatus)
		tNowVoiceState := string(tInitLogStatuses.VoiceState)
		//Periodごとに集計
		for tStart := aStart; tStart.Before(aEnd); tStart = tStart.Add(aPeriod) {
			tNowTime := tStart
			//未実装：timezoneを揃える。DBでUTCになっちゃう
			tEnd := tNowTime.Add(aPeriod)
			if tEnd.After(time.Now().In(time.UTC).Add(time.Hour * 9)) {
				tEnd = time.Now().In(time.UTC).Add(time.Hour * 9)
			}

			//ステータスログを取得
			tLogStatuses, tError := aDB.GetStatusesByUserIDAndRangeAscendingOrder(tUserID, tNowTime, tEnd)
			if tError != nil {
				return nil, tError
			}
			//ログがないならスキップ
			if len(tLogStatuses) == 0 {
				continue
			}

			tStatus := &TotalStatus{
				Start:          tNowTime,
				End:            tEnd,
				ChannelByID:    map[string]TotalChannel{},
				VoiceByState:   map[string]TotalVoiceState{},
				OnlineByStatus: map[string]TotalOnlineStatus{},
			}
			for _, tLogStatus := range tLogStatuses {
				tStatus.totalChannel(tNowChannel, tNowTime, tLogStatus.Timestamp)
				tStatus.totalVoiceState(tNowVoiceState, tNowTime, tLogStatus.Timestamp)
				tStatus.totalOnlineStatus(tNowOnlineStatus, tNowTime, tLogStatus.Timestamp)

				//集計用ステータス更新
				tNowChannel = tLogStatus.ChannelID
				tNowVoiceState = string(tLogStatus.VoiceState)
				tNowOnlineStatus = string(tLogStatus.OnlineStatus)
				tNowTime = tLogStatus.Timestamp
			}
			//最後のTimeStampからPeriodの終わりまでのオンライン集計
			tStatus.totalChannel(tNowChannel, tNowTime, tEnd)
			tStatus.totalVoiceState(tNowVoiceState, tNowTime, tEnd)
			tStatus.totalOnlineStatus(tNowOnlineStatus, tNowTime, tEnd)

			tStatuses = append(tStatuses, tStatus)
		}

		tApiStatuses.StatusesByUserID[tUserID] = tStatuses
	}

	return tApiStatuses, nil
}

func (aStatus *TotalStatus) totalChannel(aNowChannel string, aStartTime time.Time, aEndTime time.Time) {
	if aNowChannel == "" {
		return
	}
	tChannelTotal := aStatus.ChannelByID[aNowChannel]
	tChannelTotal.TotalTime += aEndTime.Sub(aStartTime)
	aStatus.ChannelByID[aNowChannel] = tChannelTotal
}
func (aStatus *TotalStatus) totalVoiceState(aNowVoiceState string, aStartTime time.Time, aEndTime time.Time) {
	tVoiceTotal := aStatus.VoiceByState[aNowVoiceState]
	tVoiceTotal.TotalTime += aEndTime.Sub(aStartTime)
	aStatus.VoiceByState[aNowVoiceState] = tVoiceTotal
}
func (aStatus *TotalStatus) totalOnlineStatus(aNowOnlineStatus string, aStartTime time.Time, aEndTime time.Time) {
	tOnlineTotal := aStatus.OnlineByStatus[aNowOnlineStatus]
	tOnlineTotal.TotalTime += aEndTime.Sub(aStartTime)
	aStatus.OnlineByStatus[aNowOnlineStatus] = tOnlineTotal
}
