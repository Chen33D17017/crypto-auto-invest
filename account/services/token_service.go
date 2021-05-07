package services

import (
	"account-tutorial/model"
	"account-tutorial/model/apperrors"
	"context"
	"crypto/rsa"
	"log"
)

type tokenService struct {
	PrivKey       *rsa.PrivateKey
	PubKey        *rsa.PublicKey
	RefreshSecret string
	IDExpirationSecs int64
	RefreshExpirationSecs int64
}

type TSConfig struct {
	PrivKey       *rsa.PrivateKey
	PubKey        *rsa.PublicKey
	RefreshSecret string
	IDExpirationSecs int64
	RefreshExpirationSecs int64
}

func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		PrivKey:       c.PrivKey,
		PubKey:        c.PubKey,
		RefreshSecret: c.RefreshSecret,
		IDExpirationSecs: c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

func (s *tokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, s.PrivKey, s.IDExpirationSecs)

	if err != nil{
		log.Printf("Error generating idToken for uid: %v. Error: %v", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecs)
	if err != nil{
		log.Printf("Error generating refresh Token for uid: %v. Error: %v", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}
	
	// TODO: store refresh tokens by calling TokenRepository methods

	return &model.TokenPair{
		IDToken: idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}
