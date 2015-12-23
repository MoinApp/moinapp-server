package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MoinApp/moinapp-server/models"
	"github.com/dgrijalva/jwt-go"
)

const (
	TokenValidity = time.Hour * 24 * 3
)

var (
	ErrAuthTokenMissing = errors.New("Auth token missing.")
	ErrTokenInBlacklist = errors.New("Token is in blacklist.")
	ErrTokenChanged     = errors.New("The token has been tampered with.")
)

// GenerateJWTToken generates a JWT token for a given UserID and signs it with
// the given private key. The token will be valid for 3 days.
func GenerateJWTToken(user models.User) (string, error) {

	if user.PrivateKey == "" {
		privateKey, err := GenerateNewPrivateKey()

		if err != nil {
			return "", err
		}

		user.PrivateKey = PrivateKeyToString(privateKey)
		models.SaveUser(&user)
	}

	privateKey, err := StringToPrivateKey(user.PrivateKey)

	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodRS256)

	token.Claims["user"] = user.ID
	token.Header["user"] = user.ID
	token.Claims["exp"] = time.Now().Add(TokenValidity).Unix()

	tokenString, err := token.SignedString(privateKey)

	return tokenString, err
}

// ValidateJWTToken validates a JWT token and returns the user from the DB.
func ValidateJWTToken(input string) (models.User, error) {
	var user models.User

	if isInBlacklist(input) {
		return models.User{}, ErrTokenInBlacklist
	}

	token, err := jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {

		// Check whether the right signing algorithm was used.
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Get the user ID
		userID := token.Header["user"]

		user = *models.FindUserById(userID)

		privateKey, err := StringToPrivateKey(user.PrivateKey)

		return privateKey.Public(), err
	})

	if token.Claims["user"] != token.Header["user"] {
		return models.User{}, ErrTokenChanged
	}

	if err == nil && token.Valid {
		return user, nil

	}
	return models.User{}, err
}

func isInBlacklist(tokenString string) bool {
	/* _, err := models.RedisInstance.Get(tokenString).Result()
	if err == redis.Nil {
		return false
	}
	return true */
	return false
}

func getRemainingTokenValidity(input string) int {
	var user models.User

	token, err := jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {

		// Check whether the right signing algorithm was used.
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Get the user ID
		userID := token.Header["user"]

		user = *models.FindUserById(userID)

		privateKey, err := StringToPrivateKey(user.PrivateKey)

		return privateKey.Public(), err
	})

	if err != nil {
		return 3600
	}

	timestamp := token.Claims["exp"]

	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + 3600)
		}
	}
	return 3600
}

// ValidateSession validates a session in a HTTP Request
func ValidateSession(r *http.Request) (models.User, error) {
	token := r.Header.Get("Session")

	if token == "" {
		return models.User{}, ErrAuthTokenMissing
	}

	return ValidateJWTToken(token)
}

// InvalidateToken invalidates a given token
func InvalidateToken(token string) error {
	//return models.RedisInstance.Set(token, token, time.Duration(getRemainingTokenValidity(token))*time.Second).Err()
	return nil
}
