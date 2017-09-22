package database

import (
	"database/sql"
	"fmt"
)

//
// LacrosseDAO
//
type LacrosseDAO struct {
	insStmt string
	db      *sql.DB
	schema  string
}

//
//
//
func NewLacrosseDAO(db *sql.DB, schema string) *LacrosseDAO {

	var insStmt = `INSERT INTO %s (
		GatewayID,
		Num, 
		Alias,
		PhenomenonTime,   
		Lon, 
		Lat,  
		Alt,  
		Temp,  
		Pressure,  
		Humidity,  
		LowBattery    
	) VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)
	`
	return &LacrosseDAO{
		insStmt: insStmt,
		db:      db,
		schema:  schema,
	}
}

//
//
//
func (d *LacrosseDAO) CreateTable(fqn string) error {

	var s = `CREATE TABLE %s (
		ID_%s int NOT NULL AUTO_INCREMENT,
		GatewayID int NOT NULL,
		Num            int, 
		Alias          varchar(32),
		PhenomenonTime bigint,   
		Lon            float(18), 
		Lat            float(18),  
		Alt            float(18),  
		Temp           float(18),  
		Pressure       float(18),  
		Humidity       float(18),  
		LowBattery     boolean,    
		PRIMARY KEY (ID_%s),
		CONSTRAINT FK_%s_Gateway FOREIGN KEY (GatewayID) REFERENCES gateway(GatewayID)
	);
	`
	sql := fmt.Sprintf(s, fqn, fqn, fqn, fqn)

	stmt, err := d.db.Prepare(sql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	return err

}

//
// ExistsTable
//
func (d *LacrosseDAO) ExistsTable(fqn string) bool {

	var tname string

	var s = "SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?"

	if err := d.db.QueryRow(s, d.schema, fqn).Scan(&tname); err != nil {
		return false
	}

	return true

}

//
// GetKnownTables
//
func (d *LacrosseDAO) GetKnownTables(tableName string) (res []string, err error) {

	var s = "SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name LIKE ?"
	var rows *sql.Rows

	if rows, err = d.db.Query(s, d.schema, tableName+"%"); err != nil {
		return nil, err
	}

	defer rows.Close()

	var tname string
	for rows.Next() {

		if err = rows.Scan(&tname); err != nil {
			break
		}

		res = append(res, tname)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

//
// {"host":"rigaer28-v3","num":5,"alias":"Wohnzimmer","phenomenontime":1505763276,"lon":5.1,"lat":5.2,"alt":5.3,"temp":18.7,"pressure":56,"humidity":-999,"lowbattery":false}
//
func (d *LacrosseDAO) Insert(gwID int, m *MeasurementEx) error {

	stmt, err := d.db.Prepare(d.insStmt)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(gwID, m.Num, m.Alias, m.PhenomenonTime, m.Lon, m.Lat,
		m.Alt, m.Temp, m.Pressure, m.Humidity, m.LowBattery)
	if err != nil {
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
