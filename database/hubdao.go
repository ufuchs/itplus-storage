package database

import "database/sql"

type HubDAO struct {
	AbstractDAO
	createStmt string
}

//
// NewHubDAO
//
func NewHubDAO(db *sql.DB, schema string) *HubDAO {

	dao := &HubDAO{
		createStmt: `CREATE TABLE %s (
			HubID       int NOT NULL AUTO_INCREMENT,
			Hostname    varchar(256),
			Alias       varchar(32),
			PRIMARY KEY (HubID)
		)`,
	}

	dao.AbstractDAO.Db = db
	dao.AbstractDAO.Schema = schema

	return dao
}
