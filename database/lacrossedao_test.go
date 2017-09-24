package database

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ufuchs/itplus/base/fcc"
)

func getMeasurement() *MeasurementEx {

	n := fcc.MeasurementDTO{

		Host:           "rigaer28-v3",
		Num:            5,
		Alias:          "Wohnzimmer",
		PhenomenonTime: 1505763276,
		Lon:            5.1,
		Lat:            5.2,
		Alt:            5.3,
		Temp:           18.7,
		Pressure:       56,
		Humidity:       -999,
		LowBattery:     false,
	}

	m := MeasurementEx{
		GatewayID: 1,
	}

	m.MeasurementDTO = n

	return &m

}

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

	fqn := dao.getFQNTablename("Lacrosse", 5)

	fmt.Println(fqn)

	err = dao.Insert(16, "Lacrosse", getMeasurement())

	if err != nil {
		t.Error(err)
	}

}
