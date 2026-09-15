package cognito

import (
	"fmt"
	"strings"

	"awsems/internal/config"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTVerifier struct {
	Issuer   string
	ClientID string
	JWKSURL  string
}

func NewJWTVerifier(cfg *config.Config) (*JWTVerifier, error) {
	jwksURL := strings.TrimRight(
		cfg.CognitoOpenIDConfigURL,
		"/",
	) + "/.well-known/jwks.json"

	return &JWTVerifier{
		Issuer:   cfg.CognitoOpenIDConfigURL,
		ClientID: cfg.CognitoClientID,
		JWKSURL:  jwksURL,
	}, nil
}

func (v *JWTVerifier) VerifyAccessToken(tokenString string) (*jwt.Token, error) {
	// Fetch the JWKS keys on each call for verification
	jwks, err := keyfunc.Get(v.JWKSURL, keyfunc.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to load Cognito JWKS: %w", err)
	}

	token, err := jwt.Parse(
		tokenString,
		jwks.Keyfunc,                            // get the kid for a particular JWT from header
		jwt.WithValidMethods([]string{"RS256"}), // validating signature
		jwt.WithIssuer(v.Issuer),
	)

	if err != nil {
		return nil, err
	}

	// fmt.Println("Token:", token)
	// converting token details to map
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify the token has the correct type
	tokenUse, ok := claims["token_use"].(string)
	if !ok || tokenUse != "access" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Verify the token has the correct client
	clientID, ok := claims["client_id"].(string)
	if !ok || clientID != v.ClientID {
		return nil, fmt.Errorf("invalid client")
	}

	return token, nil
}
