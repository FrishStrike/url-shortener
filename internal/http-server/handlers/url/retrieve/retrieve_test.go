package retrieve_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"start1/internal/http-server/handlers/url/retrieve"
	"start1/internal/http-server/mocks"
	"start1/internal/storage"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetrieveHandler(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		wantUrl    string
		wantStatus int
		respError  string
		expectMock bool
		mockError  error
	}{
		{
			name:       "Success",
			alias:      "Google",
			wantUrl:    "https://google.com",
			wantStatus: http.StatusOK,
			expectMock: true,
		},
		{
			name:       "Empty alias",
			alias:      "",
			wantUrl:    "",
			wantStatus: http.StatusBadRequest,
			respError:  "field Alias is a required field",
			expectMock: false,
		},
		{
			name:       "Invalid alias",
			alias:      "InvalidAlias",
			wantUrl:    "",
			respError:  "url is not found",
			mockError:  storage.ErrURLNotFound,
			wantStatus: http.StatusBadRequest,
			expectMock: true,
		},
		{
			name:       "GetURL error",
			alias:      "yandex.ru",
			wantUrl:    "https://yandex.ru",
			respError:  "failed to get URL",
			mockError:  errors.New("unexpected error"),
			wantStatus: http.StatusBadRequest,
			expectMock: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			urlRetrieverMock := mocks.NewURLRetriever(t)
			if tc.expectMock {
				urlRetrieverMock.On("GetURL", tc.alias).
					Return(tc.wantUrl, tc.mockError).
					Once()
			}

			handler := retrieve.New(slog.New(slog.DiscardHandler), urlRetrieverMock)
			input := fmt.Sprintf(`{"alias": "%s"}`, tc.alias)

			req, err := http.NewRequest("GET", "/get", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)
			require.Equal(t, tc.wantStatus, rr.Code)
			body := rr.Body.String()
			var resp retrieve.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
