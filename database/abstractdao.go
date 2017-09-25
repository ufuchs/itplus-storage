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

	tx, err := a.Db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}
