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
	tQuery := fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s,%s) VALUES (:discord_id,:time_stamp,:channel_id,:voice_state,:online_status)",
		statusesTable, usersTableID, statusesTableTimestamp, statusesTableChannelID, statusesTableVoiceState, statusesTableOnlineStatus)
	_, tError := aDB.DB.NamedExec(tQuery, *aStatus)
	return tError
}
