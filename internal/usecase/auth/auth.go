package auth

import (
	"bhs/internal/entity"
	"bhs/internal/usecase"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.devnw.com/ttl"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Auth -.
type Auth struct {
	repo        usecase.UserRepo
	cache       ttl.Cache[string, entity.User]
	tokenTtl    time.Duration
	tokenSecret string
}

// NewAuth -.
func NewAuth(r usecase.UserRepo, tokenSecret string) *Auth {
	tokenTtl := time.Hour * 12
	return &Auth{
		repo:        r,
		cache:       *ttl.NewCache[string, entity.User](nil, tokenTtl, true),
		tokenTtl:    tokenTtl,
		tokenSecret: tokenSecret,
	}
}

// Authenticate -.
func (a *Auth) Authenticate(ctx context.Context, username string, password string) (string, error) {
	userEntity, err := a.repo.GetUser(ctx, username)
	if err != nil {
		return "", err
	}
	if checkPasswordHash(password, userEntity.PasswordHash) {
		return a.createToken(ctx, *userEntity)
	} else {
		return "", fmt.Errorf("error checking password")
	}
}

// Register -.
func (a *Auth) Register(ctx context.Context, user entity.User) (string, error) {
	userEntity, err := a.repo.SaveUser(ctx, user)
	if err != nil {
		return "", err
	}
	return a.createToken(ctx, *userEntity)
}

// Logout -.
func (a *Auth) Logout(ctx context.Context, token string) {
	a.cache.Delete(ctx, token)
}

// HashPassword -.
func (a *Auth) HashPassword(rawPassword string) (string, error) {
	return hashPassword(rawPassword)
}

// Authorize -.
func (a *Auth) Authorize(ctx context.Context, token string) (entity.User, error) {
	cached, found := a.cache.Get(nil, token)
	if !found {
		return a.parseToken(token)
	} else {
		return cached, nil
	}
}

func (a *Auth) createToken(ctx context.Context, user entity.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(a.tokenTtl)

	tokenString, err := token.SignedString([]byte(a.tokenSecret))
	if err != nil {
		return "", fmt.Errorf("token creation error: %w", err)
	}

	userCache := entity.User{
		Id:       user.Id,
		Username: user.Username,
	}
	err = a.cache.SetTTL(nil, tokenString, userCache, a.tokenTtl)

	if err != nil {
		return "", fmt.Errorf("cache error: %w", err)
	}
	return tokenString, nil
}

func (a *Auth) parseToken(token string) (entity.User, error) {
	now := time.Now()
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.tokenSecret), nil
	})
	if err != nil {
		return entity.User{}, err
	}
	claims := tokenParsed.Claims.(jwt.MapClaims)
	expiry := claims["exp"].(time.Duration)
	if expiry.Milliseconds() < now.UnixMilli() {
		return entity.User{}, fmt.Errorf("token expired")
	}
	userData := entity.User{
		Id:       claims["id"].(int64),
		Username: claims["username"].(string),
	}

	return userData, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
