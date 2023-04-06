package aime

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeRT struct {
	rt func(*http.Request) (*http.Response, error)
}

func (rt *fakeRT) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.rt(req)
}

var _ http.RoundTripper = new(fakeRT)

func withFakeRT(rt func(*http.Request) (*http.Response, error)) ClientOption {
	return func(c *http.Client) error {
		c.Transport = &fakeRT{rt: rt}
		return nil
	}
}

func TestReq(t *testing.T) {
	request, err := ToRequest(&RequestData{Model: "gpt-3.5-turbo"})
	assert.Nil(t, err)
	assert.NotNil(t, request)

	client, err := ToClient(withFakeRT(
		func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"model":"gpt-3.5-turbo"}`)),
			}, nil
		},
	))
	assert.Nil(t, err)
	assert.NotNil(t, client)

	data, err := MakeRequest(client, request)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, data.Model, "gpt-3.5-turbo")
}
