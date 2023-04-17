package homing

import "github.com/ip2location/ip2location-go/v9"

var _ LocationService = (*locationService)(nil)

type locationService struct {
	db *ip2location.DB
}

// LocationService provides a service for obtaining location information from an IP address.
type LocationService interface {
	// GetLocation returns the location information for a given IP address.
	GetLocation(ip string) (ip2location.IP2Locationrecord, error)
}

// NewLocationService creates a new LocationService that uses the specified IP2Location database file.
func NewLocationService(dbfilepath string) (LocationService, error) {
	db, err := ip2location.OpenDB(dbfilepath)
	if err != nil {
		return nil, err
	}
	return &locationService{
		db: db,
	}, nil

}

// GetLocation returns the location information for a given IP address.
func (ls *locationService) GetLocation(ip string) (ip2location.IP2Locationrecord, error) {
	return ls.db.Get_all(ip)
}
