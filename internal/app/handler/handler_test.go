package handler

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type usecaseMock struct {
	s string // short
	o string // orig
	e bool   // err
}

func (u *usecaseMock) Shorten(_ string) string { return u.s }
func (u *usecaseMock) RestoreOrigin(_ string) (string, error) {
	if u.e {
		return "", errors.New("usecase error")
	}
	return u.o, nil
}

func TestAppHandler_HandleMain(t *testing.T) {
	t.Run("Test Handler", func(t *testing.T) {

		// Prepare fake usecase
		uc := &usecaseMock{
			s: "xyz",
			o: "http://example.com",
			e: false,
		}
		// Main App router
		h := NewAppRouter(uc)

		// POST request
		body := bytes.NewBufferString("http://example.com")
		request := httptest.NewRequest(http.MethodPost, "/", body)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, request)
		response := w.Result()

		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))

		content, err := ioutil.ReadAll(response.Body)
		require.NoError(t, err)

		assert.Equal(t, "http://localhost:8080/xyz", string(content))

		err = response.Body.Close()
		require.NoError(t, err)

		// GET request
		request = httptest.NewRequest(http.MethodGet, "/xyz", nil)
		w = httptest.NewRecorder()

		h.ServeHTTP(w, request)
		response = w.Result()

		assert.Equal(t, 307, response.StatusCode)
		assert.Equal(t, "http://example.com", response.Header.Get("Location"))

		err = response.Body.Close()
		require.NoError(t, err)

		// GET Not found
		request = httptest.NewRequest(http.MethodGet, "/abc", nil)
		w = httptest.NewRecorder()

		uc.e = true // Make usecase to return an error

		h.ServeHTTP(w, request)
		response = w.Result()

		assert.Equal(t, 404, response.StatusCode)

		err = response.Body.Close()
		require.NoError(t, err)

	})
}

func TestAppHandler_ApiJson_ShortURL(t *testing.T) {
	t.Run("Test Handler API Json ShortURL", func(t *testing.T) {

		// Prepare fake usecase
		uc := &usecaseMock{
			s: "xyz",
			o: "http://example.com",
			e: false,
		}

		// Main App router
		h := NewAppRouter(uc)

		// OK
		// Prepare request json
		requestJSON := "{ \"url\" : \"http://example.com\"}"
		body := bytes.NewBufferString(requestJSON)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", body)

		w := httptest.NewRecorder()
		h.ServeHTTP(w, request)
		response := w.Result()

		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, "application/json", response.Header.Get("Content-Type"))

		content, err := ioutil.ReadAll(response.Body)
		require.NoError(t, err)

		assert.Equal(t, "{\"result\":\"http://localhost:8080/xyz\"}", string(content))

		err = response.Body.Close()
		require.NoError(t, err)

		// Error - Invalid Json
		// Prepare bad json
		requestBadJSON := "{ \"url\" : http://example.com\"}"
		body = bytes.NewBufferString(requestBadJSON)
		request = httptest.NewRequest(http.MethodPost, "/api/shorten", body)

		w = httptest.NewRecorder()
		h.ServeHTTP(w, request)

		response = w.Result()
		assert.Equal(t, 400, response.StatusCode)

		err = response.Body.Close()
		require.NoError(t, err)

		// Error - wrong json
		requestWrongJSON := "{ \"urlBad\" : http://example.com\"}"
		body = bytes.NewBufferString(requestWrongJSON)
		request = httptest.NewRequest(http.MethodPost, "/api/shorten", body)

		w = httptest.NewRecorder()
		h.ServeHTTP(w, request)

		response = w.Result()
		assert.Equal(t, 400, response.StatusCode)

		err = response.Body.Close()
		require.NoError(t, err)

	})
}
