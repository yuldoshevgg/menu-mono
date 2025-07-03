package config

type AuthInfo string

const (
	Authorization AuthInfo = "authorization"
	VenueID       AuthInfo = "x-venue-id"
	TableNumber   AuthInfo = "x-table-number"
)
