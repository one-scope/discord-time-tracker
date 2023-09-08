package api

import (
	"time"

	"github.com/one-scope/discord-time-tracker/internal/db"
)

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

func GetTotalStatusesByUsersID(aDB *db.PostgresDB, aStart time.Time, aEnd time.Time, aPeriod time.Duration, tUserIDs []string) (*AllUsersTotalStatuses, error) {
	//返り値
	tAllUsersStatuses := &AllUsersTotalStatuses{
		Start:            aStart,
		End:              aEnd,
		Period:           aPeriod,
		StatusesByUserID: map[string][]*TotalStatus{},
	}

	tTimeNow := time.Now()

	//ユーザーごとに集計
	for _, tUserID := range tUserIDs {
		tStatuses := []*TotalStatus{}

		//初期ステータス設定
		tInitLogStatuses, tError := aDB.GetRecentStatusByUserIDAndTimestamp(tUserID, aStart)
		if tError != nil {
			return nil, tError
		}
		tNowChannel := tInitLogStatuses.ChannelID
		tNowOnlineStatus := tInitLogStatuses.OnlineStatus
		tNowVoiceState := tInitLogStatuses.VoiceState

		//Periodごとに集計
		for tStart := aStart; tStart.Before(aEnd) && tStart.Before(tTimeNow); tStart = tStart.Add(aPeriod) {
			tNowStart := tStart
			tNowEnd := tNowStart.Add(aPeriod)
			if tNowEnd.After(tTimeNow) { // 集計終了時間が現在より未来の場合は現在まで集計
				tNowEnd = tTimeNow
			}

			//ステータスログを取得
			tLogStatuses, tError := aDB.GetStatusesByUserIDAndRangeAscendingOrder(tUserID, tNowStart, tNowEnd)
			if tError != nil {
				return nil, tError
			}
			//ログがなくオフラインの場合は集計しない
			if (len(tLogStatuses) == 0) && (tNowChannel == "") && (tNowOnlineStatus == db.Offline) && (tNowVoiceState == db.VoiceOffline) {
				continue
			}

			tStatus := &TotalStatus{
				Start:          tNowStart,
				End:            tNowEnd,
				ChannelByID:    map[string]TotalChannel{},
				VoiceByState:   map[string]TotalVoiceState{},
				OnlineByStatus: map[string]TotalOnlineStatus{},
			}
			for _, tLogStatus := range tLogStatuses {
				//ログの間の集計
				tStatus.totalChannel(tNowChannel, tNowStart, tLogStatus.Timestamp)
				tStatus.totalVoiceState(tNowVoiceState, tNowStart, tLogStatus.Timestamp)
				tStatus.totalOnlineStatus(tNowOnlineStatus, tNowStart, tLogStatus.Timestamp)

				//集計用ステータス更新
				tNowChannel = tLogStatus.ChannelID
				tNowVoiceState = tLogStatus.VoiceState
				tNowOnlineStatus = tLogStatus.OnlineStatus
				tNowStart = tLogStatus.Timestamp
			}
			//最後のTimeStampからPeriodの終わりまでのオンライン集計
			tStatus.totalChannel(tNowChannel, tNowStart, tNowEnd)
			tStatus.totalVoiceState(tNowVoiceState, tNowStart, tNowEnd)
			tStatus.totalOnlineStatus(tNowOnlineStatus, tNowStart, tNowEnd)

			tStatuses = append(tStatuses, tStatus)
		}

		tAllUsersStatuses.StatusesByUserID[tUserID] = tStatuses
	}

	return tAllUsersStatuses, nil
}

func (aStatus *TotalStatus) totalChannel(aNowChannel string, aStartTime time.Time, aEndTime time.Time) {
	if aNowChannel == "" {
		aNowChannel = "offline"
	}
	tChannelTotal := aStatus.ChannelByID[aNowChannel]
	tChannelTotal.TotalTime += aEndTime.Sub(aStartTime)
	aStatus.ChannelByID[aNowChannel] = tChannelTotal
}
func (aStatus *TotalStatus) totalVoiceState(aNowVoiceState db.VoiceState, aStartTime time.Time, aEndTime time.Time) {
	tNowVoiceState := string(aNowVoiceState)
	tVoiceTotal := aStatus.VoiceByState[tNowVoiceState]
	tVoiceTotal.TotalTime += aEndTime.Sub(aStartTime)
	aStatus.VoiceByState[tNowVoiceState] = tVoiceTotal
}
func (aStatus *TotalStatus) totalOnlineStatus(aNowOnlineStatus db.OnlineStatus, aStartTime time.Time, aEndTime time.Time) {
	tNowOnlineStatus := string(aNowOnlineStatus)
	tOnlineTotal := aStatus.OnlineByStatus[tNowOnlineStatus]
	tOnlineTotal.TotalTime += aEndTime.Sub(aStartTime)
	aStatus.OnlineByStatus[tNowOnlineStatus] = tOnlineTotal
}
