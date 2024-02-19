package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/Zeta-Manu/Backend/internal/config"
)

func AuthenticationMiddleware(cfg config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		// Fetch the public JWK from the Cognito endpoint
		keySet, err := fetchPublicJWTKey(context.Background(), cfg.JWT.JWTPublicKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch public JWK"})
			c.Abort()
			return
		}

		// Verify the Token
		token, err := verifyToken(clientToken, keySet)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func verifyToken(tokenString string, keySet jwk.Set) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}
		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, errors.New("cannot look up kid header")
		}
		var publickey interface{}
		if err := keys.Raw(&publickey); err != nil {
			return nil, err
		}
		return publickey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func fetchPublicJWTKey(ctx context.Context, link string) (jwk.Set, error) {
	keySet, err := jwk.Fetch(ctx, link)
	if err != nil {
		return nil, err
	}
	return keySet, nil
}
