package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ステータステーブルを作成
func (aDB *PostgresDB) createStatusesTable() error {
	tQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT PRIMARY KEY,%s TEXT NOT NULL,%s TIMESTAMP with time zone NOT NULL,%s TEXT NOT NULL,%s TEXT NOT NULL,%s TEXT NOT NULL)",
		statusesTable, statusesTableID, usersTableID, statusesTableTimestamp, statusesTableChannelID, statusesTableVoiceState, statusesTableOnlineStatus)
	_, tError := aDB.DB.Exec(tQuery)
	return tError
}

// ステータス新規登録
func (aDB *PostgresDB) InsertStatus(aStatus *Statuslog) error {
	tQuery := fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s,%s,%s) VALUES ($1,$2,$3,$4,$5,$6)",
		statusesTable, statusesTableID, usersTableID, statusesTableTimestamp, statusesTableChannelID, statusesTableVoiceState, statusesTableOnlineStatus)
	_, tError := aDB.DB.Exec(tQuery, aStatus.ID, aStatus.UserID, aStatus.Timestamp, aStatus.ChannelID, aStatus.VoiceState, aStatus.OnlineStatus)
	return tError
}

// ユーザーIDから全てのステータスを取得
func (aDB *PostgresDB) GetAllStatusesByUserID(aUserID string) ([]*Statuslog, error) {
	tQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", statusesTable, usersTableID)
	tRows, tError := aDB.DB.Queryx(tQuery, aUserID)
	if tError != nil {
		return nil, tError
	}
	defer tRows.Close()
	var tStatuses []*Statuslog
	for tRows.Next() {
		tStatus := Statuslog{}
		if tError := tRows.Scan(&tStatus.ID, &tStatus.UserID, &tStatus.Timestamp, &tStatus.ChannelID, &tStatus.VoiceState, &tStatus.OnlineStatus); tError != nil {
			return nil, tError
		}
		tStatuses = append(tStatuses, &tStatus)
	}
	//エラーチェック
	if tError := tRows.Err(); tError != nil {
		return nil, tError
	}
	return tStatuses, nil
}

// ユーザーIDと期間を指定して昇順でステータスを取得
func (aDB *PostgresDB) GetStatusesByUserIDAndRangeAscendingOrder(aUserID string, aStart time.Time, aEnd time.Time) ([]*Statuslog, error) {
	tQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s BETWEEN $2 AND $3 ORDER BY %s ASC",
		statusesTable, usersTableID, statusesTableTimestamp, statusesTableTimestamp)
	tRows, tError := aDB.DB.Query(tQuery, aUserID, aStart, aEnd)
	if tError != nil {
		return nil, tError
	}
	defer tRows.Close()

	var tStatuses []*Statuslog
	for tRows.Next() {
		tStatus := Statuslog{}
		if tError := tRows.Scan(&tStatus.ID, &tStatus.UserID, &tStatus.Timestamp, &tStatus.ChannelID, &tStatus.VoiceState, &tStatus.OnlineStatus); tError != nil {
			return nil, tError
		}
		tStatuses = append(tStatuses, &tStatus)
		//エラーチェック
		if tError := tRows.Err(); tError != nil {
			return nil, tError
		}
	}
	return tStatuses, nil
}

// LogIDを指定してステータスを取得
func (aDB *PostgresDB) GetStatusByLogID(aLogID string) (*Statuslog, error) {
	tQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", statusesTable, statusesTableID)
	tRow := aDB.DB.QueryRowx(tQuery, aLogID)
	tStatus := Statuslog{}
	if tError := tRow.Scan(&tStatus.ID, &tStatus.UserID, &tStatus.Timestamp, &tStatus.ChannelID, &tStatus.VoiceState, &tStatus.OnlineStatus); tError != nil {
		return nil, tError
	}
	return &tStatus, nil
}

// 日付を指定してそれより前で一番近いステータスを取得.なければオフラインを返す
func (aDB *PostgresDB) GetRecentStatusByUserIDAndTimestamp(aUserID string, aTimestamp time.Time) (*Statuslog, error) {
	tQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 AND %s <= $2 ORDER BY %s DESC LIMIT 1",
		statusesTable, usersTableID, statusesTableTimestamp, statusesTableTimestamp)
	tRow := aDB.DB.QueryRowx(tQuery, aUserID, aTimestamp)
	tStatus := Statuslog{}
	if tError := tRow.Scan(&tStatus.ID, &tStatus.UserID, &tStatus.Timestamp, &tStatus.ChannelID, &tStatus.VoiceState, &tStatus.OnlineStatus); tError != nil {
		if errors.Is(tError, sql.ErrNoRows) {
			return &Statuslog{
				ChannelID:    "",
				VoiceState:   VoiceOffline,
				OnlineStatus: Offline,
			}, nil
		} else {
			return nil, tError
		}
	}
	return &tStatus, nil
}
