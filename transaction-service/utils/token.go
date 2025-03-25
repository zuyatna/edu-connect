package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func ValidateJWT(tokenString string) (string, string, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	err := godotenv.Load(".env")
	if err != nil {
		return "", "", fmt.Errorf("server error: failed to load environment variables")
	}

	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		return "", "", fmt.Errorf("server error: JWT secret is missing")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var email string
		if e, emailOk := claims["email"].(string); emailOk && e != "" {
			email = e
		} else {
			return "", "", fmt.Errorf("email not found in token claims")
		}

		var userID string
		if id, idOk := claims["id"].(string); idOk && id != "" {
			userID = id
		} else if id, idOk := claims["user_id"].(string); idOk && id != "" {
			userID = id
		}

		return userID, email, nil
	}

	return "", "", fmt.Errorf("invalid token")
}

// func ValidateToken(tokenString string) (*jwt.MapClaims, error) {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		return nil, err
// 	}

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("JWT_SECRET")), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		return &claims, nil
// 	}

// 	return nil, err
// }
