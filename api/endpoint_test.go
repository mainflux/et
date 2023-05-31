package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mainflux/callhome"
	"github.com/mainflux/callhome/mocks"
	"github.com/mainflux/mainflux/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEndpointsRetrieve(t *testing.T) {
	svc := mocks.NewService(t)
	svc.On("Retrieve", mock.Anything, callhome.PageMetadata{Limit: 10}).Return(callhome.TelemetryPage{}, nil)
	h := MakeHandler(svc, opentracing.NoopTracer{}, logger.NewMock())
	server := httptest.NewServer(h)
	client := server.Client()
	testCases := []struct {
		test       string
		limit      int
		offset     int
		statuscode int
	}{
		{"successful req", 10, 0, http.StatusOK},
		{"large-limit-size", maxLimitSize + 1, 0, http.StatusBadRequest},
		{"negative-limit-size", -1, 0, http.StatusBadRequest},
	}

	for _, testCase := range testCases {
		t.Run(testCase.test, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/telemetry?limit=%d&offset=%d", server.URL, testCase.limit, testCase.offset), nil)
			assert.Nil(t, err)
			res, err := client.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, testCase.statuscode, res.StatusCode)
		})
	}
}

func TestEndpointSave(t *testing.T) {
	body := `{
		"service": "ty",
		"mainflux_version": "1.0",
		"ip_address": "41.90.185.50",
		"last_seen":"2023-03-27T17:40:50.356401087+03:00"
		}`
	svc := mocks.NewService(t)
	svc.On("Save", mock.Anything, mock.AnythingOfType("callhome.Telemetry")).Return(nil)
	h := MakeHandler(svc, opentracing.NoopTracer{}, logger.NewMock())
	server := httptest.NewServer(h)
	client := server.Client()
	testCases := []struct {
		test, body, contetType string
		statuscode             int
	}{
		{"success", body, "application/json", http.StatusCreated},
		{"malformed-request", "{}", "application/json", http.StatusBadRequest},
		{"wrong-content-type", "{}", "application/text", http.StatusUnsupportedMediaType},
	}

	for _, testCase := range testCases {
		t.Run(testCase.test, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/telemetry", server.URL), strings.NewReader(testCase.body))
			if testCase.contetType != "" {
				req.Header.Set("Content-Type", testCase.contetType)
			}
			assert.Nil(t, err)
			res, err := client.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, testCase.statuscode, res.StatusCode)
		})
	}
}
