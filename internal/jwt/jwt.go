package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type Payload struct {
	UserID string
}

// Encode encodes a jwt token using data gotten from payload.
func Encode(secretKey []byte, payload Payload) (tokenString string, err error) {
	if len(secretKey) < 1 {
		return "", fmt.Errorf("secret key must be provided")
	}
	claims := jwt.MapClaims{
		"userid": payload.UserID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Decode decodes a jwt token string.
//
// If the jwt token is invalid it returns an error.
func Decode(secretKey []byte, tokenString string) (payload *Payload, err error) {
	if len(secretKey) < 1 {
		return nil, fmt.Errorf("secret key must be provided")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("invalid jwt token string")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		payload = &Payload{
			UserID: interfaceToStr(claims["userid"]),
		}
		return payload, nil
	}
	err = fmt.Errorf("failed to decode jwt")
	return nil, err
}

func interfaceToStr(str interface{}) string {
	switch str := str.(type) {
	case string:
		return str

	default:
		return ""
	}
}
