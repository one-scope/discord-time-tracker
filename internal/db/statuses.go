package db

import "fmt"

// ステータステーブルを作成
func (aDB *PostgresDB) CreateStatusesTable() error {
	tQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT PRIMARY KEY,%s TEXT NOT NULL,%s TEXT NOT NULL,%s TIMESTAMP NOT NULL,%s TEXT NOT NULL,%s TEXT NOT NULL,%s TEXT NOT NULL)",
		statusesTable, statusesTableID, statusesTablePreviusID, usersTableID, statusesTableTimestamp, statusesTableChannelID, statusesTableVoiceState, statusesTableOnlineStatus)
	_, tError := aDB.DB.Exec(tQuery)
	return tError
}

// ステータス新規登録
func (aDB *PostgresDB) InsertStatus(aStatus *Statuslog) error {
	tQuery := fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s) VALUES ($1,$2,$3,$4,$5,$6,$7)",
		statusesTable, statusesTableID, statusesTablePreviusID, usersTableID, statusesTableTimestamp, statusesTableChannelID, statusesTableVoiceState, statusesTableOnlineStatus)
	_, tError := aDB.DB.Exec(tQuery, aStatus.ID, aStatus.PreviusID, aStatus.UserID, aStatus.Timestamp, aStatus.ChannelID, aStatus.VoiceState, aStatus.OnlineStatus)
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
		if tError := tRows.Scan(&tStatus.ID, &tStatus.PreviusID, &tStatus.UserID, &tStatus.Timestamp, &tStatus.ChannelID, &tStatus.VoiceState, &tStatus.OnlineStatus); tError != nil {
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
