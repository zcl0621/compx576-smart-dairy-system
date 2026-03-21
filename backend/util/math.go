package util

import "math"

// convert degrees to radians
func DegreesToRadians(v float64) float64 {
	return v * math.Pi / 180
}

// calc distance in meters from 2 lat lng points
func HaversineMeters(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000.0
	dLat := DegreesToRadians(lat2 - lat1)
	dLng := DegreesToRadians(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(DegreesToRadians(lat1))*math.Cos(DegreesToRadians(lat2))*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
