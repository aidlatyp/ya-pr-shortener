package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aidlatyp/ya-pr-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewRouter() (chi.Router, error) {
	r := chi.NewRouter()

	configs := &utils.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
	}

	handler, _, err := NewHandler(configs)
	if err != nil {
		return nil, err
	}

	r.Route("/", func(r chi.Router) {
		r.Post("/api/shorten", handler.PostShortenURLHandler)
		r.Get("/{id}", handler.GetURLHandler)
		r.Post("/", handler.SaveURLHandler)
	})
	return r, nil
}

func TestPostShortenURLHandlerHandler(t *testing.T) {

	type want struct {
		statusCode   int
		response     PostURLResponse
		contentType  string
		errorMessage string
		err          string
	}

	tests := []struct {
		name        string
		request     string
		body        PostURLRequest
		requestType string
		want        want
	}{
		{
			name:    "simple positive test #1",
			request: "/api/shorten",
			body:    PostURLRequest{URL: "https://ya.ru"},
			want: want{
				statusCode:  http.StatusCreated,
				response:    PostURLResponse{Result: "http://localhost:8080/rfBd67"},
				contentType: "application/json",
			},
			requestType: http.MethodPost,
		},
		{
			name:    "simple test #2 with empty URL",
			request: "/api/shorten",
			body:    PostURLRequest{URL: "d"},
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				err:         "invalid character 'u' looking for beginning of value",
			},
			requestType: http.MethodPost,
		},
		{
			name:    "simple test #3 with uncorrect request type",
			request: "/api/shorten",
			body:    PostURLRequest{URL: "https://ya.ru"},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				err:        "unexpected end of JSON input",
			},
			requestType: http.MethodDelete,
		},
		{
			name:    "simple test #4 with uncorrect url format",
			request: "/api/shorten",
			body:    PostURLRequest{URL: "/ya.ru"},
			want: want{
				statusCode:  http.StatusBadRequest,
				err:         "invalid character 'u' looking for beginning of value",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r, err := NewRouter()
			require.NoError(t, err)

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := json.Marshal(tt.body)
			require.NoError(t, err)

			resp, body := testRequest(t, ts, tt.requestType, tt.request, bytes.NewReader(req))

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			var postURLResponse PostURLResponse
			err = json.Unmarshal([]byte(body), &postURLResponse)

			if err != nil {
				fmt.Println(err)
				assert.Equal(t, tt.want.err, err.Error())
			}

			assert.Equal(t, tt.want.response, postURLResponse)

			resp.Body.Close()
		})
	}
}

func TestSaveURLHandler(t *testing.T) {

	type want struct {
		statusCode  int
		redirectURL string
		contentType string
	}

	tests := []struct {
		name        string
		request     string
		body        string
		requestType string
		want        want
	}{
		{
			name:    "simple positive test #1",
			request: "/",
			body:    "https://ya.ru",
			want: want{
				statusCode:  http.StatusCreated,
				redirectURL: "http://localhost:8080/ti3SMt",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodPost,
		},
		{
			name:    "simple test #2 with empty URL",
			request: "/",
			body:    "",
			want: want{
				statusCode:  http.StatusBadRequest,
				redirectURL: "uncorrect URL format\n",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodPost,
		},
		{
			name:    "simple test #3 with uncorrect request type",
			request: "/",
			body:    "https://ya.ru",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				redirectURL: "",
				contentType: "",
			},
			requestType: http.MethodDelete,
		},
		{
			name:    "simple test #4 with uncorrect url format",
			request: "/",
			body:    "/ya.ru",
			want: want{
				statusCode:  http.StatusBadRequest,
				redirectURL: "uncorrect URL format\n",
				contentType: "text/plain; charset=utf-8",
			},
			requestType: http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r, _ := NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, body := testRequest(t, ts, tt.requestType, tt.request, strings.NewReader(tt.body))

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.redirectURL, body)

			resp.Body.Close()
		})
	}
}

func TestGetURLHandler(t *testing.T) {

	type want struct {
		statusCode  int
		redirectURL string
		contentType string
		body        string
	}

	tests := []struct {
		name        string
		request     string
		requestType string
		want        want
	}{
		//{
		//	name:    "simple positive test #1",
		//	request: "/test",
		//	want: want{
		//		statusCode:  http.StatusTemporaryRedirect,
		//		redirectURL: "https://yatest.ru",
		//		contentType: "text/plain; charset=utf-8",
		//	},
		//	requestType: http.MethodGet,
		//},
		{
			name:    "negative test #2 with wrong method type",
			request: "/test",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				redirectURL: "",
				contentType: "",
			},
			requestType: http.MethodDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r, _ := NewRouter()
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, body := testRequest(t, ts, tt.requestType, tt.request, nil)

			fmt.Println(body)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.redirectURL, resp.Header.Get("Location"))

			resp.Body.Close()
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, body)

	require.NoError(t, err)

	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := httpClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
