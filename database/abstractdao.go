package database

import (
	"database/sql"
)

type AbstractDAO struct {
	Db     *sql.DB
	Schema string
}

//
//
//
func NewAbstractDAO(db *sql.DB, schema string) *AbstractDAO {
	return &AbstractDAO{
		Db:     db,
		Schema: schema,
	}
}

//
//
//
func (a *AbstractDAO) CreateTable(sql string) error {

	stmt, err := a.Db.Prepare(sql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	return err

}
