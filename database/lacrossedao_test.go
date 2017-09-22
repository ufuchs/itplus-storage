package database

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestLacrosse(t *testing.T) {

	db, err := sql.Open("mysql", "itplus:abc12345@tcp(192.168.178.12:3306)/itp_home")

	defer db.Close()

	if err != nil {
		t.Error(err)
	}

	if err = db.Ping(); err != nil {
		t.Error(err.Error())
	}

	dao := NewLacrosseDAO(db, "itp_home")

	var tnames []string
	tnames, err = dao.GetKnownTables("Lacrosse_")

	for _, tn := range tnames {
		fmt.Println(tn)
	}

	// if len(tnames) == 0 {
	// 	fqn := fmt.Sprintf("%s_%0.2d", "Lacrosse", 6)
	// 	err = dao.CreateTable(fqn)
	// }

	if err != nil {
		t.Error(err)
	}

	exists := dao.ExistsTable("Lacrosse_06")

	fmt.Println(exists)

}
