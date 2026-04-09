package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zcl0621/compx576-smart-dairy-system/util"
)

func TestDouglasPeucker_StraightLine(t *testing.T) {
	points := []util.GeoPoint{
		{Lat: 0, Lng: 0},
		{Lat: 1, Lng: 1},
		{Lat: 2, Lng: 2},
		{Lat: 3, Lng: 3},
		{Lat: 4, Lng: 4},
	}
	result := util.DouglasPeucker(points, 0.1)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, points[0], result[0])
	assert.Equal(t, points[4], result[1])
}

func TestDouglasPeucker_Zigzag(t *testing.T) {
	points := []util.GeoPoint{
		{Lat: 0, Lng: 0},
		{Lat: 1, Lng: 5},
		{Lat: 2, Lng: 0},
		{Lat: 3, Lng: 5},
		{Lat: 4, Lng: 0},
	}
	result := util.DouglasPeucker(points, 0.1)
	assert.Equal(t, 5, len(result))
}

func TestDouglasPeucker_SinglePoint(t *testing.T) {
	points := []util.GeoPoint{{Lat: 1, Lng: 2}}
	result := util.DouglasPeucker(points, 0.1)
	assert.Equal(t, 1, len(result))
}

func TestDouglasPeucker_TwoPoints(t *testing.T) {
	points := []util.GeoPoint{{Lat: 0, Lng: 0}, {Lat: 1, Lng: 1}}
	result := util.DouglasPeucker(points, 0.1)
	assert.Equal(t, 2, len(result))
}

func TestDouglasPeucker_Empty(t *testing.T) {
	result := util.DouglasPeucker(nil, 0.1)
	assert.Equal(t, 0, len(result))
}

func TestDouglasPeucker_KeepMarkedPoints(t *testing.T) {
	points := []util.GeoPoint{
		{Lat: 0, Lng: 0},
		{Lat: 1, Lng: 1, Keep: true},
		{Lat: 2, Lng: 2},
		{Lat: 3, Lng: 3, Keep: true},
		{Lat: 4, Lng: 4},
	}
	result := util.DouglasPeucker(points, 0.1)
	assert.Equal(t, 4, len(result))
	assert.True(t, result[1].Keep)
	assert.True(t, result[2].Keep)
}
