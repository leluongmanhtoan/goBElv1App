package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"program/repository"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IAuth interface {
	GenerateToken(userId string, isRefeshToken bool) (string, error)
	ValidateToken(token string, isRefreshToken bool) (*jwt.Token, error)
	RevokeSession(accessToken, refreshToken string) error
}

type JwtAuth struct {
	SecretKey string
	Issuer    string
}

func (s *JwtAuth) GenerateToken(userId string, isRefeshToken bool) (string, error) {
	tokenID := uuid.NewString()
	expireAt := time.Now().Unix()
	if !isRefeshToken {
		expireAt = time.Now().Add(15 * time.Minute).Unix()
	} else {
		expireAt = time.Now().AddDate(0, 0, 7).Unix()
	}
	claims := &jwt.StandardClaims{
		Id:        tokenID,
		Subject:   userId,
		IssuedAt:  time.Now().Unix(),
		Issuer:    s.Issuer,
		ExpiresAt: expireAt,
		NotBefore: time.Now().Add(5 * time.Second).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	if isRefeshToken {
		_, err := repository.RedisClientConnection.SAdd("userId:"+userId+":tokenID_valid", tokenID)
		if err != nil {
			return "", err
		}
	}

	return token, nil
}

func (s *JwtAuth) ValidateToken(token string, isRefreshToken bool) (*jwt.Token, error) {
	if !strings.HasPrefix(token, "Bearer ") && !isRefreshToken {
		return nil, fmt.Errorf("not a Bearer authorization")
	}
	keyFunc := func(t_ *jwt.Token) (any, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(s.SecretKey), nil
	}
	tokenString := strings.TrimPrefix(token, "Bearer ")
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("can not parse token: %v", err)
	}
	claims, ok := parsedToken.Claims.(*jwt.StandardClaims)
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid claims")
	}
	if !ok {
		return nil, fmt.Errorf("can't claims token")
	}
	if isRefreshToken {
		isValidToken, err := repository.RedisClientConnection.IsMemberInSet("userId:"+claims.Subject+":tokenID_valid", claims.Id)
		if !isValidToken || err != nil {
			return nil, fmt.Errorf("refresh token is invalid")
		}
	} else {
		isBannedToken, err := repository.RedisClientConnection.IsExisted("blacklist:accessToken:" + token)
		if isBannedToken || err != nil {
			return nil, fmt.Errorf("access token is invalid")
		}
	}
	now := time.Now().Unix()
	if claims.IssuedAt > now {
		return nil, fmt.Errorf("token issued in the future")
	}
	if claims.NotBefore > now {
		return nil, fmt.Errorf("token not yet valid")
	}
	if claims.ExpiresAt < now {
		return nil, fmt.Errorf("token expired")
	}

	return parsedToken, nil
}

// Private check revoke token
// func (s *JwtAuth) isTokenRevoked(tokenId string) bool {
func (s *JwtAuth) RevokeSession(accessToken, refreshToken string) error {
	keyFunc := func(t_ *jwt.Token) (any, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(s.SecretKey), nil
	}
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, keyFunc)
	if err != nil {
		return fmt.Errorf("can not parse token: %v", err)
	}
	claims, ok := parsedToken.Claims.(*jwt.StandardClaims)
	if !parsedToken.Valid {
		return fmt.Errorf("invalid claims")
	}
	if !ok {
		return fmt.Errorf("can't claims token")
	}

	_, err = repository.RedisClientConnection.SRem("userId:"+claims.Subject+":tokenID_valid", claims.Id)
	if err != nil {
		return err
	}
	_, err = repository.RedisClientConnection.SetTTL("blacklist:accessToken:"+accessToken, "revoked", time.Duration(30*time.Minute))
	if err != nil {
		return err
	}

	return nil
}

// Password Handler for User Login and Register
type PasswordHandler struct {
	SaltSize int
}

func (p *PasswordHandler) GenerateSalt() (salt string, err error) {
	saltByte := make([]byte, p.SaltSize)
	if _, err := rand.Read(saltByte); err != nil {
		return "", err
	}
	salt = base64.StdEncoding.EncodeToString(saltByte)
	return
}

func (p *PasswordHandler) HashPassword(password, salt string) (string, error) {
	saltedPassword := password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (p *PasswordHandler) ValidatePassword(hash, password, salt string) error {
	saltedPassword := password + salt
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltedPassword))
}
