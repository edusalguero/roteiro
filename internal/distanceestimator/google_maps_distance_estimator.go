package distanceestimator

import (
	"context"
	"fmt"

	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/point"
	"googlemaps.github.io/maps"
)

type GoogleMapsDistanceEstimator struct {
	client *maps.Client
	logger logger.Logger
	cache  map[string]RouteEstimation
}

func NewGoogleMapsDistanceEstimator(apiKey string, l logger.Logger) (*GoogleMapsDistanceEstimator, error) {
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GoogleMapsDistanceEstimator{client: c, logger: l, cache: make(map[string]RouteEstimation)}, nil
}

func (g *GoogleMapsDistanceEstimator) EstimateDistance(ctx context.Context, from, to point.Point) (*RouteEstimation, error) {
	if e, ok := g.cache[fmt.Sprintf("%s-%s", from, to)]; ok {
		return &e, nil
	}

	if from.Equal(to) {
		return &RouteEstimation{
			From:     from,
			To:       to,
			Distance: float64(0),
			Duration: 0,
		}, nil
	}

	request := maps.DistanceMatrixRequest{
		Origins:       []string{fmt.Sprintf("%v,%v", from.Lat(), from.Lon())},
		Destinations:  []string{fmt.Sprintf("%v,%v", to.Lat(), to.Lon())},
		DepartureTime: `now`,
		Units:         maps.UnitsMetric,
		Mode:          maps.TravelModeDriving,
	}
	g.logger.Debugf("Requesting estimation %+v", request)
	resp, err := g.client.DistanceMatrix(ctx, &request)

	if err != nil {
		return nil, err
	}

	if len(resp.Rows) == 0 {
		return nil, fmt.Errorf("no response rows")
	}
	r := resp.Rows[0].Elements[0]

	if r.Status != "OK" {
		return nil, fmt.Errorf("invalid response element status")
	}

	e := RouteEstimation{
		From:     from,
		To:       to,
		Distance: float64(r.Distance.Meters),
		Duration: r.Duration,
	}
	g.cache[fmt.Sprintf("%s-%s", from, to)] = e
	return &e, nil
}
