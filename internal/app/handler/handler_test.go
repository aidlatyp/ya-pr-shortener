package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type usecaseMock struct {
	s string
	o string
	e bool
}

func (u *usecaseMock) Shorten(_ string) string {
	return u.s
}
func (u *usecaseMock) RestoreOrigin(_ string) (string, error) {
	if u.e {
		fmt.Println("error zz")
		return "", errors.New("usecase error")
	}
	return u.o, nil
}

func TestAppHandler_HandleMain(t *testing.T) {
	t.Run("Test Handler", func(t *testing.T) {

		// POST
		body := bytes.NewBufferString("http://example.com")
		request := httptest.NewRequest(http.MethodPost, "/", body)

		w := httptest.NewRecorder()

		uc := &usecaseMock{
			s: "xyz",
			o: "http://example.com",
			e: false,
		}

		appHandler := NewAppHandler(uc)
		h := appHandler.HandleMain()

		h.ServeHTTP(w, request)
		response := w.Result()

		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))

		content, err := ioutil.ReadAll(response.Body)
		require.NoError(t, err)

		assert.Equal(t, "http://localhost:8080/xyz", string(content))

		err = response.Body.Close()
		require.NoError(t, err)

		// GET
		request = httptest.NewRequest(http.MethodGet, "/xyz", nil)
		w = httptest.NewRecorder()

		h.ServeHTTP(w, request)
		response = w.Result()

		assert.Equal(t, 307, response.StatusCode)
		assert.Equal(t, "http://example.com", response.Header.Get("Location"))

		err = response.Body.Close()
		require.NoError(t, err)

		// Not found
		request = httptest.NewRequest(http.MethodGet, "/abc", nil)
		w = httptest.NewRecorder()

		uc.e = true

		h.ServeHTTP(w, request)
		response = w.Result()

		assert.Equal(t, 404, response.StatusCode)

		err = response.Body.Close()
		require.NoError(t, err)

	})
}
