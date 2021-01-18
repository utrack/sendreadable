package handlers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

const jwtIssuer = "sendreadable.utrack.dev"

type jwtClaims struct {
	jwt.StandardClaims

	RmJWT string `json:"rmtok"`
	ID    string `json:"id"`
}

func jwtGen(key interface{}, rmTok string) (string, error) {
	//	jwt.NewWithClaims(jwt.SigningMethodRS512,
	claims := jwtClaims{RmJWT: rmTok}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Hour * 24 * 7 * 30).Unix()
	claims.Issuer = jwtIssuer
	id, _ := uuid.NewV4()
	claims.ID = id.String()

	tok := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return tok.SignedString(key)
}
