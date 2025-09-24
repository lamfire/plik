package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"

	"github.com/root-gg/plik/server/common"
)

func TestGetUserFromJWT(t *testing.T) {
	ctx := newTestingContext(common.NewConfiguration())
	ctx.GetConfig().FeatureAuthentication = common.FeatureEnabled

	// Create a test user
	user := &common.User{
		ID:       "test-user-id",
		Provider: common.ProviderLocal,
		Login:    "testuser",
		Name:     "Test User",
		Email:    "test@example.com",
	}

	// Create authenticator with test signature key
	signatureKey := "test-signature-key"
	authenticator := &common.SessionAuthenticator{
		SignatureKey:   signatureKey,
		SecureCookies:  false,
		SessionTimeout: 3600,
		Path:           "/",
	}
	ctx.SetAuthenticator(authenticator)

	// Create user in metadata backend
	err := ctx.GetMetadataBackend().CreateUser(user)
	require.NoError(t, err, "unable to create test user")

	// Test 1: Valid JWT token
	t.Run("ValidJWTToken", func(t *testing.T) {
		// Generate JWT token
		token := jwt.New(jwt.SigningMethodHS512)
		claims := token.Claims.(jwt.MapClaims)
		claims["uid"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour).Unix()
		claims["iat"] = time.Now().Unix()

		tokenString, err := token.SignedString([]byte(signatureKey))
		require.NoError(t, err, "unable to sign JWT token")

		// Create request with JWT token
		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{}))
		require.NoError(t, err, "unable to create request")

		req.Header.Set("Authorization", "Bearer "+tokenString)
		ctx.SetReq(req)

		// Test JWT authentication
		authenticatedUser, err := getUserFromJWT(ctx)
		require.NoError(t, err, "JWT authentication should succeed")
		require.NotNil(t, authenticatedUser, "authenticated user should not be nil")
		require.Equal(t, user.ID, authenticatedUser.ID, "user ID should match")
		require.Equal(t, user.Login, authenticatedUser.Login, "user login should match")
	})

	// Test 2: Invalid JWT token
	t.Run("InvalidJWTToken", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{}))
		require.NoError(t, err, "unable to create request")

		req.Header.Set("Authorization", "Bearer invalid-token")
		ctx.SetReq(req)

		// Test JWT authentication
		authenticatedUser, err := getUserFromJWT(ctx)
		require.Error(t, err, "JWT authentication should fail")
		require.Nil(t, authenticatedUser, "authenticated user should be nil")
	})

	// Test 3: Missing Authorization header
	t.Run("MissingAuthorizationHeader", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{}))
		require.NoError(t, err, "unable to create request")

		ctx.SetReq(req)

		// Test JWT authentication
		authenticatedUser, err := getUserFromJWT(ctx)
		require.NoError(t, err, "should not error when no Authorization header")
		require.Nil(t, authenticatedUser, "authenticated user should be nil")
	})

	// Test 4: Invalid Authorization header format
	t.Run("InvalidAuthorizationHeaderFormat", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{}))
		require.NoError(t, err, "unable to create request")

		req.Header.Set("Authorization", "InvalidFormat token")
		ctx.SetReq(req)

		// Test JWT authentication
		authenticatedUser, err := getUserFromJWT(ctx)
		require.NoError(t, err, "should not error for invalid format")
		require.Nil(t, authenticatedUser, "authenticated user should be nil")
	})

	// Test 5: Expired JWT token
	t.Run("ExpiredJWTToken", func(t *testing.T) {
		// Generate expired JWT token
		token := jwt.New(jwt.SigningMethodHS512)
		claims := token.Claims.(jwt.MapClaims)
		claims["uid"] = user.ID
		claims["exp"] = time.Now().Add(-time.Hour).Unix() // Expired 1 hour ago
		claims["iat"] = time.Now().Add(-time.Hour * 2).Unix()

		tokenString, err := token.SignedString([]byte(signatureKey))
		require.NoError(t, err, "unable to sign JWT token")

		// Create request with expired JWT token
		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{}))
		require.NoError(t, err, "unable to create request")

		req.Header.Set("Authorization", "Bearer "+tokenString)
		ctx.SetReq(req)

		// Test JWT authentication
		authenticatedUser, err := getUserFromJWT(ctx)
		require.Error(t, err, "JWT authentication should fail for expired token")
		require.Nil(t, authenticatedUser, "authenticated user should be nil")
	})

	// Test 6: JWT token with wrong signature
	t.Run("JWTTokenWithWrongSignature", func(t *testing.T) {
		// Generate JWT token with wrong signature key
		token := jwt.New(jwt.SigningMethodHS512)
		claims := token.Claims.(jwt.MapClaims)
		claims["uid"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour).Unix()
		claims["iat"] = time.Now().Unix()

		tokenString, err := token.SignedString([]byte("wrong-signature-key"))
		require.NoError(t, err, "unable to sign JWT token")

		// Create request with JWT token
		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{}))
		require.NoError(t, err, "unable to create request")

		req.Header.Set("Authorization", "Bearer "+tokenString)
		ctx.SetReq(req)

		// Test JWT authentication
		authenticatedUser, err := getUserFromJWT(ctx)
		require.Error(t, err, "JWT authentication should fail for wrong signature")
		require.Nil(t, authenticatedUser, "authenticated user should be nil")
	})

	// Test 7: JWT token with missing user ID
	t.Run("JWTTokenWithMissingUserID", func(t *testing.T) {
		// Generate JWT token without user ID
		token := jwt.New(jwt.SigningMethodHS512)
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour).Unix()
		claims["iat"] = time.Now().Unix()
		// Missing "uid" claim

		tokenString, err := token.SignedString([]byte(signatureKey))
		require.NoError(t, err, "unable to sign JWT token")

		// Create request with JWT token
		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{}))
		require.NoError(t, err, "unable to create request")

		req.Header.Set("Authorization", "Bearer "+tokenString)
		ctx.SetReq(req)

		// Test JWT authentication
		authenticatedUser, err := getUserFromJWT(ctx)
		require.Error(t, err, "JWT authentication should fail for missing user ID")
		require.Nil(t, authenticatedUser, "authenticated user should be nil")
	})

	// Test 8: JWT token with non-existent user
	t.Run("JWTTokenWithNonExistentUser", func(t *testing.T) {
		// Generate JWT token with non-existent user ID
		token := jwt.New(jwt.SigningMethodHS512)
		claims := token.Claims.(jwt.MapClaims)
		claims["uid"] = "non-existent-user-id"
		claims["exp"] = time.Now().Add(time.Hour).Unix()
		claims["iat"] = time.Now().Unix()

		tokenString, err := token.SignedString([]byte(signatureKey))
		require.NoError(t, err, "unable to sign JWT token")

		// Create request with JWT token
		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{}))
		require.NoError(t, err, "unable to create request")

		req.Header.Set("Authorization", "Bearer "+tokenString)
		ctx.SetReq(req)

		// Test JWT authentication
		authenticatedUser, err := getUserFromJWT(ctx)
		require.Error(t, err, "JWT authentication should fail for non-existent user")
		require.Nil(t, authenticatedUser, "authenticated user should be nil")
	})
}

