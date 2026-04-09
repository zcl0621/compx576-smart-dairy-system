package metric

import (
	cowdto "github.com/zcl0621/compx576-smart-dairy-system/dto/cow"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/util"
	"go.uber.org/zap"
)

const stayDistanceThreshold = 5.0 // meters
const maxPathPoints = 500

// MetricQuery is the input for movement path queries
type MetricQuery struct {
	CowID       string
	MetricRange model.MetricRange
}

// epsilon per range in degrees
var dpEpsilon = map[model.MetricRange]float64{
	model.MetricRange24H: 0.00002,
	model.MetricRange7D:  0.00005,
	model.MetricRange30D: 0.0001,
	model.MetricRangeAll: 0.0002,
}

// MovementPathService returns GPS path points with stay detection
// and Douglas-Peucker simplification
func MovementPathService(q *MetricQuery) (*cowdto.MovementPathResponse, error) {
	rows, err := loadMovementRows(q.CowID, q.MetricRange)
	if err != nil {
		return nil, err
	}

	rawPoints := buildMovementPoints(q.CowID, rows)

	if len(rawPoints) == 0 {
		return &cowdto.MovementPathResponse{
			CowID:  q.CowID,
			Range:  q.MetricRange,
			Points: []cowdto.MovementPathPoint{},
		}, nil
	}

	// detect stays and build path points
	pathPoints := detectStays(rawPoints)

	// Douglas-Peucker simplification
	epsilon := dpEpsilon[q.MetricRange]
	if epsilon == 0 {
		epsilon = 0.00005
	}

	pathPoints = simplifyPath(pathPoints, epsilon)

	projectlog.L().Debug("movement path built",
		zap.String("cow_id", q.CowID),
		zap.Int("raw", len(rawPoints)),
		zap.Int("after_simplify", len(pathPoints)),
	)

	return &cowdto.MovementPathResponse{
		CowID:  q.CowID,
		Range:  q.MetricRange,
		Points: pathPoints,
	}, nil
}

// detectStays merges consecutive points within 5m into one stay point
func detectStays(points []movementPointRow) []cowdto.MovementPathPoint {
	if len(points) == 0 {
		return nil
	}

	result := make([]cowdto.MovementPathPoint, 0, len(points))
	current := cowdto.MovementPathPoint{
		Lat:         points[0].Latitude,
		Lng:         points[0].Longitude,
		Time:        points[0].CreatedAt.Unix(),
		StaySeconds: 0,
	}
	stayStart := points[0].CreatedAt

	for i := 1; i < len(points); i++ {
		dist := util.HaversineMeters(
			current.Lat, current.Lng,
			points[i].Latitude, points[i].Longitude,
		)

		if dist < stayDistanceThreshold {
			// still in stay zone, accumulate time
			current.StaySeconds = points[i].CreatedAt.Unix() - stayStart.Unix()
		} else {
			// moved away, flush current point and start new
			result = append(result, current)
			current = cowdto.MovementPathPoint{
				Lat:         points[i].Latitude,
				Lng:         points[i].Longitude,
				Time:        points[i].CreatedAt.Unix(),
				StaySeconds: 0,
			}
			stayStart = points[i].CreatedAt
		}
	}

	// flush last point
	result = append(result, current)
	return result
}

// simplifyPath applies Douglas-Peucker, keeping stay points
func simplifyPath(points []cowdto.MovementPathPoint, epsilon float64) []cowdto.MovementPathPoint {
	if len(points) <= 2 {
		return points
	}

	// convert to GeoPoints for DP, track original index
	geoPoints := make([]util.GeoPoint, len(points))
	for i, p := range points {
		geoPoints[i] = util.GeoPoint{
			Lat:     p.Lat,
			Lng:     p.Lng,
			Keep:    p.StaySeconds > 0,
			OrigIdx: i,
		}
	}

	// run DP, double epsilon if too many points (cap iterations to prevent infinite loop)
	for range 20 {
		simplified := util.DouglasPeucker(geoPoints, epsilon)
		if len(simplified) <= maxPathPoints {
			geoPoints = simplified
			break
		}
		epsilon *= 2
	}

	// rebuild path points using original index
	result := make([]cowdto.MovementPathPoint, 0, len(geoPoints))
	for _, gp := range geoPoints {
		result = append(result, points[gp.OrigIdx])
	}

	return result
}
