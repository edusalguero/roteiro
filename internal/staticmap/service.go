package staticmap

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/edusalguero/roteiro.git/internal/point"
	maps "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"

	"github.com/edusalguero/roteiro.git/internal/problem"
)

type Service interface {
	Render(solution problem.Solution) error
}

type StaticMap struct {
}

func NewService() *StaticMap {
	return &StaticMap{}
}

func (s StaticMap) Render(solution *problem.Solution) (image.Image, error) {
	mapCtx := maps.NewContext()
	mapCtx.SetSize(2000, 2000)
	for _, r := range solution.Routes {
		mapCtx.AddMarker(createMarker(r.Asset.Location, string(r.Asset.AssetID), randomColor(), 16))
		for _, req := range r.Requests {
			c := randomColor()
			mapCtx.AddMarker(createMarker(req.PickUp, string(req.RequestID), c, 16))
			mapCtx.AddMarker(createMarker(req.DropOff, string(req.RequestID), c, 16))
		}
		var positions []s2.LatLng
		for _, w := range r.Waypoints {
			positions = append(positions, s2PointFromPoint(w.Location))
		}
		mapCtx.AddPath(maps.NewPath(positions, randomColor(), 2))
	}

	for _, req := range solution.Unassigned {
		c := randomColor()
		mapCtx.AddMarker(createMarker(req.PickUp, string(req.RequestID), c, 10))
		mapCtx.AddMarker(createMarker(req.DropOff, string(req.RequestID), c, 10))
	}

	img, err := mapCtx.Render()
	if err != nil {
		return nil, err
	}
	return img, nil
}

func createMarker(p point.Point, _ string, c color.RGBA, size float64) *maps.Marker {
	m := maps.NewMarker(s2PointFromPoint(p), c, size)
	return m
}

func s2PointFromPoint(p point.Point) s2.LatLng {
	return s2.LatLngFromDegrees(p.Lat(), p.Lon())
}

func randomColor() color.RGBA {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return color.RGBA{R: uint8(r.Intn(255)), G: uint8(r.Intn(255)), B: uint8(r.Intn(255)), A: 255}
}
