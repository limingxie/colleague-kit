package auth

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256
)

const (
	JwtSecret = "JwtSecret"
)

func init() {
}

func NewToken(m map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"iss": "colleague",
		"aud": "colleague",
		"nbf": time.Now().Add(-time.Minute * 1).Unix(),
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}
	for k, v := range m {
		claims[k] = v
	}
	return jwt.NewWithClaims(jwtSigningMethod, claims).SignedString([]byte(JwtSecret))
}

func NewTokenForTest(m map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"iss": "colleague",
		"aud": "colleague",
		"nbf": time.Now().Add(-time.Minute * 1).Unix(),
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}
	for k, v := range m {
		claims[k] = v
	}
	return jwt.NewWithClaims(jwtSigningMethod, claims).SignedString([]byte(JwtSecret))
}

func Extract(token string) (jwt.MapClaims, error) {
	return ExtractWithSecret(token, JwtSecret)
}
func ExtractWithSecret(token, jwtSecret string) (jwt.MapClaims, error) {
	if token == "" {
		return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Required authorization token not found."}
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { return []byte(JwtSecret), nil })
	if err != nil {
		return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: fmt.Sprintf("Error parsing token: %v", err)}
	}

	if jwtSigningMethod != nil && jwtSigningMethod.Alg() != parsedToken.Header["alg"] {
		return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: fmt.Sprintf("Expected %s signing method but token specified %s",
			jwtSigningMethod.Alg(),
			parsedToken.Header["alg"])}
	}

	if !parsedToken.Valid {
		return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Token is invalid"}
	}

	claimInfo, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid token"}
	}
	return claimInfo, nil
}
