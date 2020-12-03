package roteiro

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-test/deep"

	httpwrapper "github.com/edusalguero/roteiro.git/internal/utils/httpserver"
	"github.com/stretchr/testify/assert"
)

func TestStatusController_Status(t *testing.T) {
	httpServerWrapper := httpwrapper.NewHTTPServerWrapper(httpwrapper.Config{
		Mode: "debug",
		Port: "9092",
	})
	httpServerWrapper.AddController(NewStatusController())

	t.Run("status ok", func(t *testing.T) {
		defer httpServerWrapper.Stop(context.Background())

		httpServerWrapper.Start()
		path := "/status"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, strings.NewReader(""))

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
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func readGoldenJSON(t *testing.T, path string) interface{} {
	t.Helper()

	golden, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden file %q: %v", path, err)
	}

	var goldenJSON interface{}
	err = json.Unmarshal(golden, &goldenJSON)
	if err != nil {
		t.Fatalf("failed to unmarshal golden json: %v", err)
	}

	return goldenJSON
}