func TestAuthenticateWithJWT(t *testing.T) {
	ctx := newTestingContext(common.NewConfiguration())
	ctx.GetConfig().FeatureAuthentication = common.FeatureEnabled

	// Create a test user
	user := &common.User{
		ID:       "test-user-id",
		Provider: common.ProviderLocal,
		Login:    "testuser",
		Name:     "Test User",
		Email:    "test@example.com",
	}

	// Create authenticator with test signature key
	signatureKey := "test-signature-key"
	authenticator := &common.SessionAuthenticator{
		SignatureKey:   signatureKey,
		SecureCookies:  false,
		SessionTimeout: 3600,
		Path:           "/",
	}
	ctx.SetAuthenticator(authenticator)

	// Create user in metadata backend
	err := ctx.GetMetadataBackend().CreateUser(user)
	require.NoError(t, err, "unable to create test user")

	// Test JWT authentication through middleware
	t.Run("JWTAuthenticationThroughMiddleware", func(t *testing.T) {
		// Generate JWT token
		token := jwt.New(jwt.SigningMethodHS512)
		claims := token.Claims.(jwt.MapClaims)
		claims["uid"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour).Unix()
		claims["iat"] = time.Now().Unix()

		tokenString, err := token.SignedString([]byte(signatureKey))
		require.NoError(t, err, "unable to sign JWT token")

		// Create request with JWT token
		req, err := http.NewRequest("POST", "/upload", bytes.NewBuffer([]byte{}))
		require.NoError(t, err, "unable to create request")

		req.Header.Set("Authorization", "Bearer "+tokenString)
		ctx.SetReq(req)

		// Create response recorder
		rr := httptest.NewRecorder()

		// Test authentication middleware
		authenticateMiddleware := Authenticate(true)
		handler := authenticateMiddleware(ctx, common.DummyHandler)
		handler.ServeHTTP(rr, req)

		// Check that user is set in context
		authenticatedUser := ctx.GetUser()
		require.NotNil(t, authenticatedUser, "user should be set in context")
		require.Equal(t, user.ID, authenticatedUser.ID, "user ID should match")
		require.Equal(t, user.Login, authenticatedUser.Login, "user login should match")
	})
}
