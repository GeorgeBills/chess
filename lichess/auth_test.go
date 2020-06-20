package lichess_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/GeorgeBills/chess/lichess"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	last *http.Request
}

func (mrt *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	mrt.last = req
	return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
}

func TestRoundTrip(t *testing.T) {
	mrt := &mockRoundTripper{}
	at := lichess.NewAuthorizingTransport("s3cr3t_t0k3n", mrt)
	req, err := http.NewRequest("GET", "https://lichess.org/api/account", http.NoBody)
	require.NoError(t, err)
	res, err := at.RoundTrip(req)
	require.NoError(t, err)

	expectedResponse := &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}
	assert.Equal(t, expectedResponse, res, "response should be passed through unchanged")

	expectedRequest := &http.Request{
		Method:     http.MethodGet,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       http.NoBody,
		Host:       "lichess.org",
		URL:        &url.URL{Scheme: "https", Host: "lichess.org", Path: "/api/account"},
		Header:     http.Header{"Authorization": {"Bearer s3cr3t_t0k3n"}},
	}
	expectedRequest = expectedRequest.WithContext(context.Background())
	assert.Equal(t, expectedRequest, mrt.last, "request should contain authorisation header")
}
