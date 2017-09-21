package database

import "database/sql"

//
//
//
type MeasurementService struct {
	db   *sql.DB
	gwID int
}

//
//
//
func NewMeasurementService(db *sql.DB, gwID int) *MeasurementService {
	return &MeasurementService{
		db:   db,
		gwID: gwID,
	}
}

//
//
//
func (s *MeasurementService) TableExists(alias string) bool {
	dao := NewMeasurementExDAO(s.db)
	b, _ := dao.TableExists(alias)
	return b
}
