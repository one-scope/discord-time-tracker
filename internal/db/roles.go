package db

import (
	"fmt"
)

func (aDB *PostgresDB) CreateRolesTable() error {
	tQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT PRIMARY KEY ,%s TEXT NOT NULL,%s BOOLEAN NOT NULL,%s BOOLEAN NOT NULL,%s BOOLEAN NOT NULL,%s INTEGER NOT NULL,%s INTEGER NOT NULL)",
		rolesTable, rolesTableID, rolesTableName, rolesTableManaged, rolesTableMentionable, rolesTableHoist, rolesTableColor, rolesTablePosition)
	_, tError := aDB.DB.Exec(tQuery)
	return tError
}

// ロールが存在するか確認する
func (aDB *PostgresDB) IsExistsRoleByID(aRoleID string) (bool, error) {
	tQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s=$1)", rolesTable, rolesTableID)
	tRow := aDB.DB.QueryRow(tQuery, aRoleID)
	var tExists bool
	tError := tRow.Scan(&tExists)
	return tExists, tError
}

// ロールを作成する
func (aDB *PostgresDB) InsertRole(aRole *Role) error {
	tQuery := fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s) VALUES ($1,$2,$3,$4,$5,$6,$7)", rolesTable, rolesTableID, rolesTableName, rolesTableManaged, rolesTableMentionable, rolesTableHoist, rolesTableColor, rolesTablePosition)
	_, tError := aDB.DB.Exec(tQuery, aRole.ID, aRole.Name, aRole.Managed, aRole.Mentionable, aRole.Hoist, aRole.Color, aRole.Position)
	return tError
}

// ロールを削除する
func (aDB *PostgresDB) DeleteRole(aRoleID string) error {
	tQuery := fmt.Sprintf("DELETE FROM %s WHERE %s=$1", rolesTable, rolesTableID)
	_, tError := aDB.DB.Exec(tQuery, aRoleID)
	return tError
}

// ロールを更新する
func (aDB *PostgresDB) UpdateRole(aRole *Role) error {
	tQuery := fmt.Sprintf("UPDATE %s SET %s=$1,%s=$2,%s=$3,%s=$4,%s=$5,%s=$6 WHERE %s=$7", rolesTable, rolesTableName, rolesTableManaged, rolesTableMentionable, rolesTableHoist, rolesTableColor, rolesTablePosition, rolesTableID)
	_, tError := aDB.DB.Exec(tQuery, aRole.Name, aRole.Managed, aRole.Mentionable, aRole.Hoist, aRole.Color, aRole.Position, aRole.ID)
	return tError
}

// 全てのロールを取得する
func (aDB *PostgresDB) GetAllRoles() ([]*Role, error) {
	tQuery := fmt.Sprintf("SELECT * FROM %s", rolesTable)
	tRows, tError := aDB.DB.Query(tQuery)
	if tError != nil {
		return nil, tError
	}
	defer tRows.Close()
	var tRoles []*Role
	for tRows.Next() {
		var tRole Role
		tError = tRows.Scan(&tRole.ID, &tRole.Name, &tRole.Managed, &tRole.Mentionable, &tRole.Hoist, &tRole.Color, &tRole.Position)
		if tError != nil {
			return nil, tError
		}
		tRoles = append(tRoles, &tRole)
	}
	return tRoles, nil
}
