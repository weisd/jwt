package jwt

import (
	"errors"
	"fmt"
	// "net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const (
	Bearer        = "Bearer"
	JWTContextKey = "JWTContextKey"
)

var UnauthorizedErr = errors.New("JWT Authorization Failed")

func Claims(value interface{}) map[string]interface{} {
	switch v := value.(type) {
	case *echo.Context:
		return v.Get(JWTContextKey).(map[string]interface{})
	default:
		return nil
	}
}

type JWTKeyFunc func(c *echo.Context) (string, error)

// A JSON Web Token middleware
func EchoJWTAuther(keyFunc JWTKeyFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {

		// Skip WebSocket
		if (c.Request().Header.Get(echo.Upgrade)) == echo.WebSocket {
			return nil
		}

		// he := echo.NewHTTPError(http.StatusUnauthorized)
		he := UnauthorizedErr

		key, err := keyFunc(c)
		if err != nil {
			return he
		}

		auth := c.Request().Header.Get("Authorization")
		l := len(Bearer)

		if len(auth) > l+1 && auth[:l] == Bearer {
			t, err := jwt.Parse(auth[l+1:], func(token *jwt.Token) (interface{}, error) {

				// Always check the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// Return the key for validation
				return []byte(key), nil
			})
			if err == nil && t.Valid {
				// Store token claims in echo.Context
				c.Set(JWTContextKey, t.Claims)
				return nil
			}
		}
		return he
	}
}

func NewToken(key string, claims ...map[string]interface{}) (string, error) {
	// New web token.
	token := jwt.New(jwt.SigningMethodHS256)

	// Set a header and a claim
	token.Header["typ"] = "JWT"
	token.Claims["exp"] = time.Now().Add(time.Second * 60).Unix()

	if len(claims) > 0 {
		for k, v := range claims[0] {
			token.Claims[k] = v
		}
	}

	// Generate encoded token
	return token.SignedString([]byte(key))
}
