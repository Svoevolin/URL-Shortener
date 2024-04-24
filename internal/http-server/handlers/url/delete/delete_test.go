package delete

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Svoevolin/url-shortener/internal/database"
	"github.com/Svoevolin/url-shortener/internal/http-server/handlers/url/delete/mocks"
	"github.com/Svoevolin/url-shortener/internal/lib/api"
	"github.com/Svoevolin/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {

	tests := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "go",
		},
		{
			name:      "non-existent alias deleting",
			alias:     "non-existent",
			mockError: database.ErrURLNotFound,
			respError: "url not found",
		},
		{
			name:      "DeleteURL error",
			alias:     "some alias",
			mockError: fmt.Errorf("unexpected error"),
			respError: "failed to delete url",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).Return(tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Delete("/{alias}", New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			_, body := api.TestRequest(t, ts, "DELETE", fmt.Sprintf("/%s", tc.alias), nil)

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
