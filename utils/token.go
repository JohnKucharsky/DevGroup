package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type TokenDetails struct {
	Token     *string
	TokenUUID string
	UserID    string
	ExpiresIn *int64
}

func CreateToken(
	userID string,
	ttl time.Duration,
	privateKey string,
) (*TokenDetails, error) {
	now := time.Now().UTC()
	td := &TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}
	*td.ExpiresIn = now.Add(ttl).Unix()
	td.UserID = userID
	tokenUUID, _ := uuid.NewV4()
	td.TokenUUID = tokenUUID.String()

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf(
			"could not decode token private key: %s",
			err.Error(),
		)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return nil, fmt.Errorf(
			"create: parse token private key: %s",
			err.Error(),
		)
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = userID
	atClaims["token_uuid"] = td.TokenUUID
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	*td.Token, err = jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		atClaims,
	).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %s", err.Error())
	}

	return td, nil
}

func ValidateToken(token string, publicKey string) (*TokenDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %s", err.Error())
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %s", err.Error())
	}

	parsedToken, err := jwt.Parse(
		token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
			}
			return key, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("validate: %s", err.Error())
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return &TokenDetails{
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
		UserID:    fmt.Sprint(claims["sub"]),
	}, nil
}
