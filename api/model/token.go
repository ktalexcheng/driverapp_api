package model

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ktalexcheng/trailbrake_api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	TokenString string
	Token       *jwt.Token
	Subject     string
	User        User
}

// Use custom type for context key to avoid collisions
type TokenKeyType string

const TokenKey TokenKeyType = "token"

var secretKey = os.Getenv("TOKEN_SECRET_KEY")

func (t *Token) VerifyToken(mg *util.MongoClient) error {
	token, err := jwt.Parse(t.TokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	if token.Valid {
		t.Token = token

		// Parse token claim
		claims, ok := t.Token.Claims.(jwt.MapClaims)
		if !ok {
			return errors.New("unable to parse token claims")
		}

		// Get subject of token claim
		subject, ok := claims["sub"].(string)
		if !ok {
			return errors.New("invalid subject")
		}

		// Subject is the full user ID
		t.Subject = subject

		// Check user exists
		_userId, err := primitive.ObjectIDFromHex(subject)
		if err != nil {
			return err
		}
		_user := User{ID: _userId}
		userExists, err := _user.CheckUserExists(mg)
		if err != nil {
			return err
		}
		if !userExists {
			return errors.New("user does not exist")
		}

		return nil
	} else {
		return errors.New("invalid token")
	}
}

func (t *Token) CreateToken(user *User) error {
	// Create new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Expiry
	claims["iat"] = time.Now().Unix()                     // Issued at
	claims["sub"] = (*user).ID                            // Set user ID as the subject

	// Generate the signed token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}

	t.Token = token
	t.TokenString = tokenString
	t.Subject = (*user).ID.Hex()
	t.User = *user

	return nil
}
