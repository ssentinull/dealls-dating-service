package common

import "math"

func toRadian(deg float64) float64 {
	return deg * math.Pi / 180
}

type Point struct {
	Lat float64
	Lon float64
}

func (p Point) toRadians() Point {
	return Point{
		Lat: toRadian(p.Lat),
		Lon: toRadian(p.Lon),
	}
}

func (p Point) Delta(point Point) Point {
	return Point{
		Lat: p.Lat - point.Lat,
		Lon: p.Lon - point.Lon,
	}
}

func HaversineDistance(src, dest Point) float64 {
	const earthRadiusMetres float64 = 6371000

	src = src.toRadians()
	dest = dest.toRadians()

	diff := src.Delta(dest)

	a := math.Pow(math.Sin(diff.Lat/2), 2) + math.Cos(src.Lat)*math.Cos(dest.Lat)*math.Pow(math.Sin(diff.Lon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * earthRadiusMetres
}
