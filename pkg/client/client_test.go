package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/mainflux/mainflux/logger"
)

func TestGetIp(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, cancel := context.WithCancel(context.Background())
		hs := New("test_svc", "test.1", logger.NewMock(), cancel)
		for _, endpoint := range ipEndpoints {
			if _, err := hs.getIP(endpoint); err != nil {
				t.Errorf("endpoint ip request unsuccessful with err : %v", err)
			}
		}
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestSend(t *testing.T) {
	t.Run("successful-send", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
				Header:     make(http.Header),
			}
		})
		hs := homingService{
			httpClient: *client,
		}
		data := telemetryData{}

		if err := hs.send(&data); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})
	t.Run("error-sending-req", func(t *testing.T) {
		client := NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString(`some error`)),
				Header:     make(http.Header),
			}
		})
		hs := homingService{
			httpClient: *client,
		}
		data := telemetryData{}

		if err := hs.send(&data); err == nil {
			t.Error("expected non nil error")
		}
	})

}
