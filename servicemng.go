package main

import (
	"context"

	"hidrive.com/ufuchs/itplus/base/zvous"
	"hidrive.com/ufuchs/itplus/storage/socket"
)

type (
	DBService struct {
		ID     int
		client *socket.Client
	}

	//
	Manager struct{}
)

//
//
//
func NewDBService(ctx context.Context, entry zvous.ServiceEntry, ID int) *DBService {

	dbservice := &DBService{
		ID:     ID,
		client: socket.NewClient(entry.ExtractConn(), ID),
	}

	go dbservice.client.Run(ctx)

	return dbservice

}
