package database

import "database/sql"

//
//
//
type MeasurementService struct {
	db     *sql.DB
	gwID   int
	schema string
}

//
//
//
func NewMeasurementService(db *sql.DB, gwID int, schema string) *MeasurementService {
	return &MeasurementService{
		db:     db,
		gwID:   gwID,
		schema: schema,
	}
}

//
//
//
func (s *MeasurementService) TableExists(alias string) error {
	dao := NewMeasurementExDAO(s.db, s.schema)
	return dao.TableExists(alias)
}
