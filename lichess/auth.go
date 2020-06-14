package lichess

import (
	"fmt"
	"net/http"
)

func NewAuthorizingTransport(token string, wrapped http.RoundTripper) *AuthorizingTransport {
	return &AuthorizingTransport{
		token:   token,
		wrapped: wrapped,
	}
}

type AuthorizingTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (t *AuthorizingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return t.wrapped.RoundTrip(req)
}
