package roteiro

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/edusalguero/roteiro.git/internal/logger"
	"github.com/edusalguero/roteiro.git/internal/model"
	"github.com/edusalguero/roteiro.git/internal/point"
	"github.com/edusalguero/roteiro.git/internal/problem"
	"github.com/edusalguero/roteiro.git/internal/solver"
	solverMock "github.com/edusalguero/roteiro.git/internal/solver/mock"
	httpwrapper "github.com/edusalguero/roteiro.git/internal/utils/httpserver"
	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSolverController_solveProblem(t *testing.T) {
	tests := []struct {
		name          string
		idGenerator   func() uuid.UUID
		prepareSolver func(t *testing.T, s *solverMock.MockService)
		statusCode    int
	}{
		{
			"when invalid json request",
			func() uuid.UUID {
				return uuid.MustParse("6e175ad7-7776-4992-94e0-b010589d0772")
			},
			func(t *testing.T, s *solverMock.MockService) {
				// Do nothing
			},
			400,
		},
		{
			"when error solving problem",
			func() uuid.UUID {
				return uuid.MustParse("6e175ad7-7776-4992-94e0-b010589d0772")
			},
			func(t *testing.T, s *solverMock.MockService) {
				s.EXPECT().
					SolveProblem(gomock.Any(), gomock.Any()).
					Return(nil, solver.ErrInAlgo)
			},
			500,
		},
		{
			"ok",
			func() uuid.UUID {
				return uuid.MustParse("83437db4-3e3b-4167-bb7b-74178b6586fd")
			},
			func(t *testing.T, s *solverMock.MockService) {
				s.EXPECT().
					SolveProblem(gomock.Any(), gomock.Any()).
					Return(&problem.Solution{
						ID: uuid.MustParse("83437db4-3e3b-4167-bb7b-74178b6586fd"),
						Solution: model.Solution{
							Metrics: model.SolutionMetrics{
								NumAssets:     1,
								NumRequests:   1,
								NumUnassigned: 0,
								Duration:      0,
								Distance:      0,
								SolvedTime:    161939,
							},
							Routes: []model.SolutionRoute{
								{
									Asset: model.Asset{
										AssetID:  "asset ID",
										Location: point.NewPoint(52.52568, 13.45345),
										Capacity: 1,
									},
									Requests: []model.Request{
										{
											RequestID: "requester ID",
											PickUp:    point.NewPoint(52.52568, 13.45345),
											DropOff:   point.NewPoint(52.52568, 13.45345),
										},
									},
									Waypoints: []model.Waypoint{
										{
											Location: point.NewPoint(52.52568, 13.45345),
											Load:     0,
											Activities: []model.Activity{
												{
													ActivityType: model.ActivityTypeStart,
													Ref:          "asset ID",
												},
												{
													ActivityType: model.ActivityTypePickUp,
													Ref:          "requester ID",
												},
												{
													ActivityType: model.ActivityTypeDropOff,
													Ref:          "requester ID",
												},
											},
										},
									},
									Metrics: model.RouteMetrics{},
								},
							},
							Unassigned: []model.Request{},
						},
					}, nil)
			},
			200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httpServerWrapper := httpwrapper.NewHTTPServerWrapper(httpwrapper.Config{
				Mode: "debug",
				Port: "9092",
			})
			defer httpServerWrapper.Stop(context.Background())
			log := logger.NewNopLogger()

			s := solverMock.NewMockService(ctrl)
			tt.prepareSolver(t, s)
			httpServerWrapper.AddController(NewSolverController(s, log, tt.idGenerator))

			w := httptest.NewRecorder()
			path := "/api/v1/problem"
			reqPath := filepath.Join("./testdata", t.Name()+".req.json")
			r := readRequestJSON(t, reqPath)
			req, _ := http.NewRequest("POST", path, bytes.NewReader(r))
			httpServerWrapper.GetGin().ServeHTTP(w, req)

			var resData interface{}
			_ = json.NewDecoder(w.Body).Decode(&resData)
			_, err := json.MarshalIndent(&resData, "", "    ")
			if err != nil {
				t.Fatalf("failed to unmarshal json response: %v", err)
			}
			goldenPath := filepath.Join("./testdata", t.Name()+".golden.json")

			goldenJSON := readGoldenJSON(t, goldenPath)
			differences := deep.Equal(goldenJSON, resData)
			if differences != nil {
				t.Errorf("response not matching golden file: %v", differences)
			}
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func readRequestJSON(t *testing.T, path string) []byte {
	t.Helper()

	req, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden file %q: %v", path, err)
	}

	return req
}
