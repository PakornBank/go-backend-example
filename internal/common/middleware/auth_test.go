package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

const (
	testSecret   = "test-secret"
	bearerPrefix = "Bearer "
)

func setupAuthTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Auth(testSecret))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": c.MustGet("user_id"),
			"email":   c.MustGet("email"),
		})
	})
	return router
}

func generateTestToken(userID string, email string, expiry time.Duration) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(testSecret))
	return signedToken
}

func TestAuth(t *testing.T) {
	const (
		testID    = "test-user-id"
		testEmail = "test@email.com"
	)

	tests := []struct {
		name           string
		generateHeader func() string
		wantCode       int
		errContains    string
	}{
		{
			name: "valid token",
			generateHeader: func() string {
				return bearerPrefix + generateTestToken(testID, testEmail, time.Hour)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "expired token",
			generateHeader: func() string {
				return bearerPrefix + generateTestToken(testID, testEmail, -time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token",
		},
		{
			name: "invalid token",
			generateHeader: func() string {
				return bearerPrefix + "invalid-token"
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token",
		},
		{
			name: "empty authorization header",
			generateHeader: func() string {
				return ""
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "authorization header required",
		},
		{
			name: "missing Bearer prifix",
			generateHeader: func() string {
				return generateTestToken(testID, testEmail, time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid authorization header format",
		},
		{
			name: "wrong signing method",
			generateHeader: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
					"user_id": testID,
					"email":   testEmail,
					"exp":     time.Now().Add(time.Hour).Unix(),
				})
				signedToken, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
				return bearerPrefix + signedToken
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token",
		},
		{
			name: "invalid token claims",
			generateHeader: func() string {
				claims := jwt.MapClaims{
					"id": testID,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString([]byte(testSecret))
				return bearerPrefix + signedToken
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token claims",
		},
		{
			name: "missing user_id claim",
			generateHeader: func() string {
				return bearerPrefix + generateTestToken("", testEmail, time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token claims",
		},
		{
			name: "missing email claim",
			generateHeader: func() string {
				return bearerPrefix + generateTestToken(testID, "", time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token claims",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupAuthTest()

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if header := tt.generateHeader(); header != "" {
				req.Header.Set("Authorization", header)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &res)
			assert.NoError(t, err)

			if tt.wantCode == http.StatusOK {
				assert.Equal(t, testID, res["user_id"])
				assert.Equal(t, testEmail, res["email"])
			} else {
				assert.Contains(t, res["error"], tt.errContains)
			}
		})
	}
}
