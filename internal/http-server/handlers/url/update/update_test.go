package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Svoevolin/url-shortener/internal/database"
	"github.com/Svoevolin/url-shortener/internal/http-server/handlers/url/update/mocks"
	"github.com/Svoevolin/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateHandler(t *testing.T) {

	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "go",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias - will be generated automatically with length of aliasLength",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			alias:     "some alias",
			url:       "",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid url",
			alias:     "some alias",
			url:       "google/com",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "Valid url, but bad url",
			alias:     "some alias",
			url:       `https:/.\214youtube.com`,
			respError: "failed to decode request body",
		},
		{
			name:      "We have no alias to this url",
			alias:     "some alias",
			url:       "https://youtube.com",
			mockError: database.ErrURLNotFound,
			respError: "url not found",
		},
		{
			name:      "UpdateURL error",
			alias:     "some alias",
			url:       "https://youtube.com",
			mockError: fmt.Errorf("unexpected error"),
			respError: "failed to update url",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlUpdaterMock := mocks.NewURLUpdater(t)

			if tc.respError == "" || tc.mockError != nil {
				urlUpdaterMock.On("UpdateURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).Once()
			}

			handler := New(slogdiscard.NewDiscardLogger(), urlUpdaterMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPut, "/url", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
