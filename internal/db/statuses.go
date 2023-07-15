package db

import "fmt"

// ステータステーブルを作成
func (aDB *PostgresDB) CreateStatusesTable() error {
	tQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT NOT NULL,%s TIMESTAMP NOT NULL,%s TEXT NOT NULL,%s TEXT NOT NULL,%s TEXT NOT NULL, UNIQUE(%s,%s))",
		statusesTable, usersTableID, statusesTableTimestamp, statusesTableChannelID, statusesTableVoiceState, statusesTableOnlineStatus, usersTableID, statusesTableTimestamp)
	_, tError := aDB.DB.Exec(tQuery)
	return tError
}

// ステータス新規登録
func (aDB *PostgresDB) InsertStatus(aStatus *Statuslog) error {
	tQuery := fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s,%s) VALUES ($1,$2,$3,$4,$5)",
		statusesTable, usersTableID, statusesTableTimestamp, statusesTableChannelID, statusesTableVoiceState, statusesTableOnlineStatus)
	_, tError := aDB.DB.Exec(tQuery, aStatus.UserID, aStatus.Timestamp, aStatus.ChannelID, aStatus.VoiceState, aStatus.OnlineStatus)
	return tError
}
