package homing

import "github.com/ip2location/ip2location-go/v9"

// LocationService service to obtain location from IP Address.
type LocationService interface {
	// GetLocation from ip address
	GetLocation(ip string) (ip2location.IP2Locationrecord, error)
}

// NewLocationService creates new location service.
func NewLocationService(dbfilepath string) (LocationService, error) {
	db, err := ip2location.OpenDB(dbfilepath)
	if err != nil {
		return nil, err
	}
	return &locationService{
		db: db,
	}, nil

}

var _ LocationService = (*locationService)(nil)

type locationService struct {
	db *ip2location.DB
}

// GetLocation implements LocationService.
func (ls *locationService) GetLocation(ip string) (ip2location.IP2Locationrecord, error) {
	return ls.db.Get_all(ip)
}
