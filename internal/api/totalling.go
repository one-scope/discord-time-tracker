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

func AggregateStatusWithinRangeByUserIDs(aDB *db.PostgresDB, aStart time.Time, aEnd time.Time, aPeriod time.Duration, tUserIDs []string) (*AllUsersTotalStatuses, error) {
	//返り値
	tAllUsersStatuses := &AllUsersTotalStatuses{
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
			tNowStart := tStart
			//未実装：timezoneを揃える。DBでUTCになっちゃう
			tNowEnd := tNowStart.Add(aPeriod)
			if tNowEnd.After(time.Now().In(time.UTC).Add(time.Hour * 9)) {
				tNowEnd = time.Now().In(time.UTC).Add(time.Hour * 9)
			}

			//ステータスログを取得
			tLogStatuses, tError := aDB.GetStatusesByUserIDAndRangeAscendingOrder(tUserID, tNowStart, tNowEnd)
			if tError != nil {
				return nil, tError
			}
			//ログがなくオフラインの場合は集計しない
			if (len(tLogStatuses) == 0) && (tNowChannel == "") && (tNowOnlineStatus == string(db.Offline)) && (tNowVoiceState == string(db.VoiceOffline)) {
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
				tStatus.totalChannel(tNowChannel, tNowStart, tLogStatus.Timestamp)
				tStatus.totalVoiceState(tNowVoiceState, tNowStart, tLogStatus.Timestamp)
				tStatus.totalOnlineStatus(tNowOnlineStatus, tNowStart, tLogStatus.Timestamp)

				//集計用ステータス更新
				tNowChannel = tLogStatus.ChannelID
				tNowVoiceState = string(tLogStatus.VoiceState)
				tNowOnlineStatus = string(tLogStatus.OnlineStatus)
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
