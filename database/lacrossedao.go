package database

import (
	"database/sql"
	"fmt"
)

//
// LacrosseDAO
//
type LacrosseDAO struct {
	AbstractDAO
	insStmt string
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
	dao := &LacrosseDAO{
		insStmt: insStmt,
	}

	dao.AbstractDAO.Db = db
	dao.AbstractDAO.Schema = schema

	return dao

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

	stmt, err := d.Db.Prepare(sql)
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

	return d.Db.QueryRow(s, d.Schema, fqn).Scan(&tname) == nil

}

//
// GetKnownTables
//
func (d *LacrosseDAO) GetKnownTables(tableName string) (res []string, err error) {

	var s = "SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name LIKE ?"
	var rows *sql.Rows

	if rows, err = d.Db.Query(s, d.Schema, tableName+"%"); err != nil {
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
//
//
func (d *LacrosseDAO) getFQNTablename(basename string, num int) string {
	return fmt.Sprintf("%s_%0.2d", "Lacrosse", 6)
}

//
//
//
func (d *LacrosseDAO) Insert(gwID int, basename string, m *MeasurementEx) error {

	fqn := d.getFQNTablename(basename, m.Num)

	insStmt := fmt.Sprintf(d.insStmt, fqn)

	stmt, err := d.Db.Prepare(insStmt)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(gwID, m.Num, m.Alias, m.PhenomenonTime, m.Lon, m.Lat,
		m.Alt, m.Temp, m.Pressure, m.Humidity, m.LowBattery)

	return err
}
