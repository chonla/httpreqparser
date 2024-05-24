package httpreqparser_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/chonla/httpreqparser"
	"github.com/stretchr/testify/assert"
)

func TestParseEmptyString(t *testing.T) {
	req := ""
	result, err := httpreqparser.Parse(req)

	assert.Equal(t, errors.New("unexpected http request"), err)
	assert.Nil(t, result)
}

func TestParseGetWithImplicitHost(t *testing.T) {
	req := `GET http://localhost:1234/path HTTP/1.0`
	result, err := httpreqparser.Parse(req)

	assert.Nil(t, err)
	assert.Equal(t, http.MethodGet, result.Method)
	assert.Equal(t, "http://localhost:1234/path", result.URL.String())
}

func TestParseGetWithExplicitHostInHeader(t *testing.T) {
	req := `GET /path HTTP/1.0
Host: localhost:1234`
	result, err := httpreqparser.Parse(req)

	assert.Nil(t, err)
	assert.Equal(t, http.MethodGet, result.Method)
	assert.Equal(t, "http://localhost:1234/path", result.URL.String())
}

func TestParseHTTPSPostWithImplicitHost(t *testing.T) {
	req := `POST https://localhost:1234/path HTTP/1.0
Content-Type: application/json
Content-Length: 15

{"key":"value"}`
	result, err := httpreqparser.Parse(req)

	expectedBody := []byte(`{"key":"value"}`)
	defer result.Body.Close()
	reqBody, _ := io.ReadAll(result.Body)

	assert.Nil(t, err)
	assert.Equal(t, http.MethodPost, result.Method)
	assert.Equal(t, "https://localhost:1234/path", result.URL.String())
	assert.Equal(t, expectedBody, reqBody)
}

func TestParseHTTPSPutWithImplicitHost(t *testing.T) {
	req := `PUT https://localhost:1234/path HTTP/1.0
Content-Type: application/json
Content-Length: 15

{"key":"value"}`
	result, err := httpreqparser.Parse(req)

	expectedBody := []byte(`{"key":"value"}`)
	defer result.Body.Close()
	reqBody, _ := io.ReadAll(result.Body)

	assert.Nil(t, err)
	assert.Equal(t, http.MethodPut, result.Method)
	assert.Equal(t, "https://localhost:1234/path", result.URL.String())
	assert.Equal(t, expectedBody, reqBody)
}

func TestParseHTTPSHEADWithImplicitHost(t *testing.T) {
	req := `HEAD https://localhost:1234/path HTTP/1.0`
	result, err := httpreqparser.Parse(req)

	assert.Nil(t, err)
	assert.Equal(t, http.MethodHead, result.Method)
	assert.Equal(t, "https://localhost:1234/path", result.URL.String())
}

func TestParseHTTPSOPTIONSWithImplicitHost(t *testing.T) {
	req := `OPTIONS https://localhost:1234/path HTTP/1.0`
	result, err := httpreqparser.Parse(req)

	assert.Nil(t, err)
	assert.Equal(t, http.MethodOptions, result.Method)
	assert.Equal(t, "https://localhost:1234/path", result.URL.String())
}

func TestParseHTTPSTRACEWithImplicitHost(t *testing.T) {
	req := `TRACE https://localhost:1234/path HTTP/1.0`
	result, err := httpreqparser.Parse(req)

	assert.Nil(t, err)
	assert.Equal(t, http.MethodTrace, result.Method)
	assert.Equal(t, "https://localhost:1234/path", result.URL.String())
}
