package helpers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/lukachi/blob-svc/internal/data"
	"log"
	"time"
)

type UserClaims struct {
	ID       string
	Username string
	jwt.RegisteredClaims
}

type JWTManager interface {
	New() JWTManager

	Gen(claims interface{}) (data.AuthTokens, error)

	ParseAccessToken(accessToken string) (*UserClaims, error)
	ParseRefreshToken(accessToken string) (*jwt.RegisteredClaims, error)

	NewAccessToken(claims UserClaims) (string, error)
	NewRefreshToken(claims jwt.RegisteredClaims) (string, error)
}

type JWT struct {
	signingKey []byte
}

func NewJwtManager(signingKey []byte) JWTManager {
	return &JWT{
		signingKey: []byte(signingKey),
	}
}

func (j *JWT) New() JWTManager {
	return NewJwtManager(j.signingKey)
}

func (j *JWT) NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString(j.signingKey)
}

func (j *JWT) NewRefreshToken(claims jwt.RegisteredClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString(j.signingKey)
}

func (j *JWT) Gen(claims interface{}) (data.AuthTokens, error) {
	newClaims := UserClaims{
		ID:       claims.(UserClaims).ID,
		Username: claims.(UserClaims).Username,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute * 30),
			},
		},
	}

	signedAccessToken, err := j.NewAccessToken(newClaims)

	if err != nil {
		log.Default().Fatal(err)

		return data.AuthTokens{
			AccessToken:  "",
			RefreshToken: "",
		}, err
	}

	refreshedClaims := jwt.RegisteredClaims{
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(time.Hour * 48),
		},
	}

	signedRefreshToken, err := j.NewRefreshToken(refreshedClaims)

	if err != nil {
		return data.AuthTokens{
			AccessToken:  signedAccessToken,
			RefreshToken: "",
		}, err
	}

	return data.AuthTokens{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
		CreatedAt:    newClaims.RegisteredClaims.IssuedAt.Time,
		ExpiresAt:    newClaims.RegisteredClaims.ExpiresAt.Time,
	}, nil
}

func (j *JWT) ParseAccessToken(accessToken string) (*UserClaims, error) {
	var userClaims UserClaims

	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &userClaims, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})

	return parsedAccessToken.Claims.(*UserClaims), err
}

func (j *JWT) ParseRefreshToken(refreshToken string) (*jwt.RegisteredClaims, error) {
	var registeredClaims jwt.RegisteredClaims

	parsedAccessToken, err := jwt.ParseWithClaims(refreshToken, &registeredClaims, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})

	return parsedAccessToken.Claims.(*jwt.RegisteredClaims), err
}
