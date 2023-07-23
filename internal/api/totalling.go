package api

import (
	"log"
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
		StatusesByUserID: nil,
	}
	//返り値のステータス部分
	tApiStatusesByUserID := map[string][]*TotalStatus{}

	//ユーザーごとに集計
	for _, tUserID := range tUserIDs {
		tStatuses := []*TotalStatus{}

		//集計用のステータス
		tNowOnlineStatus := string(db.UnknownOnline)
		tNowVoiceState := string(db.VoiceUnknown)
		tNowChannel := db.ChannelUnknown
		//Periodごとに集計
		for tStart := aStart; tStart.Before(aEnd); tStart = tStart.Add(aPeriod) {
			//未実装：timezoneを揃える。DBでUTCになっちゃう
			tEnd := tStart.Add(aPeriod)
			log.Println(time.Now().In(time.UTC).Add(time.Hour * 9))
			if tEnd.After(time.Now().In(time.UTC).Add(time.Hour * 9)) {
				tEnd = time.Now().In(time.UTC).Add(time.Hour * 9)
			}

			//ステータスログを取得
			tLogStatuses, tError := aDB.GetStatusesByUserIDAndRangeAscendingOrder(tUserID, tStart, tEnd)
			if tError != nil {
				return nil, tError
			}
			//未実装：現在のステータスから集計をする
			//未実装：現在のステータスが不明、初期値なら区間外の直近のログを取得して集計をする
			//ログがないならスキップ、未実装：直近のログがないならcontinue
			if len(tLogStatuses) == 0 {
				continue
			}

			tStatus := &TotalStatus{
				Start:          tStart,
				End:            tEnd,
				ChannelByID:    map[string]TotalChannel{},
				VoiceByState:   map[string]TotalVoiceState{},
				OnlineByStatus: map[string]TotalOnlineStatus{},
			}

			tNowTime := tStart
			for _, tLogStatus := range tLogStatuses {

				//チャンネル
				if tNowChannel == db.ChannelUnknown {
					tNowVoiceState, tError = initChannel(aDB, tLogStatus)
					if tError != nil {
						return nil, tError
					}
				}
				tStatus.totalChannel(tNowChannel, tNowTime, tLogStatus.Timestamp)

				//ボイス
				if tNowVoiceState == string(db.VoiceUnknown) {
					tNowVoiceState, tError = initVoiceState(aDB, tLogStatus)
					if tError != nil {
						return nil, tError
					}
				}
				tStatus.totalVoiceState(tNowVoiceState, tNowTime, tLogStatus.Timestamp)

				//オンライン
				//最初の常態を取得
				if tNowOnlineStatus == string(db.UnknownOnline) {
					tNowOnlineStatus, tError = initOnlineStatus(aDB, tLogStatus)
					if tError != nil {
						return nil, tError
					}
				}
				//オンライン集計
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

		tApiStatusesByUserID[tUserID] = tStatuses
	}

	tApiStatuses.StatusesByUserID = tApiStatusesByUserID
	return tApiStatuses, nil
}

func initChannel(aDB *db.PostgresDB, aLogStatus *db.Statuslog) (string, error) {
	if aLogStatus.PreviusID != "" {
		tStatusLog, tError := aDB.GetStatusByLogID(aLogStatus.PreviusID)
		if tError != nil {
			return string(db.ChannelUnknown), tError
		}
		return tStatusLog.ChannelID, nil
	}
	return "", nil
}

func initVoiceState(aDB *db.PostgresDB, aLogStatus *db.Statuslog) (string, error) {
	if aLogStatus.PreviusID != "" {
		tStatusLog, tError := aDB.GetStatusByLogID(aLogStatus.PreviusID)
		if tError != nil {
			return string(db.VoiceUnknown), tError
		}
		return string(tStatusLog.VoiceState), nil
	}

	return string(db.VoiceOffline), nil
}

func initOnlineStatus(aDB *db.PostgresDB, aLogStatus *db.Statuslog) (string, error) {
	if aLogStatus.PreviusID != "" {
		tStatusLog, tError := aDB.GetStatusByLogID(aLogStatus.PreviusID)
		if tError != nil {
			return string(db.UnknownOnline), tError
		}
		return string(tStatusLog.OnlineStatus), nil
	}

	return string(db.Offline), nil
}

func (aStatus *TotalStatus) totalChannel(aNowChannel string, aStartTime time.Time, aEndTime time.Time) {
	if aNowChannel == db.ChannelUnknown || aNowChannel == "" {
		return
	}
	tChannelTotal := aStatus.ChannelByID[aNowChannel]
	tChannelTotal.TotalTime += aEndTime.Sub(aStartTime)
	aStatus.ChannelByID[aNowChannel] = tChannelTotal
}

func (aStatus *TotalStatus) totalVoiceState(aNowVoiceState string, aStartTime time.Time, aEndTime time.Time) {
	if aNowVoiceState == string(db.VoiceUnknown) {
		return
	}
	tVoiceTotal := aStatus.VoiceByState[aNowVoiceState]
	tVoiceTotal.TotalTime += aEndTime.Sub(aStartTime)
	aStatus.VoiceByState[aNowVoiceState] = tVoiceTotal
}

func (aStatus *TotalStatus) totalOnlineStatus(aNowOnlineStatus string, aStartTime time.Time, aEndTime time.Time) {
	log.Println("totalOnlineStatus")
	log.Println("StartTime-EndTime", aStartTime, aEndTime)
	if aNowOnlineStatus == string(db.UnknownOnline) {
		return
	}
	tOnlineTotal := aStatus.OnlineByStatus[aNowOnlineStatus]
	tOnlineTotal.TotalTime += aEndTime.Sub(aStartTime)
	aStatus.OnlineByStatus[aNowOnlineStatus] = tOnlineTotal
}
