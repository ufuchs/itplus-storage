package database

import (
	"database/sql"
	"testing"
)

func TestMeasurement(t *testing.T) {

	db, err := sql.Open("mysql", "itplus:abc12345@tcp(192.168.178.12:3306)/itp_home")

	defer db.Close()

	if err != nil {
		t.Error(err)
	}

	if err = db.Ping(); err != nil {
		t.Error(err.Error())
	}

	dao := NewMeasurementExDAO(db, "itp_home")

	err = dao.TableExists("hub")

	if err != nil {
		t.Error(err)
	}

}
