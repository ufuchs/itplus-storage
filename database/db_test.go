package database

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

/*

{"host":"rigaer28-v3","num":5,"alias":"Wohnzimmer","phenomenontime":1505763276,"lon":5.1,"lat":5.2,"alt":5.3,"temp":18.7,"pressure":56,"humidity":-999,"lowbattery":false}

*/

func Testgateway(t *testing.T) {

	db, err := sql.Open("mysql", "itplus:abc12345@tcp(192.168.178.12:3306)/itp-home")

	defer db.Close()

	if err != nil {
		t.Error(err)
	}

	if err = db.Ping(); err != nil {
		t.Error(err.Error())
	}

	g := NewGatewayDAO(db)

	list, err := g.RetrieveAll(1)
	if err != nil {
		t.Error(err.Error())
	}

	for _, gw := range list {
		fmt.Println(gw)
	}

}
