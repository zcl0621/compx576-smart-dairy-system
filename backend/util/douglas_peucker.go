package util

import "math"

// GeoPoint is a lat/lng point for Douglas-Peucker simplification
// Keep=true means this point is never removed (used for stay points)
type GeoPoint struct {
	Lat      float64
	Lng      float64
	Keep     bool
	OrigIdx  int // original index, used to reconstruct after simplification
}

// DouglasPeucker simplifies a polyline by removing points closer
// than epsilon to the line between neighbors. Points with Keep=true
// are always retained.
func DouglasPeucker(points []GeoPoint, epsilon float64) []GeoPoint {
	if len(points) <= 2 {
		return points
	}
	return dpRecurse(points, epsilon)
}

func dpRecurse(points []GeoPoint, epsilon float64) []GeoPoint {
	if len(points) <= 2 {
		return points
	}

	maxDist := 0.0
	maxIdx := 0
	first := points[0]
	last := points[len(points)-1]

	for i := 1; i < len(points)-1; i++ {
		d := perpendicularDist(points[i], first, last)
		if points[i].Keep && d <= epsilon {
			d = epsilon + 1
		}
		if d > maxDist {
			maxDist = d
			maxIdx = i
		}
	}

	if maxDist <= epsilon {
		for i := 1; i < len(points)-1; i++ {
			if points[i].Keep {
				left := dpRecurse(points[:i+1], epsilon)
				right := dpRecurse(points[i:], epsilon)
				return append(left, right[1:]...)
			}
		}
		return []GeoPoint{first, last}
	}

	left := dpRecurse(points[:maxIdx+1], epsilon)
	right := dpRecurse(points[maxIdx:], epsilon)
	return append(left, right[1:]...)
}

// perpendicular distance from point to line(a, b) in degrees
func perpendicularDist(p, a, b GeoPoint) float64 {
	dx := b.Lng - a.Lng
	dy := b.Lat - a.Lat
	if dx == 0 && dy == 0 {
		return math.Sqrt((p.Lng-a.Lng)*(p.Lng-a.Lng) + (p.Lat-a.Lat)*(p.Lat-a.Lat))
	}
	return math.Abs(dy*(p.Lng-a.Lng)-dx*(p.Lat-a.Lat)) / math.Sqrt(dx*dx+dy*dy)
}
