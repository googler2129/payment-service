package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/mercor/payment-service/constants"
	"github.com/mercor/payment-service/pkg/config"
	"github.com/mercor/payment-service/pkg/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthenticateJWT(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken, err := extractToken(c)
		if err != nil {
			log.Errorf("Not able to extract token from header")
			log.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(bearerToken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			rsaKey, cusErr := getRSAPublicKey(c, config.GetString(c, "authentication.rsaPublicKey"))
			if cusErr != nil {
				return nil, cusErr
			}

			return rsaKey, cusErr
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Errorf("Token is invalid :: %v", token)
			log.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*JWTClaim)

		c.Set(constants.UserDetails, &claims.UserDetails)
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	bearerToken := c.Request.Header.Get(constants.Authorization)
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1], nil
	}
	return "", errors.New("no auth token found")
}

type JWTClaim struct {
	jwt.StandardClaims
	UserDetails UserDetails
}

type UserDetails struct {
	ID    string
	Email string
}

func getRSAPublicKey(ctx context.Context, publicKey string) (rsaPubKey *rsa.PublicKey, err error) {
	pubPem, _ := pem.Decode([]byte(publicKey))

	if pubPem.Type != "PUBLIC KEY" {
		err = errors.New(fmt.Sprintf("RSA public key is of the wrong type, Pem Type :%s", pubPem.Type))
		return
	}

	parsedKey, parseErr := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if parseErr != nil {
		err = errors.New(fmt.Sprintf("RSA public key is of the wrong type, Pem Type :%s", pubPem.Type))
		return
	}

	rsaPubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		err = errors.New("unable to parse pub key")
		return
	}

	return
}
