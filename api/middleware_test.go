package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/vipeergod123/simple_bank/token"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMarker token.Marker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMarker.CreateToken(username, duration)
	require.NoError(t, err)
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Marker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		//test case go here
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "username", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UnAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, "unsupported type", "username", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, "", "username", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Marker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "username", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(server.tokenMarker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMarker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})
	}
}
