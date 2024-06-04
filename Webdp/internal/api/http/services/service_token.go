package services

import (
	"time"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/repo/postgres"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService struct {
	postg         postgres.TokenPostgres
	signKey       []byte
	signingMethod jwt.SigningMethod
}

type JWTTokenClaims struct {
	Handle string   `json:"handle"`
	Roles  []string `json:"roles"`
	jwt.StandardClaims
}

func (j *JWTTokenClaims) Valid() error {
	if j.ExpiresAt == 0 || j.Handle == "" {
		return errors.ErrInvalidToken
	}
	return nil
}

func NewTokenService(tokenRepo postgres.TokenPostgres, signingKey []byte, signingMethod jwt.SigningMethod) TokenService {
	return TokenService{postg: tokenRepo, signKey: signingKey, signingMethod: signingMethod}
}

func (t TokenService) IssueNewTokenFor(user string, roles []string) (string, error) {
	c := JWTTokenClaims{}
	var repoFun func(string, string) error

	exists, err := t.postg.UserTokenExists(user)
	if err != nil {
		return "", errors.WrapDBError(err, "token exists", user)
	}

	if exists {
		repoFun = t.postg.UpdateToken
	} else {
		repoFun = t.postg.SaveToken
	}

	c.ExpiresAt = time.Now().Add(time.Hour).Unix()
	c.Handle = user
	c.Roles = roles

	token := jwt.NewWithClaims(t.signingMethod, &c)
	signed, err := token.SignedString(t.signKey)
	if err != nil {
		return "", err
	}

	// maybe move to repo that ensures that only one token per user is saved
	if err := repoFun(user, signed); err != nil {
		return "", errors.WrapDBError(err, "token", user)
	}
	return signed, nil
}

func (t TokenService) LogOffUser(userHandle string) error {
	if err := t.postg.DeleteUserToken(userHandle); err != nil {
		return errors.WrapDBError(err, "delete token for", userHandle)
	}
	return nil
}
