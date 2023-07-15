package db

import "fmt"

// ステータス新規登録
func (aDB *PostgresDB) InsertStatus(aStatus *Statuslog) error {
	tQuery := fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s,%s) VALUES (:discord_id,:time_stamp,:channel_id,:voice_state,:online_status)",
		statusesTable, usersTableID, statusesTableTimestamp, statusesTableChannelID, statusesTableVoiceState, statusesTableOnlineStatus)
	_, tError := aDB.DB.NamedExec(tQuery, *aStatus)
	return tError
}
