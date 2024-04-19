package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Svoevolin/url-shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/Svoevolin/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty URL",
			alias:     "some_alias",
			url:       "",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			alias:     "some_alias",
			url:       "some invalid URL",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to save url",
			mockError: errors.New("unexpected error"),
		},
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).Once()
			}

			handler := New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "save", bytes.NewReader([]byte(input)))
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
