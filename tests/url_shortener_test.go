//go:build functional

package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/Svoevolin/url-shortener/internal/http-server/handlers/url/save"
	"github.com/Svoevolin/url-shortener/internal/lib/api"
	"github.com/Svoevolin/url-shortener/internal/lib/random"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const host = "localhost:8080"

func TestURLShortenerHappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomStrings(10),
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().Status(http.StatusOK).JSON().Object().
		ContainsKey("alias")
}

func TestSaveToRedirectToUpdateToRedirectToDeleteToNoRedirect(t *testing.T) {
	testCases := []struct {
		name         string
		url          string
		alias        string
		updatedAlias string
		error        string
	}{
		{
			name:         "All steps successfully with given aliases",
			url:          gofakeit.URL(),
			alias:        gofakeit.Word() + gofakeit.Word(),
			updatedAlias: gofakeit.Word(),
		},
		{
			name: "All steps successfully: aliases for save and update will be generated automatically",
			url:  gofakeit.URL(),
		},
		{
			name:  "Invalid url, dropped on saving",
			url:   "invalid",
			error: "field URL is not a valid URL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// Redirect

			testRedirect(t, alias, tc.url)

			// Update

			resp = e.PUT("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.updatedAlias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			updatedAlias := tc.updatedAlias

			if tc.updatedAlias != "" {
				resp.Value("alias").String().IsEqual(tc.updatedAlias)
			} else {
				resp.Value("alias").String().NotEmpty()

				updatedAlias = resp.Value("alias").String().Raw()
			}

			// Redirect

			testRedirectNotFound(t, alias)

			testRedirect(t, updatedAlias, tc.url)

			// Delete

			resp = e.DELETE(fmt.Sprintf("/url/%s", updatedAlias)).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).JSON().Object()

			resp.Value("status").String().IsEqual("OK")

			// Redirect

			testRedirectNotFound(t, updatedAlias)

		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	_, err := api.GetRedirect(u.String())
	require.ErrorIs(t, err, api.ErrInvalidStatusCode)
}
