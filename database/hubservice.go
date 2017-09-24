package database

type HubService struct {
	hub []Hub
}

//
//
//
func NewHubService() *HubService {
	s := &HubService{
		hub: []Hub{},
	}
	return s
}

//
//
//
func (s *HubService) Retrieve() error {
	return nil
}
